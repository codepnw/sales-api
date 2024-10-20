package products

import "time"

type Product struct {
	ProductID  string     `db:"product_id" json:"productId"`
	Name       string     `db:"name" json:"name"`
	Desc       string     `db:"desc" json:"desc"`
	Price      float64    `db:"price" json:"price"`
	Discount   uint       `db:"discount" json:"discount"`
	Stock      uint       `db:"stock" json:"stock"`
	CategoryID uint       `db:"category_id" json:"categoryId"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}

type ProductRequest struct {
	Name       string  `json:"name" form:"name"`
	Desc       string  `json:"desc" form:"desc"`
	Price      float64 `json:"price" form:"price"`
	Discount   uint    `json:"discount" form:"discount"`
	Stock      uint    `json:"stock" form:"stock"`
	CategoryID uint    `json:"categoryId" form:"category_id"`
}
