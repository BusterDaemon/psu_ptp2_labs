package db

import (
	"product_api/entity"

	"gorm.io/gorm"
)

type Databaser interface {
	CreateDBConnection(db_path string) error
	AddProduct(product *entity.Product) error
	DeleteProduct(id string) error
	SearchProduct(name string) ([]entity.Product, error)
	UpdateProduct(product *entity.Product) error
	GetProductByID(id string) (entity.Product, error)
	DropDBAndReinsert(prds []entity.Product) error
}

type Database struct {
	connection *gorm.DB
}
