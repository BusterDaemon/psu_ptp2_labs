package db

import (
	"product_api/entity"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (d *Database) CreateDBConnection() error {
	db, err := gorm.Open(sqlite.Open("store.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	db = db.Debug()

	db.AutoMigrate(&entity.Product{})

	d.connection = db

	return nil
}

func (d Database) AddProduct(product *entity.Product) error {
	tx := d.connection.Create(product)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (d Database) DeleteProduct(id string) error {
	tx := d.connection.Delete(&entity.Product{Id: id})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (d Database) SearchProduct(name string) ([]entity.Product, error) {
	var products []entity.Product

	tx := d.connection.Where("name LIKE ? COLLATE NOCASE", "%"+name+"%").Find(&products)
	if tx.Error != nil {
		switch tx.Error {
		case gorm.ErrRecordNotFound:
			return nil, nil
		default:
			return nil, tx.Error
		}
	}

	return products, nil
}

func (d Database) UpdateProduct(product *entity.Product) error {
	tx := d.connection.Save(product)
	if tx.Error != nil {
		switch tx.Error {
		case gorm.ErrRecordNotFound:
			return nil
		default:
			return tx.Error
		}
	}

	return nil
}

func (d Database) GetProductByID(id string) (entity.Product, error) {
	var product entity.Product
	tx := d.connection.First(&product, "id = ?", id)
	if tx.Error != nil {
		return entity.Product{}, tx.Error
	}

	return product, nil
}
