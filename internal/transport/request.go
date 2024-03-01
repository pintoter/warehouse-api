package transport

import "github.com/pintoter/warehouse-api/internal/model"

type reserveProductsReq struct {
	Products []model.ReserveProductReq `json:"products"`
}

type releaseProductsReq struct {
	Products []model.ReleaseProductReq `json:"products"`
}

type showProductsReq struct {
	WarehouseId int `json:"warehouse_id"`
}
