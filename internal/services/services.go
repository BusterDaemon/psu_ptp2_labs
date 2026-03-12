package services

import (
	"encoding/json"
	"os"
	"product_api/entity"
	"product_api/internal/config"
	"product_api/internal/db"
	"product_api/internal/imgconv"
	"slices"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/sdrapkin/guid"
)

type Servicer interface {
	Add(c fiber.Ctx) error
	Remove(c fiber.Ctx) error
	Edit(c fiber.Ctx) error
	Search(c fiber.Ctx) error
	InitFromFile() error
	WriteToFile() error
}

type APIService struct {
	data     db.Databaser
	config   config.Configer
	Products []entity.Product
}

func NewAPIService(cf config.Configer, cf_path string) (Servicer, error) {
	var database db.Database
	err := cf.ReadConfig(cf_path)
	if err != nil {
		return nil, err
	}
	err = database.CreateDBConnection(cf.GetDBPath())
	if err != nil {
		return nil, err
	}

	return &APIService{config: cf, data: &database}, nil
}

func (ap *APIService) InitFromFile() error {
	f, err := os.ReadFile(ap.config.GetInitFilePath())
	if err != nil {
		return err
	}

	err = json.Unmarshal(f, &ap.Products)
	if err != nil {
		return err
	}

	return nil
}

func (ap *APIService) Add(c fiber.Ctx) error {
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

	ap.Products = append(ap.Products, product)

	// err = ap.data.AddProduct(&product)
	// if err != nil {
	// 	log.Errorw("Can't add record to database", err)
	// 	return c.SendStatus(fiber.StatusInternalServerError)
	// }

	err = ap.WriteToFile()
	if err != nil {
		log.Error("Невозможно произвести запись в базу данных", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func (ap *APIService) Remove(c fiber.Ctx) error {
	gid := c.FormValue("id")

	ap.Products = slices.DeleteFunc(ap.Products, func(p entity.Product) bool {
		return p.Id == gid
	})

	// err := ap.data.DeleteProduct(gid)
	// if err != nil {
	// 	log.Errorw("Can't delete record from database:", err)
	// 	return c.SendStatus(fiber.StatusInternalServerError)
	// }

	return c.SendStatus(fiber.StatusAccepted)
}

func (ap *APIService) Edit(c fiber.Ctx) error {
	var product entity.Product
	gid := c.FormValue("id")
	name := c.FormValue("name")
	definition := c.FormValue("definition")
	price_key := c.FormValue("price")

	// product, err := ap.data.GetProductByID(gid)
	product = func() entity.Product {
		for _, v := range ap.Products {
			if v.Id == gid {
				return v
			}
		}
		return entity.Product{}
	}()

	if product.Id == "" {
		return c.SendStatus(fiber.StatusNotFound)
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

	// ap.data.UpdateProduct(&product)
	for _, v := range ap.Products {
		if v.Id == gid {
			v = product
			break
		}
	}

	err = ap.WriteToFile()
	if err != nil {
		log.Error("Невозможно произвести запись в базу данных", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusAccepted)
}

func (ap APIService) Search(c fiber.Ctx) error {
	var products []entity.Product

	name := c.Query("name")

	for _, v := range ap.Products {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			products = append(products, v)
		}
	}
	// products, err := ap.data.SearchProduct(name)
	// if err != nil {
	// 	switch err {
	// 	case gorm.ErrRecordNotFound:
	// 		return c.SendStatus(fiber.StatusNotFound)
	// 	default:
	// 		log.Errorw("Can't retrieve records from database:", err)
	// 		return c.SendStatus(fiber.StatusInternalServerError)
	// 	}
	// }

	if len(products) == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(products)
}

func (ap APIService) WriteToFile() error {
	return ap.data.DropDBAndReinsert(ap.Products)
}
