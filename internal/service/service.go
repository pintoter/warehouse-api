package service

import (
	"net/http"

	"github.com/pintoter/warehouse-api/internal/service/model"
)

type ProductService interface {
	ReserveProducts(r *http.Request, args *model.ReserveProductsReq, reply *model.ReserveProductsResp) error
	ReleaseProducts(r *http.Request, args *model.ReleaseProductsReq, reply *model.ReleaseProductsResp) error
	GetProductsByWarehouse(r *http.Request, args *model.ShowProductsReq, reply *[]model.Product) error
}
