package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/codepnw/sales-api/entities"
	"github.com/codepnw/sales-api/pkg/logs"
	"github.com/codepnw/sales-api/pkg/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tableCategories string = entities.TableCategories()

type ErrCategory string

const (
	createCategoryErr ErrCategory = "category-001"
	getOneCategoryErr ErrCategory = "category-002"
	getAllCategoryErr ErrCategory = "category-003"
	updateCategoryErr ErrCategory = "category-004"
	deleteCategoryErr ErrCategory = "category-005"
)

type ICategory interface {
	CreateCategory(c *gin.Context)
	GetCategory(c *gin.Context)
	GetCategories(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
}

type category struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) ICategory {
	return &category{db: db}
}

func (h *category) CreateCategory(c *gin.Context) {
	req := entities.Category{}

	if err := c.ShouldBindJSON(&req); err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(createCategoryErr),
			"category format invalid",
		)
		return
	}

	if req.Name == "" {
		logs.Error("name is empty")
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(createCategoryErr),
			"category name is required",
		)
		return
	}

	category := entities.Category{
		Name: req.Name,
	}

	// database create
	err := h.db.Table(tableCategories).Create(&category).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(createCategoryErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusCreated, category)
}

func (h *category) GetCategory(c *gin.Context) {
	id := c.Param("categoryId")
	var category entities.Category

	// database select
	err := h.db.Table(tableCategories).Select("id, name").Where("id = ?", id).First(&category).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(getOneCategoryErr),
			err.Error(),
		)
		return
	}

	data := entities.CategoryResponse{
		Id:   category.Id,
		Name: category.Name,
	}

	responses.NewResponse(c).Success(http.StatusOK, data)
}

func (h *category) GetCategories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var categories []entities.CategoryResponse

	if limit == 0 {
		limit = 10
	}

	// database select
	err := h.db.Table(tableCategories).Select("id, name").Limit(limit).Offset(skip).Find(&categories).Count(&count).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(getAllCategoryErr),
			err.Error(),
		)
		return
	}

	response := entities.CategoryListResponse{
		Categories: categories,
		Meta: &entities.QueryMeta{
			Total: int(count),
			Limit: limit,
			Skip:  skip,
		},
	}

	responses.NewResponse(c).Success(http.StatusOK, response)
}

func (h *category) UpdateCategory(c *gin.Context) {
	id := c.Param("categoryId")
	var category entities.CategoryUpdate

	// database find id
	err := h.db.Table(tableCategories).Find(&category, "id = ?", id).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(updateCategoryErr),
			err.Error(),
		)
		return
	}

	if category.Id == 0 {
		logs.Error("id not found")
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(updateCategoryErr),
			"id not found",
		)
		return
	}

	updateData := entities.CategoryUpdate{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(updateCategoryErr),
			err.Error(),
		)
		return
	}

	category.Name = updateData.Name
	category.UpdatedAt = time.Now()

	// database save update
	err = h.db.Table(tableCategories).Save(&category).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(updateCategoryErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusOK, category)
}

func (h *category) DeleteCategory(c *gin.Context) {
	id := c.Param("categoryId")
	var category entities.Category

	// database find id
	err := h.db.Table(tableCategories).Where("id = ?", id).First(&category).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(updateCategoryErr),
			err.Error(),
		)
		return
	}

	// database delete
	err = h.db.Table(tableCategories).Where("id = ?", id).Delete(&category).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(updateCategoryErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusNoContent, nil)
}
