package routes

import (
	"github.com/codepnw/sales-api/database"
	"github.com/codepnw/sales-api/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, version string) {
	cashierRoute(router, version)
	categoryRoute(router, version)
	paymentRoutes(router, version)
}

func cashierRoute(router *gin.Engine, version string) {
	h := handlers.NewCashierHandler(database.GetDB())
	g := router.Group(version + "/cashiers")

	g.GET("/", h.GetCashiers)
	g.GET("/:cashierId", h.GetCashierDetails)
	g.POST("/", h.CreateCashier)
	g.PATCH("/:cashierId", h.UpdateCashier)
	g.DELETE("/:cashierId", h.DeleteCashier)
}

func categoryRoute(router *gin.Engine, version string) {
	h := handlers.NewCategoryHandler(database.GetDB())
	g := router.Group(version + "/categories")
	id := ":categoryId"

	g.GET("/", h.GetCategories)
	g.GET("/"+id, h.GetCategory)
	g.POST("/", h.CreateCategory)
	g.PATCH("/"+id, h.UpdateCategory)
	g.DELETE("/"+id, h.DeleteCategory)
}

func paymentRoutes(router *gin.Engine, version string) {
	h := handlers.NewPaymentHandler(database.GetDB())
	g := router.Group(version + "/payments")
	paramId := "/:paymentId"

	g.POST("/", h.CreatePayment)
	g.GET("/", h.GetPayments)
	g.GET(paramId, h.GetPaymentDetails)
	g.PATCH(paramId, h.UpdatePayment)
	g.DELETE(paramId, h.DeletePayment)
}
