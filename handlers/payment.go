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

var tablePayment = entities.TablePayment()
var tablePaymentType = entities.TablePaymentType()

type paymentErr string

const (
	paymentCreateErr    paymentErr = "payment-001"
	paymentGetErr       paymentErr = "payment-002"
	paymentGetDetailErr paymentErr = "payment-003"
	paymentUpdateErr    paymentErr = "payment-004"
	paymentDeleteErr    paymentErr = "payment-005"
)

type IPayment interface {
	CreatePayment(c *gin.Context)
	GetPayments(c *gin.Context)
	GetPaymentDetails(c *gin.Context)
	UpdatePayment(c *gin.Context)
	DeletePayment(c *gin.Context)
}

type payment struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) IPayment {
	return &payment{db: db}
}

func (h *payment) CreatePayment(c *gin.Context) {
	req := entities.PaymentRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(paymentCreateErr),
			err.Error(),
		)
		return
	}

	if req.Name == "" || req.Type == "" {
		logs.Error("req is empty")
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(paymentCreateErr),
			"name and type is required",
		)
		return
	}

	payTypes := entities.PaymentType{}

	// check payment type
	err := h.db.Table(tablePaymentType).Where("name", req.Type).First(&payTypes).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentCreateErr),
			err.Error(),
		)
		return
	}

	paymentData := entities.Payment{
		Name:          req.Name,
		Type:          req.Type,
		PaymentTypeId: payTypes.Id,
		Logo:          req.Logo,
	}

	// create payment
	err = h.db.Table(tablePayment).Create(&paymentData).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentCreateErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusCreated, paymentData)
}

func (h *payment) GetPayments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	skip, _ := strconv.Atoi(c.Query("skip"))
	var count int64
	var payments []entities.Payment

	if limit == 0 {
		limit = 10
	}

	// select payment
	err := h.db.Table(tablePayment).Select("id ,name,type,payment_type_id,logo,created_at,updated_at").Limit(limit).Offset(skip).Find(&payments).Count(&count).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentGetErr),
			err.Error(),
		)
		return
	}

	responseData := entities.PaymentListResponse{
		Payments: payments,
		Meta: &entities.QueryMeta{
			Total: int(count),
			Limit: limit,
			Skip:  skip,
		},
	}

	responses.NewResponse(c).Success(http.StatusOK, responseData)
}

func (h *payment) GetPaymentDetails(c *gin.Context) {
	paymentId := c.Param("paymentId")
	payment := entities.Payment{}

	err := h.db.Table(tablePayment).Where("id = ?", paymentId).First(&payment).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentGetDetailErr),
			err.Error(),
		)
		return
	}

	if payment.Name == "" {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(paymentGetDetailErr),
			"payment id not found",
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusOK, payment)
}

func (h *payment) UpdatePayment(c *gin.Context) {
	paymentId := c.Param("paymentId")
	payment := entities.Payment{}

	// find payment
	err := h.db.Table(tablePayment).Find(&payment).Where("id = ?", paymentId).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentUpdateErr),
			err.Error(),
		)
		return
	}

	if payment.Id == 0 {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(paymentUpdateErr),
			"payment id not found",
		)
		return
	}

	updateData := entities.Payment{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusBadRequest,
			string(paymentUpdateErr),
			err.Error(),
		)
		return
	}

	if updateData.Name == "" {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(paymentUpdateErr),
			"payment name is required",
		)
		return
	}

	var paymentTypeId int

	if updateData.Type == "CASH" {
		paymentTypeId = 1
	}
	if updateData.Type == "TRANSFER" {
		paymentTypeId = 2
	}
	if updateData.Type == "ETC" {
		paymentTypeId = 3
	}

	payment.Name = updateData.Name
	payment.Type = updateData.Type
	payment.PaymentTypeId = paymentTypeId
	payment.Logo = updateData.Logo
	payment.UpdatedAt = time.Now()

	// save payment
	err = h.db.Table(tablePayment).Save(&payment).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentUpdateErr),
			err.Error(),
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusOK, payment)
}

func (h *payment) DeletePayment(c *gin.Context) {
	paymentId := c.Param("paymentId")
	payment := entities.Payment{}

	err := h.db.Table(tablePayment).First(&payment, paymentId).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentDeleteErr),
			err.Error(),
		)
		return
	}

	if payment.Name == "" {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusNotFound,
			string(paymentDeleteErr),
			"payment id not found",
		)
		return
	}

	err = h.db.Table(tablePayment).Delete(&payment).Error
	if err != nil {
		logs.Error(err)
		responses.NewResponse(c).Error(
			http.StatusInternalServerError,
			string(paymentDeleteErr),
			"payment delete failed",
		)
		return
	}

	responses.NewResponse(c).Success(http.StatusNoContent, nil)
}
