package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/rcjeferson/go-api-products/internal/db"
	"github.com/rcjeferson/go-api-products/internal/model"
	"github.com/redis/go-redis/v9"
)

type ProductRepository struct {
	connection *sql.DB
	rdb        *redis.Client
	useCache   bool
}

func NewProductRepository(connection *sql.DB, useCache bool, rdbClient *redis.Client) ProductRepository {
	return ProductRepository{
		connection: connection,
		rdb:        rdbClient,
		useCache:   useCache,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	var productList []model.Product
	var productObject model.Product
	rdbKey := "products"

	if pr.useCache {
		slog.Info("Using cache on GetProducts Repository!")

		cachedProducts, err := db.GetCache(rdbKey, pr.rdb)
		if err == nil {
			err := json.Unmarshal([]byte(cachedProducts), &productList)
			if err != nil {
				slog.Error("Error while Unmarshal JSON object from cache on GetProducts Repository: ", err)
			}

			return productList, nil
		}
	} else {
		slog.Info("Not using cache on GetProducts Repository!")
	}

	query := "SELECT id, name, price FROM product"

	rows, err := pr.connection.Query(query)
	if err != nil {
		slog.Error("Error while getting products on GetProducts Repository: ", err)
		return []model.Product{}, err
	}

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

	if pr.useCache {
		mProducts, _ := json.Marshal(productList)
		db.SetCache(rdbKey, mProducts, 60, pr.rdb)
	}

	return productList, nil
}

func (pr *ProductRepository) GetProductById(product *model.Product) error {
	rdbKey := fmt.Sprintf("product:%d", product.ID)

	if pr.useCache {
		slog.Info("Using cache on GetProductById Repository!")

		cachedProduct, err := db.GetCache(rdbKey, pr.rdb)
		if err == nil {
			err := json.Unmarshal([]byte(cachedProduct), &product)
			if err != nil {
				slog.Error("Error while Unmarshal JSON object from cache on GetProductById Repository: ", err)
			}

			return nil
		}
	} else {
		slog.Info("Not using cache on GetProductById Repository!")
	}

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

	if pr.useCache {
		mProduct, _ := json.Marshal(product)
		db.SetCache(rdbKey, mProduct, 60, pr.rdb)
	}

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
