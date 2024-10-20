package prodrepositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/codepnw/sales-api/modules/products"
	"github.com/jmoiron/sqlx"
)

type IProductRepo interface {
	CreateProduct(product *products.Product) (*products.Product, error)
	GetProducts() ([]*products.Product, error)
	GetProduct(productID string) (*products.Product, error)
	UpdateProduct(product *products.Product) (*products.Product, error)
	DeleteProduct(productID string) error
}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) IProductRepo {
	return &productRepo{db: db}
}

func (r *productRepo) CreateProduct(product *products.Product) (*products.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
		INSERT INTO products ("name", "desc", "price", "discount", "stock", "category_id", "created_at", "updated_at")
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING "product_id";
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		product.Name,
		product.Desc,
		product.Price,
		product.Discount,
		product.Stock,
		product.CategoryID,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ProductID)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepo) GetProducts() ([]*products.Product, error) {
	prods := make([]*products.Product, 0)

	query := `
		SELECT "product_id", "name", "desc", "price", "discount", "stock", "category_id", "created_at", "updated_at"
		FROM "products";
	`
	err := r.db.Select(&prods, query)
	if err != nil {
		return nil, err
	}

	return prods, nil
}

func (r *productRepo) GetProduct(productID string) (*products.Product, error) {
	prod := products.Product{}

	query := `
		SELECT "product_id", "name", "desc", "price", "discount", "stock", "category_id", "created_at", "updated_at"
		FROM "products"
		WHERE "product_id" = $1
		LIMIT 1;
	`
	err := r.db.Get(&prod, query, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product_id not found")
		}
		return nil, err
	}

	return &prod, nil
}

func (r *productRepo) UpdateProduct(product *products.Product) (*products.Product, error) {
	query := `
		UPDATE "products"
		SET 
			"name" = COALESCE(NULLIF($1, ''), "name"),
			"desc" = COALESCE(NULLIF($2, ''), "desc"),
			"price" = COALESCE(NULLIF($3, 0.0), "price"),
			"discount" = COALESCE(NULLIF($4, 0), "discount"),
			"stock" = COALESCE(NULLIF($5, 0), "stock"),
			"updated_at" = $6
		WHERE "product_id" = $7;
	`
	_, err := r.db.ExecContext(
		context.Background(),
		query,
		product.Name,
		product.Desc,
		product.Price,
		product.Discount,
		product.Stock,
		product.UpdatedAt,
		product.ProductID,
	)
	if err != nil {
		return nil, err
	}

	p, err := r.GetProduct(product.ProductID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *productRepo) DeleteProduct(productID string) error {
	query := `DELETE FROM "products" WHERE "product_id" = $1;`

	result, err := r.db.ExecContext(context.Background(), query, productID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
