package handlers

import (
	"net/http"
	"strconv"

	"github.com/codepnw/sales-api/entities"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "data bad request",
			"error":   err.Error(),
		})
		return
	}

	if req.Name == "" || req.Passcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "name and passcode is required",
			"error":   nil,
		})
		return
	}

	cashier := entities.Cashier{
		Name:     req.Name,
		Passcode: req.Passcode,
	}

	// database create
	h.db.Create(&cashier)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "cashier created",
	})
}



func (h *cashier) GetCashiers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var cashiers []entities.CashiersResponse

	// database select
	h.db.Select("*").Limit(limit).Offset(skip).Find(&cashiers).Count(&count)

	metaMap := map[string]any{
		"total": count,
		"limit": limit,
		"skip":  skip,
	}

	cashiersData := map[string]any{
		"cashiers": cashiers,
		"meta":     metaMap,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data":    cashiersData,
	})
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
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "cashier id not found",
			"error":   map[string]any{},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data":    data,
	})
}

func (h *cashier) UpdateCashier(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database find
	h.db.Find(&cashier, "id = ?", cashierId)

	if cashier.Name == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "cashier not found",
		})
		return
	}

	var updateData entities.Cashier
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "cashier name is required",
			"error":   err.Error(),
		})
		return
	}

	cashier.Name = updateData.Name
	// database save
	h.db.Save(&cashier)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "success",
		"data":    cashier,
	})
}

func (h *cashier) DeleteCashier(c *gin.Context) {
	cashierId := c.Param("cashierId")
	var cashier entities.Cashier

	// database find id
	h.db.Where("id = ?", cashierId).First(&cashier)

	if cashier.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "cashier not found",
		})
		return
	}
	// database delete
	h.db.Where("id = ?", cashierId).Delete(&cashier)

	c.JSON(http.StatusNoContent, gin.H{
		"success": true,
		"message": "message",
	})
}
