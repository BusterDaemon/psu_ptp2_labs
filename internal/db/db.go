package db

import (
	"product_api/entity"

	"gorm.io/gorm"
)

type Databaser interface {
	CreateDBConnection() (*gorm.DB, error)
	AddProduct(product *entity.Product) error
	DeleteProduct(id string) error
	SearchProduct(name string) ([]entity.Product, error)
	UpdateProduct(product *entity.Product) error
	GetProductByID(id string) (entity.Product, error)
}

type Database struct {
	connection *gorm.DB
}
