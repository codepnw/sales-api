package catrepositories

import (
	"context"
	"time"

	"github.com/codepnw/sales-api/modules/categories"
	"github.com/jmoiron/sqlx"
)

type ICategoryRepo interface {
	CreateCategory(category *categories.Category) (*categories.Category, error)
	GetOneCategory(categoryId int) (*categories.Category, error)
	GetAllCategories() ([]*categories.Category, error)
	UpdateCategory(category *categories.Category) (*categories.Category, error)
	DeleteCategory(categoryId int) error
}

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) ICategoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CreateCategory(category *categories.Category) (*categories.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	query := `
		INSERT INTO "categories" ("title", "desc")
		VALUES ($1, $2)
		RETURNING "category_id";
	`
	err := r.db.QueryRowContext(ctx, query, category.Title, category.Desc).Scan(&category.CategoryId)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *categoryRepo) GetOneCategory(categoryId int) (*categories.Category, error) {
	category := categories.Category{}

	query := `
		SELECT * FROM "categories"
		WHERE "category_id" = $1
		LIMIT 1;
	`
	err := r.db.Get(&category, query, categoryId)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepo) GetAllCategories() ([]*categories.Category, error) {
	categories := []*categories.Category{}

	query := `SELECT * FROM "categories";`
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepo) UpdateCategory(category *categories.Category) (*categories.Category, error) {
	query := `
		UPDATE "categories"
		SET
			"title" = COALESCE(NULLIF($1, ''), "title"),
			"desc" = COALESCE(NULLIF($2, ''), "desc")
		WHERE "category_id" = $3;
	`
	_, err := r.db.ExecContext(context.Background(), query, category.Title, category.Desc, category.CategoryId)
	if err != nil {
		return nil, err
	}

	c, err := r.GetOneCategory(category.CategoryId)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *categoryRepo) DeleteCategory(categoryId int) error {
	query := `DELETE FROM "categories" WHERE "category_id" = $1;`

	result, err := r.db.ExecContext(context.Background(), query, categoryId)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}
