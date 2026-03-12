package main

import (
	"product_api/internal/config"
	"product_api/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	var (
		service_cf config.GoodOldConfig
	)

	app := fiber.New()
	api := app.Group("/api")
	srvs, err := services.NewAPIService(&service_cf, "appsettings.json")
	if err != nil {
		log.Fatal(err)
	}
	srvs.InitFromFile()

	// POST Add
	api.Post("/add", func(c fiber.Ctx) error {
		// return c.SendStatus(fiber.StatusAccepted)
		return srvs.Add(c)
	})

	// DELETE Remove
	api.Delete("/remove", func(c fiber.Ctx) error {
		// return c.SendStatus(fiber.StatusAccepted)
		return srvs.Remove(c)
	})

	// PATCH Edit
	api.Patch("/edit", func(c fiber.Ctx) error {
		// return c.SendStatus(fiber.StatusAccepted)
		return srvs.Edit(c)
	})

	// GET Search
	api.Get("/search", func(c fiber.Ctx) error {
		// return c.SendStatus(fiber.StatusAccepted)
		return srvs.Search(c)
	})

	log.Fatal(app.Listen(":1111"))
}
