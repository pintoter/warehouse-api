package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/pintoter/warehouse-api/pkg/logger"
)

func updateBuilder(warehouseId, productId, quantity int) (string, []interface{}, error) {
	builder := sq.Update(warehouseProduct).
		Where(sq.Eq{"warehouse_id": warehouseId, "product_id": productId}).
		Set("quantity", quantity).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) UpdateWarehouseQuantity(ctx context.Context, warehouseId, productId, quantity int) error {
	query, args, err := updateBuilder(warehouseId, productId, quantity)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// func updateWarehouseQuantityWithAddBuilder(warehouseId, productId, quantity int) (string, []interface{}, error) {
// 	builder := sq.Update(warehouseProduct).
// 		Where(sq.Eq{"warehouse_id": warehouseId, "product_id": productId}).
// 		Set(sq.Eq{"quantity": quantity}).
// 		Set("quantity", quantity).
// 		PlaceholderFormat(sq.Dollar)

// 	return builder.ToSql()
// }

func (r *repo) UpdateWarehouseQuantityWithAdd(ctx context.Context, warehouseId, productId, quantity int) error {
	// query, args, err := updateWarehouseQuantityWithAddBuilder(warehouseId, productId, quantity)
	// if err != nil {
	// 	return err
	// }
	logger.DebugKV(ctx, "UpdateWarehouseQuantityWithAdd", "warehouseId", warehouseId, "quantity", quantity)
	query := "UPDATE warehouse_product SET quantity = quantity + $1 WHERE product_id = $2 AND warehouse_id = $3"
	_, err := r.db.ExecContext(ctx, query, quantity, productId, warehouseId)
	if err != nil {
		return err
	}

	return nil
}
