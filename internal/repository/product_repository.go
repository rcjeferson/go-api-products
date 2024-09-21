package repository

import (
	"database/sql"
	"log/slog"

	"github.com/rcjeferson/go-api-products/internal/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT id, name, price FROM product"

	rows, err := pr.connection.Query(query)
	if err != nil {
		slog.Error("Error while getting prodycts on GetProducts Repository: ", err)
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObject model.Product

	for rows.Next() {
		err = rows.Scan(
			&productObject.ID,
			&productObject.Name,
			&productObject.Price,
		)

		if err != nil {
			slog.Error("Error while getting products on GetProducts Repository: ", err)
			return []model.Product{}, err
		}

		productList = append(productList, productObject)
	}
	rows.Close()

	return productList, nil
}

func (pr *ProductRepository) GetProductById(product *model.Product) error {
	query, err := pr.connection.Prepare("SELECT id, name, price FROM product WHERE id = $1")
	if err != nil {
		slog.Error("Error while prepare query on GetProductById Repository: ", err)
		return err
	}

	err = query.QueryRow(product.ID).Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			product.ID = 0
			return nil
		}

		slog.Error("Error while getting product on GetProductById Repository: ", err)
		return err
	}

	query.Close()

	return nil
}

func (pr *ProductRepository) CreateProduct(product *model.Product) error {
	stmt, err := pr.connection.Prepare("INSERT INTO product(name, price) VALUES($1, $2) RETURNING id")

	if err != nil {
		slog.Error("Error while prepare sql statement on CreateProduct Repository: ", err)
		return err
	}

	err = stmt.QueryRow(product.Name, product.Price).Scan(&product.ID)

	if err != nil {
		slog.Error("Error while executing query on CreateProduct Repository: ", err)
		return err
	}

	stmt.Close()

	return nil
}
