package product

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pintoter/warehouse-api/internal/dbutil"
	"github.com/pintoter/warehouse-api/internal/repository"
	"github.com/pintoter/warehouse-api/internal/service"
	"github.com/pintoter/warehouse-api/internal/service/model"
)

const (
	rejected = "rejected: "
	reserved = "reserved"
	released = "released"
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
	)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	go func() {
		for _, product := range products {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
					var errTx error

					if product.Quantity <= 0 {
						errTx = model.ErrInvalidInput
						return errTx
					}

					countProductsInAllWhs, err := s.repo.GetTotalQuantityOfProducts(ctx, product.Code)
					if err != nil {
						errTx = model.ErrInvalidInput
						return errTx
					}

					if countProductsInAllWhs < product.Quantity {
						errTx = model.ErrInvalidQuantity
						return errTx
					}

					productsByWarehouses, err := s.repo.GetProductsByWarehousesByCode(ctx, product.Code)
					if err != nil {
						errTx = model.ErrInvalidInput
						return errTx
					}

					for _, productsByWarehouse := range productsByWarehouses {
						var leftQuantityOnWarehouse, quantityForReservation int

						switch {
						case productsByWarehouse.Quantity >= product.Quantity:
							leftQuantityOnWarehouse = productsByWarehouse.Quantity - product.Quantity
							quantityForReservation = 0
						case productsByWarehouse.Quantity < product.Quantity:
							leftQuantityOnWarehouse = 0
							quantityForReservation = product.Quantity - productsByWarehouse.Quantity
						}

						if err := s.repo.UpdateWarehouseQuantity(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, leftQuantityOnWarehouse); err != nil {
							errTx = model.ErrInternalServer
							return errTx
						}

						_, err = s.repo.CreateReservation(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, product.Quantity-quantityForReservation, reservationId)
						if err != nil {
							errTx = model.ErrInternalServer
							return errTx
						}
						if product.Quantity == 0 {
							break
						} else {
							product.Quantity = quantityForReservation
						}
					}
					return nil
				})

				if err != nil {
					outputCh <- model.ReserveProductResp{Code: product.Code, Status: rejected + err.Error()}
				} else {
					outputCh <- model.ReserveProductResp{Code: product.Code, Status: reserved}
				}
			}()
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

func (s *Service) ReleaseProducts(r *http.Request, args *model.ReleaseProductsReq, reply *model.ReleaseProductsResp) error {
	var (
		products     = args.Products
		productsInfo []model.ReleaseProductResp
		wg           sync.WaitGroup
		outputCh     = make(chan model.ReleaseProductResp)
	)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	go func() {
		for _, product := range products {
			wg.Add(1)
			go func() {
				defer wg.Done()

				err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
					var errTx error

					if product.Quantity <= 0 {
						errTx = model.ErrInvalidInput
						return errTx
					}

					quantityProductsInReservation, err := s.repo.GetTotalQuantityOfReservation(ctx, product.ReservationId, product.Code)
					if err != nil {
						errTx = model.ErrInvalidInput
						return errTx
					}

					if quantityProductsInReservation < product.Quantity {
						errTx = model.ErrInvalidReservationQuantity
						return errTx
					}

					productsByWarehousesInReservation, err := s.repo.GetProductsByReservationByIdAndCode(ctx, product.ReservationId, product.Code)
					if err != nil {
						errTx = model.ErrInvalidInput
						return errTx
					}

					for _, productsByWarehouseInResevation := range productsByWarehousesInReservation {
						var remainInReservation, addToWarehouse int

						switch {
						case productsByWarehouseInResevation.Quantity >= product.Quantity:
							remainInReservation = productsByWarehouseInResevation.Quantity - product.Quantity
							addToWarehouse = product.Quantity
							product.Quantity = 0
						case productsByWarehouseInResevation.Quantity < product.Quantity:
							remainInReservation = 0
							addToWarehouse = productsByWarehouseInResevation.Quantity
							product.Quantity -= productsByWarehouseInResevation.Quantity
						}

						if err := s.repo.UpdateReservationQuantity(ctx, productsByWarehouseInResevation.ID, remainInReservation); err != nil {
							errTx = model.ErrInternalServer
							return errTx
						}

						if err := s.repo.UpdateWarehouseQuantityWithAdd(ctx, productsByWarehouseInResevation.WarehouseId, productsByWarehouseInResevation.ProductId, addToWarehouse); err != nil {
							errTx = model.ErrInternalServer
							return errTx
						}

						if product.Quantity == 0 {
							break
						}
					}
					return nil
				})

				if err != nil {
					outputCh <- model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: rejected + err.Error()}
				} else {
					outputCh <- model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: released}
				}
			}()
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
