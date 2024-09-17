package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rcjeferson/go-api-products/internal/controller"
	"github.com/rcjeferson/go-api-products/internal/db"
	"github.com/rcjeferson/go-api-products/internal/repository"
	"github.com/rcjeferson/go-api-products/internal/usecase"
)

func main() {
	server := gin.Default()

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	productRepository := repository.NewProductRepository(dbConnection)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productController := controller.NewProductController(productUseCase)

	server.GET("/products", productController.GetProducts)
	server.POST("/product", productController.CreateProduct)

	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.Run(":8000")
}
