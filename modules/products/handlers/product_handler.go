package prodhandlers

import (
	"net/http"
	"strings"

	"github.com/codepnw/sales-api/modules/products"
	prodservices "github.com/codepnw/sales-api/modules/products/services"
	"github.com/codepnw/sales-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type productHandler struct {
	service prodservices.IProductService
}

func NewProductHandler(service prodservices.IProductService) *productHandler {
	return &productHandler{service: service}
}

type productErr string

const (
	createError productErr = "products-001"
	getOneError productErr = "products-002"
	getAllError productErr = "products-003"
	updateError productErr = "products-004"
	deleteError productErr = "products-005"
)

func (h *productHandler) CreateProduct(c *gin.Context) {
	request := products.ProductRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewResponse(c).Error(
			http.StatusBadRequest,
			string(createError),
			err.Error(),
		)
		return
	}

	product, err := h.service.CreateProduct(&request)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(createError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusCreated, product)
}

func (h *productHandler) GetProducts(c *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(getAllError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, products)
}

func (h *productHandler) GetProduct(c *gin.Context) {
	id := c.Param("productId")

	product, err := h.service.GetProduct(id)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(getOneError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, product)
}

func (h *productHandler) UpdateProduct(c *gin.Context) {
	id := strings.Trim(c.Param("productId"), " ")
	request := products.ProductRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewResponse(c).Error(
			http.StatusBadRequest,
			string(updateError),
			err.Error(),
		)
		return
	}

	p, err := h.service.UpdateProduct(id, &request)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(updateError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, &p)
}

func (h *productHandler) DeleteProduct(c *gin.Context) {
	id := strings.Trim(c.Param("productId"), " ")

	err := h.service.DeleteProduct(id)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(deleteError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusNoContent, nil)
}
