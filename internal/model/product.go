package model

type Product struct {
	ID    int     `json:"id" uri:"id" binding:"required"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
