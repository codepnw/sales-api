package cathandlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/codepnw/sales-api/modules/categories"
	catservices "github.com/codepnw/sales-api/modules/categories/services"
	"github.com/codepnw/sales-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type categoryHandler struct {
	service catservices.ICategoryService
}

func NewCategoryHandler(service catservices.ICategoryService) *categoryHandler {
	return &categoryHandler{service: service}
}

type categoryErr string

const (
	createError categoryErr = "category-001"
	getOneError categoryErr = "category-002"
	getAllError categoryErr = "category-003"
	updateError categoryErr = "category-004"
	deleteError categoryErr = "category-005"
)

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	request := categories.Category{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewResponse(c).Error(
			http.StatusBadRequest,
			string(createError),
			err.Error(),
		)
		return
	}

	category, err := h.service.CreateCategory(&request)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(createError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusCreated, category)
}

func (h *categoryHandler) GetOneCategory(c *gin.Context) {
	idStr := strings.Trim(c.Param("categoryId"), " ")
	id, _ := strconv.Atoi(idStr)

	category, err := h.service.GetOneCategory(id)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(getOneError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, category)
}

func (h *categoryHandler) GetAllCategory(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(getAllError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, categories)
}

func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	idStr := strings.Trim(c.Param("categoryId"), " ")
	id, _ := strconv.Atoi(idStr)
	request := categories.Category{}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewResponse(c).Error(
			http.StatusBadRequest,
			string(updateError),
			err.Error(),
		)
		return
	}

	category, err := h.service.UpdateCategory(id, &request)
	if err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(updateError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusOK, category)
}

func (h *categoryHandler) DeleteCategory(c *gin.Context) {
	idStr := strings.Trim(c.Param("categoryId"), " ")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.DeleteCategory(id); err != nil {
		utils.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(deleteError),
			err.Error(),
		)
		return
	}

	utils.NewResponse(c).Success(http.StatusNoContent, nil)
}
