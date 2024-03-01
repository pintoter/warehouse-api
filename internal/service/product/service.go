package product

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pintoter/warehouse-api/internal/model"
	"github.com/pintoter/warehouse-api/internal/repository"
	"github.com/pintoter/warehouse-api/internal/service"
	"github.com/pintoter/warehouse-api/pkg/logger"
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

func (s *Service) ReserveProducts(ctx context.Context, products []model.ReserveProductReq) model.ReserveProductsResp {
	layer := "service.ReserveProduct"
	var (
		ProductsInfo  []model.ReserveProductResp
		wg            sync.WaitGroup
		reservationId string = uuid.New().String()
	)

	wg.Add(len(products))
	for _, product := range products {
		go func() {
			defer wg.Done()

			if product.Quantity <= 0 {
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: model.ErrInvalidInput.Error()})
				return
			}

			countProductsInAllWhs, err := s.repo.GetTotalQuantityOfProducts(ctx, product.Code)
			logger.DebugKV(ctx, "res", "layer", layer, "countProductsInAllWhs", countProductsInAllWhs)
			if err != nil {
				// return model.ErrInvalidCode
				logger.DebugKV(ctx, "res", "layer", layer, "err", err.Error())
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: err.Error()})
				return
			}

			if countProductsInAllWhs < product.Quantity {
				logger.DebugKV(ctx, "res", "layer", layer, "err", "countProductsInAllWhs < product.Quantity")
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: model.ErrInvalidQuantity.Error()})
				return
			}

			productsByWarehouses, err := s.repo.GetProductsByWarehousesByCode(ctx, product.Code)
			logger.DebugKV(ctx, "res", "layer", layer, "productsByWarehouses", productsByWarehouses)
			if err != nil {
				logger.DebugKV(ctx, "res", "layer", layer, "err", err.Error())
				ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: err.Error()})
				return
			}

			for _, productsByWarehouse := range productsByWarehouses {
				logger.DebugKV(ctx, "res", "layer", layer, "productsByWarehouse", productsByWarehouse)
				var leftQuantityOnWarehouse, quantityForReservation int
				switch {
				case productsByWarehouse.Quantity >= product.Quantity:
					leftQuantityOnWarehouse = productsByWarehouse.Quantity - product.Quantity
					quantityForReservation = 0
				case productsByWarehouse.Quantity < product.Quantity:
					leftQuantityOnWarehouse = 0
					quantityForReservation = product.Quantity - productsByWarehouse.Quantity
				}
				logger.DebugKV(ctx, "res", "layer", layer, "leftQuantityOnWarehouse", quantityForReservation)
				if err := s.repo.UpdateWarehouseQuantity(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, leftQuantityOnWarehouse); err != nil {
					logger.DebugKV(ctx, "res", "layer", layer, "err", err.Error())
					ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: err.Error()})
					break
					/*
						ADD ROLLBACK
					*/
				}

				_, err = s.repo.CreateReservation(ctx, productsByWarehouse.WarehouseId, productsByWarehouse.ProductId, product.Quantity-quantityForReservation, reservationId)
				if err != nil {
					logger.DebugKV(ctx, "res", "layer", layer, "err", err.Error())
					ProductsInfo = append(ProductsInfo, model.ReserveProductResp{Code: product.Code, Status: err.Error()})
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

	logger.DebugKV(ctx, "res", "layer", layer, "ProductsInfo", ProductsInfo)
	return model.ReserveProductsResp{
		ReservationId:           reservationId,
		ReservationProductsInfo: ProductsInfo,
	}
}

func (s *Service) ReleaseProducts(ctx context.Context, products []model.ReleaseProductReq) model.ReleaseProductsResp {
	var (
		productsInfo []model.ReleaseProductResp
		wg           sync.WaitGroup
	)

	wg.Add(len(products))
	for _, product := range products {
		go func() {
			defer wg.Done()

			if product.Quantity <= 0 {
				productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: model.ErrInvalidInput.Error()})
				return
			}

			quantityProductsInReservation, err := s.repo.GetTotalQuantityOfReservation(ctx, product.ReservationId, product.Code)
			if err != nil {
				logger.DebugKV(ctx, "res", "quantityProductsInReservation", quantityProductsInReservation)
				productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: model.ErrInvalidInput.Error()})
				return
			}

			if quantityProductsInReservation < product.Quantity {
				logger.DebugKV(ctx, "res", "quantityProductsInReservation", quantityProductsInReservation)
				productsInfo = append(
					productsInfo,
					model.ReleaseProductResp{
						ReservationId: product.ReservationId,
						Code:          product.Code,
						Status:        model.ErrInvalidReservationQuantity.Error(),
					})
				return
			}

			logger.DebugKV(ctx, "res", "product.ReservationId", product.ReservationId, "product.Code", product.Code)
			productsByWarehousesInReservation, err := s.repo.GetProductsByReservationByIdAndCode(ctx, product.ReservationId, product.Code)
			if err != nil {
				logger.DebugKV(ctx, "res", "productsByWarehousesInReservation", productsByWarehousesInReservation)
				productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: model.ErrInvalidInput.Error()})
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

				logger.DebugKV(ctx, "res", "productsByWarehouseInResevation.ID", productsByWarehouseInResevation.ID, "remainInReservation", remainInReservation)
				if err := s.repo.UpdateReservationQuantity(ctx, productsByWarehouseInResevation.ID, remainInReservation); err != nil {
					productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: err.Error()})
					break
				}

				logger.DebugKV(ctx, "res", "productsByWarehouseInResevation.WarehouseId", productsByWarehouseInResevation.WarehouseId, "productsByWarehouseInResevation.ProductId", productsByWarehouseInResevation.ProductId, "addToWarehouse", addToWarehouse)
				if err := s.repo.UpdateWarehouseQuantityWithAdd(ctx, productsByWarehouseInResevation.WarehouseId, productsByWarehouseInResevation.ProductId, addToWarehouse); err != nil {
					productsInfo = append(productsInfo, model.ReleaseProductResp{ReservationId: product.ReservationId, Code: product.Code, Status: err.Error()})
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

	return model.ReleaseProductsResp{ReleaseProductsInfo: productsInfo}
}

func (s *Service) GetProductsByWarehouse(ctx context.Context, id int) ([]model.Product, error) {
	layer := "service.GetProductByWH"
	products, err := s.repo.GetProductsByWarehouseId(ctx, id)
	if err != nil {
		return nil, err
	}
	logger.DebugKV(ctx, "res", "layer", layer, "products", products)

	return products, nil
}
