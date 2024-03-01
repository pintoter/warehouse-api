package product

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pintoter/warehouse-api/internal/repository"
	"github.com/pintoter/warehouse-api/internal/service"
	"github.com/pintoter/warehouse-api/internal/service/model"
)

const (
	rejected = "rejected: "
)

type Service struct {
	repo repository.Repository
	// txManager db.TxManager
}

func NewService(repo repository.Repository) service.ProductService {
	return &Service{
		repo: repo,
	}
}

func (s *Service) ReserveProducts(r *http.Request, args *model.ReserveProductsReq, reply *model.ReserveProductsResp) error {
	var (
		products      = args.Products
		ProductsInfo  []model.ReserveProductResp
		wg            sync.WaitGroup
		reservationId string = uuid.New().String()
	)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	wg.Add(len(products))
	for _, product := range products {
		go func() {
			defer wg.Done()

			if product.Quantity <= 0 {
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()})
				return
			}

			countProductsInAllWhs, err := s.repo.GetTotalQuantityOfProducts(ctx, product.Code)
			if err != nil {
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()})
				return
			}

			if countProductsInAllWhs < product.Quantity {
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInvalidQuantity.Error()})
				return
			}

			productsByWarehouses, err := s.repo.GetProductsByWarehousesByCode(ctx, product.Code)
			if err != nil {
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()})
				return
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
					ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInternalServer.Error()})
					break
					/*
						ADD ROLLBACK
					*/
				}

				_, err = s.repo.CreateReservation(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, product.Quantity-quantityForReservation, reservationId)
				if err != nil {
					ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: rejected + model.ErrInternalServer.Error()})
					break
					/*
						ADD ROLLBACK
					*/
				}
				if product.Quantity == 0 {
					break
				} else {
					product.Quantity = quantityForReservation
				}
			}
			ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: "reserved"})
		}()
	}

	wg.Wait()

	*reply = model.ReserveProductsResp{
		ReservationId:           reservationId,
		ReservationProductsInfo: ProductsInfo,
	}

	return nil
}

func (s *Service) ReleaseProducts(r *http.Request, args *model.ReleaseProductsReq, reply *model.ReleaseProductsResp) error {
	var (
		products     = args.Products
		productsInfo []model.ReleaseProductResp
		wg           sync.WaitGroup
	)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	wg.Add(len(products))
	for _, product := range products {
		go func() {
			defer wg.Done()

			if product.Quantity <= 0 {
				productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()})
				return
			}

			quantityProductsInReservation, err := s.repo.GetTotalQuantityOfReservation(ctx, product.ReservationId, product.Code)
			if err != nil {
				productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: rejected + model.ErrInvalidInput.Error()})
				return
			}

			if quantityProductsInReservation < product.Quantity {
				productsInfo = append(
					productsInfo,
					model.ReleaseProductResp{
						ReservationId: product.ReservationId,
						Code:          product.Code,
						Status:        rejected + model.ErrInvalidReservationQuantity.Error(),
					})
				return
			}

			productsByWarehousesInReservation, err := s.repo.GetProductsByReservationByIdAndCode(ctx, product.ReservationId, product.Code)
			if err != nil {
				productsInfo = append(
					productsInfo,
					model.ReleaseProductResp{
						ReservationId: product.ReservationId,
						Code:          product.Code,
						Status:        rejected + model.ErrInvalidInput.Error(),
					},
				)
				return
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
					productsInfo = append(
						productsInfo,
						model.ReleaseProductResp{
							ReservationId: product.ReservationId,
							Code:          product.Code,
							Status:        rejected + model.ErrInternalServer.Error(),
						},
					)
					break
				}

				if err := s.repo.UpdateWarehouseQuantityWithAdd(ctx, productsByWarehouseInResevation.WarehouseId, productsByWarehouseInResevation.ProductId, addToWarehouse); err != nil {
					productsInfo = append(
						productsInfo,
						model.ReleaseProductResp{
							ReservationId: product.ReservationId,
							Code:          product.Code,
							Status:        rejected + model.ErrInternalServer.Error(),
						},
					)
					break
				}

				if product.Quantity == 0 {
					break
				}
			}
			productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: "released"})
		}()
	}

	wg.Wait()

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
