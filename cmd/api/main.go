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
	// api := app.Group("/api")
	srvs, err := services.NewAPIService(&service_cf, "appsettings.json", app)
	if err != nil {
		log.Fatal(err)
	}
	srvs.InitFromFile()
	srvs.AddGetController("/api/search", func(c fiber.Ctx) error {
		return srvs.Search(c)
	})
	srvs.AddDeleteController("/api/remove", func(c fiber.Ctx) error {
		return srvs.Remove(c)
	})
	srvs.AddPatchController("/api/edit", func(c fiber.Ctx) error {
		return srvs.Edit(c)
	})
	srvs.AddPostController("/api/add", func(c fiber.Ctx) error {
		return srvs.Add(c)
	})
	srvs.AddGetController("/api/get_all", func(c fiber.Ctx) error {
		return srvs.GetAllProducts(c)
	})

	log.Fatal(srvs.GetMyApp().Listen(":1111"))
}
