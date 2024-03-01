package product

import (
	"database/sql"

	"github.com/pintoter/warehouse-api/internal/repository"
)

const (
	product          = "product"
	warehouse        = "warehouse"
	warehouseProduct = "warehouse_product"
	reservation      = "reservation"
)

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) repository.Repository {
	return &repo{
		db: db,
	}
}
