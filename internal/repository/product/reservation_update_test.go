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

func TestUpdateReservationQuantity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	r := NewRepository(sqlxDB)

	type args struct {
		id       int
		quantity int
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
				id:       1,
				quantity: 4,
			},
			mockBehavior: func(args args) {
				expectedQuery := "UPDATE reservation SET quantity = $1 WHERE id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.quantity, args.id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			name: "Failed",
			args: args{
				id:       100,
				quantity: 1,
			},
			mockBehavior: func(args args) {
				expectedQuery := "UPDATE reservation SET quantity = $1 WHERE id = $2"
				mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
					WithArgs(args.quantity, args.id).
					WillReturnError(errors.New("any error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior(tt.args)
			err := r.UpdateReservationQuantity(context.Background(), tt.args.id, tt.args.quantity)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
