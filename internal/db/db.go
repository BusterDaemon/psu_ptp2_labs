package db

import (
	"database/sql"
	"product_api/entity"

	"gorm.io/gorm"
)

type Databaser interface {
	CreateDBConnection(db_path string) error
	AddProduct(product *entity.Product) error
	DeleteProduct(id string) (int64, error)
	SearchProduct(name string) ([]entity.Product, error)
	UpdateProduct(product *entity.Product) (int64, error)
	GetProductByID(id string) (entity.Product, error)
	DropDBAndReinsert(prds []entity.Product) error
	GetAllProducts() ([]entity.Product, error)
}

type Database struct {
	connection *gorm.DB
}

type SQLiteDatabase struct {
	connection *sql.DB
}
