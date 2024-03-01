package model

type ProductsOnActiveWarehouse struct {
	WarehouseId int
	ProductId   int
	Quantity    int
}

type ProductsInReservation struct {
	ID          int
	WarehouseId int
	ProductId   int
	Quantity    int
}
