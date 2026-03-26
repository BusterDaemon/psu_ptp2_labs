package tests

import (
	"math"
	"product_api/entity"
	"product_api/internal/db"
	"testing"

	"github.com/sdrapkin/guid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSqliteConnection(t *testing.T) {
	s := db.SQLiteDatabase{}
	err := s.CreateDBConnection(":memory:")

	assert.Nil(t, err)
}

func TestAddSqliteProduct(t *testing.T) {
	s := db.SQLiteDatabase{}
	err := s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	p := entity.Product{
		Id:         guid.NewString(),
		Definition: "GOIDA GOIDA",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)
}

func TestDeleteSqliteProduct(t *testing.T) {
	s := db.SQLiteDatabase{}
	err := s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	gid := guid.NewString()
	p := entity.Product{
		Id:         gid,
		Definition: "",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)

	err = s.DeleteProduct(gid)
	assert.Nil(t, err)
}

func TestSearchSqliteProduct(t *testing.T) {
	var (
		s  db.SQLiteDatabase
		ps []entity.Product
	)

	err := s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	gid := guid.NewString()
	p := entity.Product{
		Id:         gid,
		Definition: "",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)

	ps, err = s.SearchProduct("Пригожин Женя")
	assert.Nil(t, err)
	assert.Greater(t, len(ps), 0)

	ps, err = s.SearchProduct("Пригожин")
	assert.Nil(t, err)
	assert.Greater(t, len(ps), 0)

	ps, err = s.SearchProduct("Женя")
	assert.Nil(t, err)
	assert.Greater(t, len(ps), 0)
}

func TestUpdateSqliteProduct(t *testing.T) {
	s := db.SQLiteDatabase{}
	err := s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	gid := guid.NewString()
	p := entity.Product{
		Id:         gid,
		Definition: "",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)

	p.Name = "ГОЙДА ГОЙДА, РАКЕТЫ ЛЕТЯТ"
	p.Price = 88.88

	err = s.UpdateProduct(&p)
	assert.Nil(t, err)
}

func TestGetByIDSqliteProduct(t *testing.T) {
	var (
		s  db.SQLiteDatabase
		p2 entity.Product
	)
	err := s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	gid := guid.NewString()
	p := entity.Product{
		Id:         gid,
		Definition: "",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)

	p2, err = s.GetProductByID(gid)
	assert.Nil(t, err)
	assert.False(t, p2.Id == "")

}

func TestDropDBAndReinsertSqlite(t *testing.T) {
	var (
		s    db.SQLiteDatabase
		p    entity.Product
		prds []entity.Product
		err  error
	)

	err = s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	p = entity.Product{
		Id:         guid.NewString(),
		Definition: "",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)

	prds = append(prds, entity.Product{
		Id:         guid.NewString(),
		Definition: "СЛАВНЫЙ РУССКИЙ ГОРОД",
		Name:       "ХАРЬКОВ",
		Price:      float32(math.Inf(1)),
		Image:      "",
	})
	prds = append(prds, entity.Product{
		Id:         guid.NewString(),
		Definition: "НАЦИОНАЛЬНАЯ ИДЕЯ",
		Name:       "ГОЙДА",
		Price:      float32(math.Inf(-1)),
		Image:      "",
	})
	prds = append(prds, entity.Product{
		Id:         guid.NewString(),
		Definition: "ГЛАВНЫЙ ЗЛОДЕЙ",
		Name:       "ЛНР",
		Price:      42.69,
		Image:      "",
	})

	err = s.DropDBAndReinsert(prds)
	assert.Nil(t, err)
}

func TestGetAllProductsSqlite(t *testing.T) {
	var (
		s    db.SQLiteDatabase
		p    entity.Product
		prds []entity.Product
		err  error
	)

	err = s.CreateDBConnection(":memory:")
	assert.Nil(t, err)

	p = entity.Product{
		Id:         guid.NewString(),
		Definition: "",
		Name:       "Пригожин Женя",
		Price:      14.88,
		Image:      "",
	}
	err = s.AddProduct(&p)
	assert.Nil(t, err)

	prds = append(prds, entity.Product{
		Id:         guid.NewString(),
		Definition: "СЛАВНЫЙ РУССКИЙ ГОРОД",
		Name:       "ХАРЬКОВ",
		Price:      float32(math.Inf(1)),
		Image:      "",
	})
	prds = append(prds, entity.Product{
		Id:         guid.NewString(),
		Definition: "НАЦИОНАЛЬНАЯ ИДЕЯ",
		Name:       "ГОЙДА",
		Price:      float32(math.Inf(-1)),
		Image:      "",
	})
	prds = append(prds, entity.Product{
		Id:         guid.NewString(),
		Definition: "ГЛАВНЫЙ ЗЛОДЕЙ",
		Name:       "ЛНР",
		Price:      42.69,
		Image:      "",
	})

	err = s.DropDBAndReinsert(prds)
	assert.Nil(t, err)

	prds, err = s.GetAllProducts()
	assert.Nil(t, err)
	assert.Equal(t, len(prds), 3)
}
