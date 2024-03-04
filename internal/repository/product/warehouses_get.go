package product

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	repoModel "github.com/pintoter/warehouse-api/internal/repository/model"
	"github.com/pintoter/warehouse-api/internal/service/model"
)

func getProductsByWarehouseIdBuilder(id int) (string, []interface{}, error) {
	builder := sq.Select("p.id", "p.name", "p.size", "p.code", "wp.quantity").
		From(warehouseProduct + " wp").
		Join(product + " p ON p.id = wp.product_id").
		Where(sq.Eq{"wp.warehouse_id": id}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) GetProductsByWarehouseId(ctx context.Context, id int) ([]model.Product, error) {
	query, args, err := getProductsByWarehouseIdBuilder(id)
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

	var products []model.Product
	for rows.Next() {
		var product model.Product

		err = rows.Scan(&product.ID, &product.Name, &product.Size, &product.Code, &product.Quantity)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}

func getTotalQuantityOfProductsBuilder(code string) (string, []interface{}, error) {
	builder := sq.Select("SUM(quantity)").
		From(warehouseProduct + " wp").
		Join(product + " p ON p.id = wp.product_id").
		Join(warehouse + " w ON w.id = wp.warehouse_id").
		Where(sq.Eq{"p.code": code, "w.availability": true}).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) GetTotalQuantityOfProducts(ctx context.Context, code string) (int, error) {
	query, args, err := getTotalQuantityOfProductsBuilder(code)
	if err != nil {
		return 0, err
	}

	var count int
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getProductsByWarehousesByCodeBuilder(code string) (string, []interface{}, error) {
	builder := sq.Select("wp.warehouse_id, wp.product_id, wp.quantity").
		From(warehouseProduct + " wp").
		Join(product + " p ON p.id = wp.product_id").
		Join(warehouse + " w ON w.id = wp.warehouse_id").
		Where(sq.Eq{"p.code": code, "w.availability": true}).
		OrderBy("wp.quantity DESC").
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}

func (r *repo) GetProductsByWarehousesByCode(ctx context.Context, code string) ([]repoModel.ProductsOnActiveWarehouse, error) {
	query, args, err := getProductsByWarehousesByCodeBuilder(code)
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

	var ProductsByWHs []repoModel.ProductsOnActiveWarehouse
	for rows.Next() {
		var ProductsByWH repoModel.ProductsOnActiveWarehouse

		err = rows.Scan(&ProductsByWH.WarehouseId, &ProductsByWH.ProductId, &ProductsByWH.Quantity)
		if err != nil {
			return nil, err
		}

		ProductsByWHs = append(ProductsByWHs, ProductsByWH)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return ProductsByWHs, nil
}

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
