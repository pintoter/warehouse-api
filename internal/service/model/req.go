package model

type ReserveProductReq struct {
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}

type ReserveProductsReq struct {
	Products []ReserveProductReq `json:"products"`
}

type ReserveProductResp struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

type ReserveProductsResp struct {
	ReservationId           string               `json:"reservation_id"`
	ReservationProductsInfo []ReserveProductResp `json:"reservation_products_info"`
}

type ReleaseProductReq struct {
	ReservationId string `json:"reservation_id"`
	Code          string `json:"code"`
	Quantity      int    `json:"quantity"`
}

type ReleaseProductsReq struct {
	Products []ReleaseProductReq `json:"products"`
}

type ReleaseProductResp struct {
	ReservationId string `json:"reservation_id"`
	Code          string `json:"code"`
	Status        string `json:"status"`
}

type ReleaseProductsResp struct {
	ReleaseProductsInfo []ReleaseProductResp `json:"release_products_info"`
}

type ShowProductsReq struct {
	WarehouseId int `json:"warehouse_id"`
}
