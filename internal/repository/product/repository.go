package product

import (
	"github.com/jmoiron/sqlx"
	"github.com/pintoter/warehouse-api/internal/repository"
)

const (
	product          = "product"
	warehouse        = "warehouse"
	warehouseProduct = "warehouse_product"
	reservation      = "reservation"
)

type repo struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.Repository {
	return &repo{
		db: db,
	}
}
