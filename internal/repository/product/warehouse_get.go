package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

func getWarehouseAvailabilityBuilder(id int) (string, []interface{}, error) {
	builder := sq.Select("availability").
		From(warehouse).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) GetWarehouseAvailabilityById(ctx context.Context, id int) (bool, error) {
	query, args, err := getWarehouseAvailabilityBuilder(id)
	if err != nil {
		return false, err
	}

	var isAvailable bool
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&isAvailable)
	if err != nil {
		return false, err
	}

	return isAvailable, nil
}
