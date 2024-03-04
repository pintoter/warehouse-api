package product

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pintoter/warehouse-api/internal/dbutil"
	"github.com/pintoter/warehouse-api/internal/repository"
	repoModel "github.com/pintoter/warehouse-api/internal/repository/model"
	"github.com/pintoter/warehouse-api/internal/service"
	"github.com/pintoter/warehouse-api/internal/service/model"
	"github.com/pintoter/warehouse-api/pkg/semaphore"
)

const (
	rejected        = "rejected: "
	reserved        = "reserved"
	released        = "released"
	goroutinesLimit = 10
)

type Service struct {
	repo      repository.Repository
	txManager dbutil.TxManager
}

func NewService(repo repository.Repository, txManager dbutil.TxManager) service.ProductService {
	return &Service{
		repo:      repo,
		txManager: txManager,
	}
}

func (s *Service) ReserveProducts(r *http.Request, args *model.ReserveProductsReq, reply *model.ReserveProductsResp) error {
	var (
		products      = args.Products
		productsInfo  []model.ReserveProductResp
		wg            sync.WaitGroup
		outputCh             = make(chan model.ReserveProductResp)
		reservationId string = uuid.New().String()
		sema          *semaphore.Semaphore
	)

	if len(products) == 0 {
		*reply = model.ReserveProductsResp{}
		return model.ErrInvalidInput
	}

	if len(products) > goroutinesLimit {
		sema = semaphore.New(goroutinesLimit)
	} else {
		sema = semaphore.New(len(products))
	}

	go func() {
		for _, product := range products {
			wg.Add(1)
			go s.processReservation(r.Context(), outputCh, &wg, sema, product, reservationId)
		}
		wg.Wait()
		close(outputCh)
	}()

	for res := range outputCh {
		productsInfo = append(productsInfo, res)
	}

	*reply = model.ReserveProductsResp{
		ReservationId:           reservationId,
		ReservationProductsInfo: productsInfo,
	}

	return nil
}

func (s *Service) processReservation(ctx context.Context, outputCh chan<- model.ReserveProductResp, wg *sync.WaitGroup, sema *semaphore.Semaphore, product model.ReserveProductReq, reservationId string) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	defer wg.Done()

	sema.Acquire()
	defer sema.Release()

	if product.Quantity <= 0 {
		outputCh <- model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()}
		return
	}

	err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		quantityProductsOnActiveWhs, err := s.repo.GetTotalQuantityOfProducts(ctx, product.Code)
		if err != nil {
			return model.ErrInvalidInput
		}

		if quantityProductsOnActiveWhs < product.Quantity {
			return model.ErrInvalidQuantity
		}

		// Get active warehouses sorted by quantity of products with warehouse code
		productsByWarehouses, err := s.repo.GetProductsByWarehousesByCode(ctx, product.Code)
		if err != nil {
			return model.ErrInternalServer
		}

		err = s.startReservation(ctx, productsByWarehouses, reservationId, product.Quantity)
		if err != nil {
			return err
		}

		return nil
	})

	switch {
	case ctx.Err() != nil:
		outputCh <- model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInternalServer.Error()}
	case err != nil:
		outputCh <- model.ReserveProductResp{Code: product.Code, Status: rejected + err.Error()}
	default:
		outputCh <- model.ReserveProductResp{Code: product.Code, Status: reserved}
	}
}

func (s *Service) startReservation(ctx context.Context, productsByWarehouses []repoModel.ProductsOnActiveWarehouse, reservationId string, quantity int) error {
	var err error
	// Begin reserving products from warehouses, starting from the warehouse with the maximum values
	for _, productsByWarehouse := range productsByWarehouses {
		var quantityLeftOnWarehouse, quantityForReservation int

		if productsByWarehouse.Quantity >= quantity {
			quantityLeftOnWarehouse = productsByWarehouse.Quantity - quantity
			quantityForReservation = quantity
			quantity = 0
		} else {
			quantityLeftOnWarehouse = 0
			quantityForReservation = productsByWarehouse.Quantity
			quantity -= productsByWarehouse.Quantity
		}

		err = s.repo.UpdateWarehouseQuantity(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, quantityLeftOnWarehouse)
		if err != nil {
			err = model.ErrInternalServer
			break
		}

		_, err = s.repo.CreateReservation(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, quantityForReservation, reservationId)
		if err != nil {
			err = model.ErrInternalServer
			break
		}

		if quantity == 0 {
			break
		}
	}
	return err
}

