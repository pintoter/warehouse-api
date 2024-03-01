package service

import (
	"context"

	"github.com/pintoter/warehouse-api/internal/model"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type ProductService interface {
	ReserveProducts(ctx context.Context, products []model.ReserveProductReq) model.ReserveProductsResp
	ReleaseProducts(ctx context.Context, products []model.ReleaseProductReq) model.ReleaseProductsResp
	GetProductsByWarehouse(ctx context.Context, id int) ([]model.Product, error)
}
