package db

import (
	"database/sql"
	"product_api/entity"

	_ "github.com/mattn/go-sqlite3"
)

func (dbase *SQLiteDatabase) CreateDBConnection(path string) error {
	var err error

	dbase.connection, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	_, err = dbase.connection.Exec(`CREATE TABLE IF NOT EXISTS PRODUCTS(
	ID TEXT CONSTRAINT "ident" PRIMARY KEY DESC ON CONFLICT ABORT,
	DEFINITION TEXT,
	NAME TEXT NOT NULL ON CONFLICT ABORT,
	PRICE REAL DEFAULT 0.0,
	IMAGE TEXT
);
CREATE UNIQUE INDEX IF NOT EXISTS IDENT ON PRODUCTS (ID);`)
	if err != nil {
		return err
	}

	return nil
}

func (dbase SQLiteDatabase) AddProduct(p *entity.Product) error {
	_, err := dbase.connection.Exec("INSERT INTO PRODUCTS VALUES (?, ?, ?, ?, ?)", p.Id, p.Definition, p.Name, p.Price, p.Image)
	if err != nil {
		return err
	}

	return nil
}

func (dbase SQLiteDatabase) DeleteProduct(id string) (int64, error) {
	var rws int64

	res, err := dbase.connection.Exec("DELETE FROM PRODUCTS WHERE ID = ?", id)
	if err != nil {
		return 0, err
	}

	if rws, err = res.RowsAffected(); err != nil {
		return -1, err
	}

	return rws, nil
}

func (dbase SQLiteDatabase) SearchProduct(name string) ([]entity.Product, error) {
	var (
		products []entity.Product
		err      error
	)

	rows, err := dbase.connection.Query(`SELECT * FROM PRODUCTS WHERE NAME LIKE '%' || ? || '%'`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p entity.Product

		if err = rows.Scan(&p.Id, &p.Definition, &p.Name, &p.Price, &p.Image); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (dbase SQLiteDatabase) UpdateProduct(product *entity.Product) (int64, error) {
	res, err := dbase.connection.Exec("UPDATE OR ABORT PRODUCTS SET (DEFINITION, NAME, PRICE, IMAGE) = (?, ?, ?, ?) WHERE ID = ?",
		product.Definition, product.Name, product.Price, product.Image, product.Id)
	if err != nil {
		return 0, err
	}

	rws, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rws, nil
}

func (dbase SQLiteDatabase) GetProductByID(id string) (entity.Product, error) {
	var p entity.Product
	row := dbase.connection.QueryRow("SELECT * FROM PRODUCTS WHERE ID = ?", id)
	err := row.Scan(&p.Id, &p.Definition, &p.Name, &p.Price, &p.Image)
	if err != nil {
		return entity.Product{}, err
	}
	return p, nil
}

func (dbase SQLiteDatabase) DropDBAndReinsert(prds []entity.Product) error {
	_, err := dbase.connection.Exec("DELETE FROM PRODUCTS")
	if err != nil {
		return err
	}
	tx, err := dbase.connection.Begin()
	if err != nil {
		return err
	}
	for _, v := range prds {
		_, err := tx.Exec("INSERT INTO PRODUCTS VALUES (?, ?, ?, ?, ?)", v.Id, v.Definition, v.Name, v.Price, v.Image)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (dbase SQLiteDatabase) GetAllProducts() ([]entity.Product, error) {
	var (
		prds []entity.Product
		err  error
	)
	rows, err := dbase.connection.Query("SELECT * FROM PRODUCTS")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p entity.Product
		if err = rows.Scan(&p.Id, &p.Definition, &p.Name, &p.Price, &p.Image); err != nil {
			return nil, err
		}

		prds = append(prds, p)
	}

	return prds, nil
}
