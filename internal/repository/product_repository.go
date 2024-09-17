package repository

import (
	"database/sql"
	"fmt"
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
		fmt.Println(err)
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
			fmt.Println(err)
			return []model.Product{}, err
		}

		productList = append(productList, productObject)
	}
	rows.Close()

	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int
	stmt, err := pr.connection.Prepare("INSERT INTO product(name, price) VALUES($1, $2) RETURNING id")

	if err != nil {
		slog.Error("Error while prepare sql statement on CreateProduct Repository: ", err)
		return 0, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(product.Name, product.Price).Scan(&id)

	if err != nil {
		slog.Error("Error while executing query on CreateProduct Repository: ", err)
		return 0, err
	}

	return id, err
}
