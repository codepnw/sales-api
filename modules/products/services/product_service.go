package prodservices

import (
	"database/sql"
	"fmt"

	"github.com/codepnw/sales-api/modules/products"
	prodrepositories "github.com/codepnw/sales-api/modules/products/repositories"
	"github.com/codepnw/sales-api/pkg/logs"
	"github.com/codepnw/sales-api/pkg/utils"
)

type IProductService interface {
	CreateProduct(prod *products.ProductRequest) (*products.Product, error)
	GetProducts() ([]*products.Product, error)
	GetProduct(productId string) (*products.Product, error)
	UpdateProduct(productId string, req *products.ProductRequest) (*products.Product, error)
	DeleteProduct(productId string) error
}

type productService struct {
	repository prodrepositories.IProductRepo
}

func NewProductService(repository prodrepositories.IProductRepo) IProductService {
	return &productService{repository: repository}
}

func (s *productService) CreateProduct(req *products.ProductRequest) (*products.Product, error) {
	var stock uint
	if req.Stock == 0 {
		stock = 1
	}

	if req.Price <= 0 {
		logs.Error("price is zero")
		return nil, fmt.Errorf("price is zero")
	}

	product := products.Product{
		Name:       req.Name,
		Desc:       req.Desc,
		Price:      req.Price,
		Discount:   req.Discount,
		Stock:      stock,
		CategoryID: req.CategoryID,
		CreatedAt:  utils.LocalTime(),
		UpdatedAt:  utils.LocalTime(),
	}

	p, err := s.repository.CreateProduct(&product)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("failed create product")
	}

	return p, nil
}

func (s *productService) GetProducts() ([]*products.Product, error) {
	p, err := s.repository.GetProducts()
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("failed get products")
	}

	return p, nil
}

func (s *productService) GetProduct(productId string) (*products.Product, error) {
	p, err := s.repository.GetProduct(productId)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("failed get product")
	}

	return p, nil
}

func (s *productService) UpdateProduct(productId string, req *products.ProductRequest) (*products.Product, error) {
	product := products.Product{
		ProductID:  productId,
		Name:       req.Name,
		Desc:       req.Desc,
		Price:      req.Price,
		Discount:   req.Discount,
		Stock:      req.Stock,
		CategoryID: req.CategoryID,
		UpdatedAt:  utils.LocalTime(),
	}

	p, err := s.repository.UpdateProduct(&product)
	if err != nil {
		logs.Error(err)
		return nil, err
	}

	return p, nil
}

func (s *productService) DeleteProduct(productId string) error {
	err := s.repository.DeleteProduct(productId)
	if err != nil {
		logs.Error(err)
		if err == sql.ErrNoRows {
			return fmt.Errorf("product id not found")
		}
		return err
	}
	return nil
}
