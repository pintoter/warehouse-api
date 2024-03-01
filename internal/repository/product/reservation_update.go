package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func updateReservationQuantityBuilder(id, quantity int) (string, []interface{}, error) {
	builder := sq.Update(reservation).
		Where(sq.Eq{"id": id}).
		Set("quantity", quantity).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) UpdateReservationQuantity(ctx context.Context, id, quantity int) error {
	query, args, err := updateReservationQuantityBuilder(id, quantity)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
