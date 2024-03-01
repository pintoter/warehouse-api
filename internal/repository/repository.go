package repository

import (
	"context"

	"github.com/pintoter/warehouse-api/internal/model"
	repoModel "github.com/pintoter/warehouse-api/internal/repository/model"
)

type ProductRepository interface {
	GetProductsByWarehouseId(ctx context.Context, id int) ([]model.Product, error)
	GetTotalQuantityOfProducts(ctx context.Context, code string) (int, error)
	GetProductsByWarehousesByCode(ctx context.Context, code string) ([]repoModel.ProductsOnActiveWarehouse, error)
}

type WarehouseRepository interface {
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
	ProductRepository
	WarehouseRepository
	ReservationRepository
}
