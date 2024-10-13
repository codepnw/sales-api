package handlers

import (
	"net/http"
	"strconv"

	"github.com/codepnw/sales-api/entities"
	"github.com/codepnw/sales-api/pkg/logs"
	"github.com/codepnw/sales-api/pkg/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tableCashier string = entities.TableCashier()

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
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(cashierCreateErr),
			"cashier request invalid",
		)
		return
	}

	if req.Name == "" || req.Passcode == "" {
		logs.Error("request is empty")
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
	err := h.db.Table(tableCashier).Create(&cashier).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierCreateErr),
			"failed create cashier",
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusCreated, "cashier created")
}

func (h *cashier) GetCashiers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var cashiers []entities.CashierResponse

	if limit == 0 {
		limit = 10
	}

	// database select
	err := h.db.Table(tableCashier).Select("id, name").Limit(limit).Offset(skip).Find(&cashiers).Count(&count).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierGetAllErr),
			err.Error(),
		)
	}

	cashiersData := entities.CashierListResponse{
		Cashiers: cashiers,
		Meta: &entities.CashiersMeta{
			Total: int(count),
			Limit: limit,
			Skip:  skip,
		},
	}

	responses.NewResponse(c).Success(http.StatusOK, cashiersData)
}

func (h *cashier) GetCashierDetails(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database select
	err := h.db.Table(tableCashier).Select("id, name").Where("id=?", cashierId).First(&cashier).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierGetOneErr),
			err.Error(),
		)
		return
	}

	data := entities.CashierResponse{
		Id:   cashier.Id,
		Name: cashier.Name,
	}

	responses.NewResponse(c).Success(http.StatusOK, data)
}

func (h *cashier) UpdateCashier(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.CashierResponse

	// database find
	err := h.db.Table(tableCashier).Find(&cashier, "id = ?", cashierId).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierUpdateErr),
			err.Error(),
		)
		return
	}

	if cashier.Name == "" {
		logs.Error("name is empty")
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(cashierUpdateErr),
			"name is required",
		)
		return
	}

	var updateData entities.CashierResponse
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
	err = h.db.Table(tableCashier).Save(&cashier).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierUpdateErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusOK, cashier)
}

func (h *cashier) DeleteCashier(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database find id
	err := h.db.Table(tableCashier).Where("id = ?", cashierId).First(&cashier).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierDeleteErr),
			err.Error(),
		)
		return
	}

	if cashier.Id == 0 {
		logs.Error("id not found")
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(cashierDeleteErr),
			"id not found",
		)
		return
	}

	// database delete
	err = h.db.Table(tableCashier).Where("id = ?", cashierId).Delete(&cashier).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(cashierUpdateErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusNoContent, nil)
}
