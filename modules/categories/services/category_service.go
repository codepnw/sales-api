package catservices

import (
	"database/sql"
	"fmt"

	"github.com/codepnw/sales-api/modules/categories"
	catrepositories "github.com/codepnw/sales-api/modules/categories/repositories"
	"github.com/codepnw/sales-api/pkg/logs"
)

type ICategoryService interface {
	CreateCategory(request *categories.Category) (*categories.Category, error)
	GetOneCategory(categoryId int) (*categories.Category, error)
	GetAllCategories() ([]*categories.Category, error)
	UpdateCategory(categoryId int, category *categories.Category) (*categories.Category, error)
	DeleteCategory(categoryId int) error
}

type categoryService struct {
	repo catrepositories.ICategoryRepo
}

func NewCategoryService(repo catrepositories.ICategoryRepo) ICategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) CreateCategory(request *categories.Category) (*categories.Category, error) {
	category := categories.Category{
		Title: request.Title,
		Desc:  request.Desc,
	}

	result, err := s.repo.CreateCategory(&category)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return result, nil
}

func (s *categoryService) GetOneCategory(categoryId int) (*categories.Category, error) {
	result, err := s.repo.GetOneCategory(categoryId)
	if err != nil {
		logs.Error(err)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category id not found")
		}
		return nil, fmt.Errorf("failed get category")
	}
	return result, nil
}

func (s *categoryService) GetAllCategories() ([]*categories.Category, error) {
	result, err := s.repo.GetAllCategories()
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("failed get category")
	}
	return result, nil
}

func (s *categoryService) UpdateCategory(categoryId int, category *categories.Category) (*categories.Category, error) {
	request := categories.Category{
		CategoryId: categoryId,
		Title:      category.Title,
		Desc:       category.Desc,
	}

	result, err := s.repo.UpdateCategory(&request)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("failed update category")
	}

	return result, nil
}

func (s *categoryService) DeleteCategory(categoryId int) error {
	if err := s.repo.DeleteCategory(categoryId); err != nil {
		logs.Error(err)
		return fmt.Errorf("failed delete category")
	}
	return nil
}
