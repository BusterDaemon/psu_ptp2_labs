package main

import (
	"log"
	"product_api/internal/services"
)

func main() {
	app := services.NewServiceSetup()

	log.Fatal(app.Listen(":1111"))

}
