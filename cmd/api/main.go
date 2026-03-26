package main

import (
	"log"
	"product_api/internal/db"
	"product_api/internal/services"
)

func main() {
	app := services.NewServiceSetup(&db.SQLiteDatabase{})

	log.Fatal(app.Listen(":1111"))

}
