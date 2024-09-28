package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rcjeferson/go-api-products/internal/controller"
	"github.com/rcjeferson/go-api-products/internal/db"
	"github.com/rcjeferson/go-api-products/internal/repository"
	"github.com/rcjeferson/go-api-products/internal/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("The .env file not found!")
	}
}

func main() {
	server := gin.Default()

	rdb, _ := db.ConnectRedis()
	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	useCache, err := strconv.ParseBool(os.Getenv("USE_CACHE"))
	if err != nil {
		slog.Error("Failed to get USE_CACHE environment variable!")
	}

	// Product endpoints
	productRepository := repository.NewProductRepository(dbConnection, useCache, rdb)
	productUseCase := usecase.NewProductUseCase(productRepository)
	productController := controller.NewProductController(productUseCase)

	server.GET("/products", productController.GetProducts)
	server.GET("/products/:id", productController.GetProductById)
	server.POST("/products", productController.CreateProduct)

	// Ping and health endpoints
	server.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	healthController := controller.NewHealthController(dbConnection)
	server.GET("/health", healthController.Check)

	server.Run(":8000")
}