func (s *Service) ReleaseProducts(r *http.Request, args *model.ReleaseProductsReq, reply *model.ReleaseProductsResp) error {
	var (
		products     = args.Products
		productsInfo []model.ReleaseProductResp
		wg           sync.WaitGroup
		outputCh     = make(chan model.ReleaseProductResp)
		sema         *semaphore.Semaphore
	)

	if len(products) == 0 {
		*reply = model.ReleaseProductsResp{}
		return model.ErrInvalidInput
	}

	if len(products) > goroutinesLimit {
		sema = semaphore.New(goroutinesLimit)
	} else {
		sema = semaphore.New(len(products))
	}

	go func() {
		for _, product := range products {
			wg.Add(1)
			go s.processRelease(r.Context(), outputCh, &wg, sema, product)
		}

		wg.Wait()
		close(outputCh)
	}()

	for res := range outputCh {
		productsInfo = append(productsInfo, res)
	}

	*reply = model.ReleaseProductsResp{ReleaseProductsInfo: productsInfo}
	return nil
}

func (s *Service) processRelease(ctx context.Context, outputCh chan<- model.ReleaseProductResp, wg *sync.WaitGroup, sema *semaphore.Semaphore, product model.ReleaseProductReq) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	defer wg.Done()

	sema.Acquire()
	defer sema.Release()

	if product.Quantity <= 0 {
		outputCh <- model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()}
		return
	}

	err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		quantityProductsInReservation, err := s.repo.GetTotalQuantityOfReservation(ctx, product.ReservationId, product.Code)
		if err != nil {
			return model.ErrInvalidInput
		}

		if quantityProductsInReservation < product.Quantity {
			return model.ErrInvalidReservationQuantity
		}

		productsByWarehousesInReservation, err := s.repo.GetProductsByReservationByIdAndCode(ctx, product.ReservationId, product.Code)
		if err != nil {
			return model.ErrInvalidInput
		}

		s.startRelease(ctx, productsByWarehousesInReservation, product.Quantity)

		return nil
	})

	if err != nil {
		outputCh <- model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: rejected + err.Error()}
	} else {
		outputCh <- model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: released}
	}
}

func (s *Service) startRelease(ctx context.Context, productsByWarehousesInReservation []repoModel.ProductsInReservation, quantity int) error {
	var err error
	for _, productsByWarehouseInResevation := range productsByWarehousesInReservation {
		var remainInReservation, addToWarehouse int

		if productsByWarehouseInResevation.Quantity >= quantity {
			remainInReservation = productsByWarehouseInResevation.Quantity - quantity
			addToWarehouse = quantity
			quantity = 0
		} else {
			remainInReservation = 0
			addToWarehouse = productsByWarehouseInResevation.Quantity
			quantity -= productsByWarehouseInResevation.Quantity
		}

		err = s.repo.UpdateReservationQuantity(ctx,
			productsByWarehouseInResevation.ID,
			remainInReservation,
		)
		if err != nil {
			err = model.ErrInternalServer
			break
		}

		err = s.repo.UpdateWarehouseQuantityWithAdd(ctx,
			productsByWarehouseInResevation.WarehouseId,
			productsByWarehouseInResevation.ProductId,
			addToWarehouse,
		)
		if err != nil {
			err = model.ErrInternalServer
			break
		}

		if quantity == 0 {
			break
		}
	}
	return err
}

func (s *Service) GetProductsByWarehouse(r *http.Request, args *model.ShowProductsReq, reply *[]model.Product) error {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	products, err := s.repo.GetProductsByWarehouseId(ctx, args.WarehouseId)
	if err != nil {
		return err
	}

	*reply = products
	return nil
}
