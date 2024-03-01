package product

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	repoModel "github.com/pintoter/warehouse-api/internal/repository/model"
	"github.com/stretchr/testify/assert"
)

func TestGetTotalQuantityOfReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		reservationId string
		productCode   string
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantQuantity int
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				expectedExecInReservation := "SELECT SUM(r.quantity) FROM reservation r JOIN product p ON p.id = r.product_id WHERE p.code = $1 AND r.reservation_id = $2"
				mock.ExpectQuery(regexp.QuoteMeta(expectedExecInReservation)).
					WithArgs(
						args.productCode,
						args.reservationId,
					).WillReturnRows(sqlmock.NewRows([]string{"SUM(quantity)"}).AddRow(50))
			},
			args: args{
				reservationId: "1",
				productCode:   "1337",
			},
			wantQuantity: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotQuantity, err := r.GetTotalQuantityOfReservation(context.Background(), tt.args.reservationId, tt.args.productCode)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantQuantity, gotQuantity)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetProductsByReservationByIdAndCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		reservationId string
		productCode   string
	}

	type mockBehavior func(args args)

	products := []repoModel.ProductsInReservation{
		{
			ID:          1,
			WarehouseId: 1,
			ProductId:   1,
			Quantity:    5,
		},
		{
			ID:          2,
			WarehouseId: 2,
			ProductId:   1,
			Quantity:    3,
		},
		{
			ID:          3,
			WarehouseId: 3,
			ProductId:   1,
			Quantity:    1,
		},
	}

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantProducts []repoModel.ProductsInReservation
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				expectedExecInReservation := `SELECT r.id, r.warehouse_id, r.product_id, r.quantity 
				FROM reservation r 
				JOIN product p ON p.id = r.product_id
				JOIN warehouse w ON w.id = r.warehouse_id
				WHERE p.code = $1 AND r.reservation_id = $2
				ORDER BY r.quantity DESC`
				mock.ExpectQuery(regexp.QuoteMeta(expectedExecInReservation)).
					WithArgs(
						args.productCode,
						args.reservationId,
					).WillReturnRows(
					sqlmock.NewRows(
						[]string{"r.id", "r.warehouse_id", "r.product_id", "r.quantity"},
					).AddRow(products[0].ID, products[0].WarehouseId, products[0].ProductId, products[0].Quantity).
						AddRow(products[1].ID, products[1].WarehouseId, products[1].ProductId, products[1].Quantity).
						AddRow(products[2].ID, products[2].WarehouseId, products[2].ProductId, products[2].Quantity))
			},
			args: args{
				reservationId: "1",
				productCode:   "1337",
			},
			wantProducts: products,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotProducts, err := r.GetProductsByReservationByIdAndCode(context.Background(), tt.args.reservationId, tt.args.productCode)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantProducts, gotProducts)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
