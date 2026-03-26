package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"product_api/entity"
	"product_api/internal/db"
	"product_api/internal/services"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetAllController(t *testing.T) {
	app := services.NewServiceSetup(&db.SQLiteDatabase{})

	req, _ := http.NewRequest("GET", "http://127.0.0.1:1111/api/get_all", nil)

	res, err := app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestAddController(t *testing.T) {
	app := services.NewServiceSetup(&db.SQLiteDatabase{})

	data := url.Values{
		"name":  {"SVOй Vайбик"},
		"price": {"1488.0"},
	}

	req, _ := http.NewRequest("POST", "http://127.0.0.1:1111/api/add", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)

	assert.Equal(t, 202, res.StatusCode)
}

func TestSearchController(t *testing.T) {
	app := services.NewServiceSetup(&db.SQLiteDatabase{})

	data := url.Values{
		"name":  {"SVOй Vайбик"},
		"price": {"1488.0"},
	}

	req, _ := http.NewRequest("POST", "http://127.0.0.1:1111/api/add", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)
	assert.Equal(t, 202, res.StatusCode)

	req, _ = http.NewRequest("GET", "http://127.0.0.1:1111/api/search?name=svoй", nil)
	res, err = app.Test(req, fiber.TestConfig{})

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
}

func TestEditController(t *testing.T) {
	app := services.NewServiceSetup(&db.SQLiteDatabase{})

	data := url.Values{
		"name":  {"SVOй Vайбик"},
		"price": {"1488.0"},
	}

	req, _ := http.NewRequest("POST", "http://127.0.0.1:1111/api/add", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)
	assert.Equal(t, 202, res.StatusCode)

	req, _ = http.NewRequest("GET", "http://127.0.0.1:1111/api/search?name=svoй", nil)
	res, err = app.Test(req, fiber.TestConfig{})

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)

	assert.Nil(t, err)

	var p []entity.Product
	err = json.Unmarshal(body, &p)

	assert.Nil(t, err)

	data = url.Values{
		"id":         {p[0].Id},
		"name":       {"Евгений Пригожин"},
		"definition": {"Всем либерахам кричу я в ответ: УЖЕ БАХМУТА НЕТ!"},
	}

	req, _ = http.NewRequest("PATCH", "http://127.0.0.1:1111/api/edit", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)
	assert.Equal(t, 202, res.StatusCode)
}

func TestDeleteController(t *testing.T) {
	app := services.NewServiceSetup(&db.SQLiteDatabase{})

	data := url.Values{
		"name":  {"SVOй Vайбик"},
		"price": {"1488.0"},
	}

	req, _ := http.NewRequest("POST", "http://127.0.0.1:1111/api/add", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)
	assert.Equal(t, 202, res.StatusCode)

	req, _ = http.NewRequest("GET", "http://127.0.0.1:1111/api/search?name=svoй", nil)
	res, err = app.Test(req, fiber.TestConfig{})

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	body, err := io.ReadAll(res.Body)

	assert.Nil(t, err)

	var p []entity.Product
	err = json.Unmarshal(body, &p)

	assert.Nil(t, err)

	data = url.Values{
		"id": {p[0].Id},
	}

	req, _ = http.NewRequest("DELETE", "http://127.0.0.1:1111/api/remove", strings.NewReader(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = app.Test(req, fiber.TestConfig{})
	assert.Nil(t, err)
	assert.Equal(t, 202, res.StatusCode)

}
