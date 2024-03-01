package repository

import (
	"context"

	repoModel "github.com/pintoter/warehouse-api/internal/repository/model"
	"github.com/pintoter/warehouse-api/internal/service/model"
)

type WarehousesRepository interface {
	GetProductsByWarehouseId(ctx context.Context, id int) ([]model.Product, error)
	GetProductsByWarehousesByCode(ctx context.Context, code string) ([]repoModel.ProductsOnActiveWarehouse, error)
	GetTotalQuantityOfProducts(ctx context.Context, code string) (int, error)
	GetWarehouseAvailabilityById(ctx context.Context, warehouseId int) (bool, error)
	UpdateWarehouseQuantity(ctx context.Context, warehouseId, productId, quantity int) error
	UpdateWarehouseQuantityWithAdd(ctx context.Context, warehouseId, productId, quantity int) error
}

type ReservationRepository interface {
	CreateReservation(ctx context.Context, warehouseId, productId, quantity int, reservationId string) (int, error)
	GetTotalQuantityOfReservation(ctx context.Context, reservationId string, productCode string) (int, error)
	GetProductsByReservationByIdAndCode(ctx context.Context, reservationId, code string) ([]repoModel.ProductsInReservation, error)
	UpdateReservationQuantity(ctx context.Context, id, quantity int) error
}

type Repository interface {
	WarehousesRepository
	ReservationRepository
}
