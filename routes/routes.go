package routes

import (
	"github.com/codepnw/sales-api/database"
	cathandlers "github.com/codepnw/sales-api/modules/categories/handlers"
	catrepositories "github.com/codepnw/sales-api/modules/categories/repositories"
	catservices "github.com/codepnw/sales-api/modules/categories/services"
	prodhandlers "github.com/codepnw/sales-api/modules/products/handlers"
	prodrepositories "github.com/codepnw/sales-api/modules/products/repositories"
	prodservices "github.com/codepnw/sales-api/modules/products/services"
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, version string) {
	productRoutes(router, version)
	categoryRoutes(router, version)
}

func productRoutes(router *gin.Engine, version string) {
	repo := prodrepositories.NewProductRepository(database.GetPostgresDB())
	srv := prodservices.NewProductService(repo)
	h := prodhandlers.NewProductHandler(srv)
	g := router.Group(version + "/products")
	paramId := "/:productId"

	g.POST("/", h.CreateProduct)
	g.GET("/", h.GetProducts)
	g.GET(paramId, h.GetProduct)
	g.PATCH(paramId, h.UpdateProduct)
	g.DELETE(paramId, h.DeleteProduct)
}

func categoryRoutes(router *gin.Engine, version string) {
	repo := catrepositories.NewCategoryRepository(database.GetPostgresDB())
	srv := catservices.NewCategoryService(repo)
	h := cathandlers.NewCategoryHandler(srv)
	g := router.Group(version + "/categories")
	paramId := "/:categoryId"

	g.POST("/", h.CreateCategory)
	g.GET("/", h.GetAllCategory)
	g.GET(paramId, h.GetOneCategory)
	g.PATCH(paramId, h.UpdateCategory)
	g.DELETE(paramId, h.DeleteCategory)
}
