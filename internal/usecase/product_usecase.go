package usecase

import (
	"log/slog"

	"github.com/rcjeferson/go-api-products/internal/model"
	"github.com/rcjeferson/go-api-products/internal/repository"
)

type ProductUseCase struct {
	repository repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return ProductUseCase{
		repository: repo,
	}
}

func (pu *ProductUseCase) GetProducts() ([]model.Product, error) {
	return pu.repository.GetProducts()
}

func (pu *ProductUseCase) CreateProduct(product model.Product) (model.Product, error) {
	id, err := pu.repository.CreateProduct(product)

	if err != nil {
		slog.Error("Error while creating a product on CreateProduct UseCase: ", err)
		return model.Product{}, err
	}

	product.ID = id

	return product, nil
}
