package main

import (
	"product_api/entity"
	"product_api/internal/db"
	"product_api/internal/imgconv"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/sdrapkin/guid"
	"gorm.io/gorm"
)

func main() {
	app := fiber.New()
	api := app.Group("/api")
	var deber db.Database
	err := deber.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	// POST Add
	api.Post("/add", func(c fiber.Ctx) error {
		definition := c.FormValue("definition")
		name := c.FormValue("name")
		price_key := c.FormValue("price")
		image, err := c.FormFile("image")
		var b64img string

		if image != nil {
			if err != nil {
				log.Errorw("Can't get image from form", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			b64img, err = imgconv.ConvertToBase64(image)
			if err != nil {
				log.Errorw("Can't convert image to base64", err)
				return c.SendStatus(fiber.StatusBadRequest)
			}
		}

		price, err := strconv.ParseFloat(price_key, 32)
		if err != nil {
			log.Errorw("Can't parse price", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		product := entity.Product{
			Id:         guid.NewString(),
			Definition: definition,
			Name:       name,
			Price:      float32(price),
			Image:      b64img,
		}
		err = deber.AddProduct(&product)
		if err != nil {
			log.Errorw("Can't add record to database", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusAccepted)
	})

	// DELETE Remove
	api.Delete("/remove", func(c fiber.Ctx) error {
		gid := c.FormValue("id")

		err := deber.DeleteProduct(gid)
		if err != nil {
			log.Errorw("Can't delete record from database:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusAccepted)
	})

	// PATCH Edit
	api.Patch("/edit", func(c fiber.Ctx) error {
		gid := c.FormValue("id")
		name := c.FormValue("name")
		definition := c.FormValue("definition")
		price_key := c.FormValue("price")

		product, err := deber.GetProductByID(gid)
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return c.SendStatus(fiber.StatusNotFound)
			default:
				log.Errorw("Can't update the record:", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		image, err := c.FormFile("image")
		if image != nil {
			if err != nil {
				log.Errorw("Can't get image from form", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			b64img, err := imgconv.ConvertToBase64(image)
			if err != nil {
				log.Errorw("Can't convert image to base64", err)
				return c.SendStatus(fiber.StatusBadRequest)
			}
			product.Image = b64img
		}

		if name != "" {
			product.Name = name
		}
		if definition != "" {
			product.Definition = definition
		}
		if price_key != "" {
			price, err := strconv.ParseFloat(price_key, 32)
			if err != nil {
				log.Errorw("Can't parse price", err)
				return c.SendStatus(fiber.StatusBadRequest)
			}
			product.Price = float32(price)
		}

		deber.UpdateProduct(&product)

		return c.SendStatus(fiber.StatusAccepted)
	})

	// GET Search
	api.Get("/search", func(c fiber.Ctx) error {
		name := c.Query("name")

		products, err := deber.SearchProduct(name)
		if err != nil {
			switch err {
			case gorm.ErrRecordNotFound:
				return c.SendStatus(fiber.StatusNotFound)
			default:
				log.Errorw("Can't retrieve records from database:", err)
				return c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		if len(products) == 0 {
			return c.SendStatus(fiber.StatusNotFound)
		}

		return c.JSON(products)
	})

	log.Fatal(app.Listen(":1111"))
}
