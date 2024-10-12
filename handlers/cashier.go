package handlers

import (
	"net/http"
	"strconv"

	"github.com/codepnw/sales-api/entities"
	"github.com/codepnw/sales-api/pkg/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrCashier string

const (
	cashierCreateErr ErrCashier = "cashier-001"
	cashierGetAllErr ErrCashier = "cashier-002"
	cashierGetOneErr ErrCashier = "cashier-003"
	cashierUpdateErr ErrCashier = "cashier-004"
	cashierDeleteErr ErrCashier = "cashier-005"
)

type ICashier interface {
	CreateCashier(c *gin.Context)
	GetCashiers(c *gin.Context)
	GetCashierDetails(c *gin.Context)
	UpdateCashier(c *gin.Context)
	DeleteCashier(c *gin.Context)
}

type cashier struct {
	db *gorm.DB
}

func NewCashierHandler(db *gorm.DB) ICashier {
	return &cashier{db: db}
}

func (h *cashier) CreateCashier(c *gin.Context) {
	var req entities.Cashier

	if err := c.ShouldBindJSON(&req); err != nil {
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(cashierCreateErr),
			err.Error(),
		)
		return
	}

	if req.Name == "" || req.Passcode == "" {
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(cashierCreateErr),
			"name and passcode is required",
		)
		return
	}

	cashier := entities.Cashier{
		Name:     req.Name,
		Passcode: req.Passcode,
	}

	// database create
	h.db.Create(&cashier)

	responses.NewResponse(c).Success(http.StatusCreated, "cashier created")
}

func (h *cashier) GetCashiers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var cashiers []entities.Cashiers

	if limit == 0 {
		limit = 10
	}

	// database select
	h.db.Select("id, name").Limit(limit).Offset(skip).Find(&cashiers).Count(&count)

	metaMap := map[string]any{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	cashiersData := map[string]any{
		"cashiers": cashiers,
		"meta":     metaMap,
	}

	responses.NewResponse(c).Success(http.StatusOK, cashiersData)
}

func (h *cashier) GetCashierDetails(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database select
	h.db.Select("id, name").Where("id=?", cashierId).First(&cashier)

	data := make(map[string]any)
	data["cashierId"] = cashier.Id
	data["name"] = cashier.Name

	if cashier.Id == 0 {
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(cashierGetOneErr),
			"id not found",
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusOK, data)
}

func (h *cashier) UpdateCashier(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database find
	h.db.Find(&cashier, "id = ?", cashierId)

	if cashier.Name == "" {
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(cashierUpdateErr),
			"name is required",
		)
		return
	}

	var updateData entities.Cashier
	if err := c.ShouldBindJSON(&updateData); err != nil {
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(cashierUpdateErr),
			err.Error(),
		)
		return
	}

	cashier.Name = updateData.Name
	// database save
	h.db.Save(&cashier)

	responses.NewResponse(c).Success(http.StatusOK, cashier)
}

func (h *cashier) DeleteCashier(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database find id
	h.db.Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(cashierDeleteErr),
			"id not found",
		)
		return
	}
	// database delete
	h.db.Where("id = ?", cashierId).Delete(&cashier)

	responses.NewResponse(c).Success(http.StatusNoContent, nil)
}
