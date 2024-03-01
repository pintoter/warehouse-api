package product

import (
	"context"
	"errors"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUpdateWarehouseQuantity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		warehouseId int
		productId   int
		quantity    int
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Success",
			args: args{
				warehouseId: 1,
				productId:   4,
				quantity:    8,
			},
			mockBehavior: func(args args) {
				expectedQuery := "UPDATE warehouse_product SET quantity = $1 WHERE product_id = $2 AND warehouse_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.quantity, args.productId, args.warehouseId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Failed",
			args: args{
				warehouseId: 1,
				productId:   4,
				quantity:    8,
			},
			mockBehavior: func(args args) {
				expectedQuery := "UPDATE warehouse_product SET quantity = $1 WHERE product_id = $2 AND warehouse_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.quantity, args.productId, args.warehouseId).
					WillReturnError(errors.New("any error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := r.UpdateWarehouseQuantity(context.Background(), tt.args.warehouseId, tt.args.productId, tt.args.quantity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func UpdateWarehouseQuantityWithAdd(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		warehouseId int
		productId   int
		quantity    int
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "Success",
			args: args{
				warehouseId: 1,
				productId:   4,
				quantity:    8,
			},
			mockBehavior: func(args args) {
				expectedQuery := "UPDATE warehouse_product SET quantity = quantity + $1 WHERE product_id = $2 AND warehouse_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.quantity, args.warehouseId, args.productId).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Failed",
			args: args{
				warehouseId: 1,
				productId:   4,
				quantity:    8,
			},
			mockBehavior: func(args args) {
				expectedQuery := "UPDATE warehouse_product SET quantity = quantity + $1 WHERE product_id = $2 AND warehouse_id = $3"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.quantity, args.warehouseId, args.productId).
					WillReturnError(errors.New("any error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := r.UpdateWarehouseQuantityWithAdd(context.Background(), tt.args.warehouseId, tt.args.productId, tt.args.quantity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
