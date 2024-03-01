package product

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pintoter/warehouse-api/internal/service/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetProductsByWarehouseId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		id int
	}

	type mockBehavior func(args args)

	id := 1
	products := []model.Product{
		{
			ID:       1,
			Name:     "Adidas",
			Size:     "L",
			Code:     "1234",
			Quantity: 3,
		},
		{
			ID:       2,
			Name:     "Nike",
			Size:     "L",
			Code:     "12",
			Quantity: 2,
		},
		{
			ID:       3,
			Name:     "Puma",
			Size:     "L",
			Code:     "123",
			Quantity: 1,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantProducts []model.Product
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "name", "size", "code", "quantity"}).
					AddRow(
						products[0].ID,
						products[0].Name,
						products[0].Size,
						products[0].Code,
						products[0].Quantity,
					).AddRow(
					products[1].ID,
					products[1].Name,
					products[1].Size,
					products[1].Code,
					products[1].Quantity,
				).AddRow(
					products[2].ID,
					products[2].Name,
					products[2].Size,
					products[2].Code,
					products[2].Quantity,
				)

				expectedQuery := `SELECT p.id, p.name, p.size, p.code, wp.quantity FROM warehouse_product wp 
					JOIN product p ON p.id = wp.product_id
					WHERE wp.warehouse_id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnRows(rows)
			},
			args:         args{id: id},
			wantProducts: products,
		},
		{
			name: "Failed",
			mockBehavior: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "name", "size", "code", "quantity"})

				expectedQuery := `SELECT p.id, p.name, p.size, p.code, wp.quantity FROM warehouse_product wp JOIN product p ON p.id = wp.product_id WHERE wp.warehouse_id = $1`
				mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(args.id).WillReturnError(errors.New("some error")).WillReturnRows(rows)
			},
			args:         args{id: id},
			wantProducts: []model.Product{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotProducts, err := r.GetProductsByWarehouseId(context.Background(), tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantProducts, gotProducts)
			}
		})
		assert.NoError(t, mock.ExpectationsWereMet())
	}
}
