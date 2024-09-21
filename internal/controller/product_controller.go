package controller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rcjeferson/go-api-products/internal/model"
	"github.com/rcjeferson/go-api-products/internal/usecase"
)

type productController struct {
	productUseCase usecase.ProductUseCase
}

func NewProductController(usecase usecase.ProductUseCase) productController {
	return productController{
		productUseCase: usecase,
	}
}

func (p *productController) GetProducts(ctx *gin.Context) {
	products, err := p.productUseCase.GetProducts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, products)
}

func (p *productController) GetProductById(ctx *gin.Context) {
	var product model.Product

	err := ctx.ShouldBindUri(&product)

	if err != nil {
		slog.Error("Error to bind JSON on Get By ID Request: ", err)

		response := model.Response{
			Message: "id must be a integer number",
		}

		ctx.JSON(http.StatusBadRequest, response)

		return
	}

	err = p.productUseCase.GetProductById(&product)
	if err != nil {
		slog.Error("Error to get product by id on GetProductById Controller: ", err)

		response := model.Response{
			Message: "error to get product",
		}

		ctx.JSON(http.StatusInternalServerError, response)

		return
	}

	if product.ID == 0 {
		response := model.Response{
			Message: "product not found",
		}
		ctx.JSON(http.StatusNotFound, response)

		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *productController) CreateProduct(ctx *gin.Context) {
	var product model.Product

	err := ctx.BindJSON(&product)

	if err != nil {
		slog.Error("Error to bind JSON on Post Request: ", err)

		response := model.Response{
			Message: "invalid parameters",
		}

		ctx.JSON(http.StatusBadRequest, response)

		return
	}

	err = p.productUseCase.CreateProduct(&product)

	if err != nil {
		slog.Error("Error on usecase while creating product: ", err)

		response := model.Response{
			Message: "error to create a new product",
		}

		ctx.JSON(http.StatusInternalServerError, response)

		return
	}

	ctx.JSON(http.StatusCreated, product)

}
