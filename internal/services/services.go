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
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/sdrapkin/guid"
)

type Servicer interface {
	Add(c fiber.Ctx) error
	Remove(c fiber.Ctx) error
	Edit(c fiber.Ctx) error
	Search(c fiber.Ctx) error
	InitFromFile() error
	WriteToFile() error
	GetMyApp() *fiber.App
	GetAllProducts(c fiber.Ctx) error
	InfoByID(c fiber.Ctx) error
}

type APIService struct {
	data     db.Databaser
	config   config.Configer
	Products []entity.Product
	api      *fiber.App
}

func NewAPIService(cf config.Configer, cf_path string, api *fiber.App, dbase db.Databaser) (Servicer, error) {
	err := cf.ReadConfig(cf_path)
	if err != nil {
		return nil, err
	}
	err = dbase.CreateDBConnection(cf.GetDBPath())
	if err != nil {
		return nil, err
	}

	return &APIService{config: cf, data: dbase, api: api, Products: []entity.Product{}}, nil
}

func NewServiceSetup(dbase db.Databaser) *fiber.App {
	var service_cf config.GoodOldConfig

	app := fiber.New()
	srvs, err := NewAPIService(&service_cf, "appsettings.json", app, dbase)
	if err != nil {
		log.Fatal(err)
	}
	srvs.InitFromFile()
	srvs.GetMyApp().Use(cors.New(
		cors.Config{
			AllowOrigins:        []string{"*"},
			AllowHeaders:        []string{"Origin", "Content-Type", "Accept"},
			AllowMethods:        []string{"GET", "POST", "DELETE", "PATCH"},
			AllowPrivateNetwork: true,
		},
	))

	srvs.GetMyApp().Use("/api", func(c fiber.Ctx) error {
		log.Debugf("API request from %s. URI is: %s", c.Req().IP(), c.Req().RequestCtx().URI().String())

		return c.Next()
	})

	apiGroup := srvs.GetMyApp().Group("/api")

	apiGroup.Get("/search", srvs.Search)
	apiGroup.Get("/get_all", srvs.GetAllProducts)
	apiGroup.Get("/get_id", srvs.InfoByID)
	apiGroup.Post("/add", srvs.Add)
	apiGroup.Patch("/edit", srvs.Edit)
	apiGroup.Delete("/remove", srvs.Remove)

	return app
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

	err = ap.data.AddProduct(&product)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(product)
}

func (ap *APIService) Remove(c fiber.Ctx) error {
	gid := c.FormValue("id")

	rws, err := ap.data.DeleteProduct(gid)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if rws == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	ap.Products = slices.DeleteFunc(ap.Products, func(p entity.Product) bool {
		return p.Id == gid
	})

	return c.SendStatus(fiber.StatusAccepted)
}

func (ap *APIService) Edit(c fiber.Ctx) error {
	var product entity.Product
	gid := c.FormValue("id")
	name := c.FormValue("name")
	definition := c.FormValue("definition")
	price_key := c.FormValue("price")

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

	rws, err := ap.data.UpdateProduct(&product)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if rws == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	ap.Products, _ = ap.data.GetAllProducts()

	return c.JSON(product)
}

func (ap APIService) Search(c fiber.Ctx) error {
	var (
		products []entity.Product
		err      error
	)

	name := c.Query("name")
	if name == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	for _, v := range ap.Products {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(name)) {
			products = append(products, v)
		}
	}

	if len(products) == 0 {
		products, err = ap.data.SearchProduct(name)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if len(products) == 0 {
			return c.SendStatus(fiber.StatusNotFound)
		}
	}

	ap.Products = append(ap.Products, products...)
	return c.JSON(products)
}

func (ap APIService) InfoByID(c fiber.Ctx) error {
	var (
		prd entity.Product
		err error
	)
	id := c.Query("id")
	if id == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if prd, err = ap.data.GetProductByID(id); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(prd)
}

func (ap APIService) WriteToFile() error {
	return ap.data.DropDBAndReinsert(ap.Products)
}

func (ap *APIService) AddGetController(path string, handler func(c fiber.Ctx) error) {
	ap.api.Get(path, handler)
}

func (ap *APIService) AddPostController(path string, handler func(c fiber.Ctx) error) {
	ap.api.Post(path, handler)
}

func (ap *APIService) AddPatchController(path string, handler func(c fiber.Ctx) error) {
	ap.api.Patch(path, handler)
}

func (ap *APIService) AddDeleteController(path string, handler func(c fiber.Ctx) error) {
	ap.api.Delete(path, handler)
}

func (ap *APIService) GetMyApp() *fiber.App {
	return ap.api
}

func (ap APIService) GetAllProducts(c fiber.Ctx) error {
	prds, err := ap.data.GetAllProducts()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if len(prds) == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.JSON(prds)
}
