package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	repoModel "github.com/pintoter/warehouse-api/internal/repository/model"
	"github.com/pkg/errors"
)

func getTotalQuantityOfReservationBuilder(reservationId, productCode string) (string, []interface{}, error) {
	builder := sq.Select("SUM(r.quantity)").
		From(reservation + " r").
		Join(product + " p ON p.id = r.product_id").
		Where(sq.Eq{"r.reservation_id": reservationId, "p.code": productCode}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) GetTotalQuantityOfReservation(ctx context.Context, reservationId string, productCode string) (int, error) {
	query, args, err := getTotalQuantityOfReservationBuilder(reservationId, productCode)
	if err != nil {
		return 0, err
	}

	var totalQuantity int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&totalQuantity)
	if err != nil {
		return 0, err
	}

	return totalQuantity, nil
}

func getProductsByReservationByCodeBuilder(reservationId, code string) (string, []interface{}, error) {
	builder := sq.Select("r.id, r.warehouse_id, r.product_id, r.quantity").
		From(reservation + " r").
		Join(product + " p ON p.id = r.product_id").
		Join(warehouse + " w ON w.id = r.warehouse_id").
		Where(sq.Eq{"p.code": code, "r.reservation_id": reservationId}).
		OrderBy("r.quantity DESC").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) GetProductsByReservationByIdAndCode(ctx context.Context, reservationId, code string) ([]repoModel.ProductsInReservation, error) {
	query, args, err := getProductsByReservationByCodeBuilder(reservationId, code)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var productsInReservation []repoModel.ProductsInReservation
	for rows.Next() {
		var ProductInReservation repoModel.ProductsInReservation

		err = rows.Scan(&ProductInReservation.ID, &ProductInReservation.WarehouseId, &ProductInReservation.ProductId, &ProductInReservation.Quantity)
		if err != nil {
			return nil, errors.Wrap(err, "GetProductsInReservation.rows.Scan")
		}

		productsInReservation = append(productsInReservation, ProductInReservation)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return productsInReservation, nil
}
