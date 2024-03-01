package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func createReservationBuilder(warehouseId, productId, quantity int, reservationId string) (string, []interface{}, error) {
	builder := sq.Insert(reservation).
		Columns("reservation_id", "warehouse_id", "product_id", "quantity").
		Values(reservationId, warehouseId, productId, quantity).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) CreateReservation(ctx context.Context, warehouseId, productId, quantity int, reservationId string) (int, error) {
	query, args, err := createReservationBuilder(warehouseId, productId, quantity, reservationId)
	if err != nil {
		return 0, err
	}

	var id int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
