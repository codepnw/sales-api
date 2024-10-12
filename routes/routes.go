package routes

import (
	"github.com/codepnw/sales-api/database"
	"github.com/codepnw/sales-api/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, version string) {
	cashierRoute(router, version)
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
