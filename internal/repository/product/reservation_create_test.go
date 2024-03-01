package product

import (
	"context"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateReservation(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		warehouseId   int
		productId     int
		quantity      int
		reservationId string
	}

	type mockBehavior func(args args)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantId       int
		wantErr      bool
	}{
		{
			name: "Success",
			mockBehavior: func(args args) {
				expectedQueryInReservation := "INSERT INTO reservation (reservation_id,warehouse_id,product_id,quantity) VALUES ($1,$2,$3,$4) RETURNING id"
				mock.ExpectQuery(regexp.QuoteMeta(expectedQueryInReservation)).
					WithArgs(
						args.reservationId,
						args.warehouseId,
						args.productId,
						args.quantity,
					).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			args: args{
				warehouseId:   1,
				productId:     1,
				quantity:      2,
				reservationId: "1337",
			},
			wantId: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)

			gotId, err := r.CreateReservation(context.Background(), tt.args.warehouseId, tt.args.productId, tt.args.quantity, tt.args.reservationId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantId, gotId)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
