package db

import (
	"product_api/entity"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (d *Database) CreateDBConnection(db_path string) error {
	db, err := gorm.Open(sqlite.Open(db_path), &gorm.Config{})
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

func (d Database) DeleteProduct(id string) (int64, error) {
	tx := d.connection.Delete(&entity.Product{Id: id})
	if tx.Error != nil {
		return 0, tx.Error
	}

	return tx.RowsAffected, nil
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

func (d Database) UpdateProduct(product *entity.Product) (int64, error) {
	tx := d.connection.Save(product)
	if tx.Error != nil {
		switch tx.Error {
		case gorm.ErrRecordNotFound:
			return 0, nil
		default:
			return 0, tx.Error
		}
	}

	return tx.RowsAffected, nil
}

func (d Database) GetProductByID(id string) (entity.Product, error) {
	var product entity.Product
	tx := d.connection.First(&product, "id = ?", id)
	if tx.Error != nil {
		return entity.Product{}, tx.Error
	}

	return product, nil
}

func (d Database) DropDBAndReinsert(prds []entity.Product) error {
	var mt sync.Mutex

	mt.Lock()
	err := d.connection.Migrator().DropTable(&entity.Product{})
	if err != nil {
		return err
	}
	err = d.connection.AutoMigrate(&entity.Product{})
	if err != nil {
		return err
	}

	for _, v := range prds {
		d.connection.Create(v)
	}
	mt.Unlock()

	return nil
}

func (d Database) GetAllProducts() ([]entity.Product, error) {
	var prds []entity.Product
	tx := d.connection.Find(&prds)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return prds, nil
}
