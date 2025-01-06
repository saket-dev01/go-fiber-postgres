package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/saket-dev01/go-fiber-postgres/models"
	"github.com/saket-dev01/go-fiber-postgres/storage"
	"gorm.io/gorm"
)

// `json:"author"` -> so that go understands to map the fields in the Struct to the json
type Book struct {
	Title 		string		`json:"title"`
	Author 		string		`json:"author"`
	Publisher 	string		`json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func(r *Repository) CreateBook(context *fiber.Ctx) error{
	book := Book{}

	err:=context.BodyParser(&book) // body will be converted to book

	if err!=nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "request failed",
			},
		)
		return err
	}

	errMsg :=r.DB.Create(&book).Error
	if errMsg!=nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "could not create book",
			},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book created",
		},
	)
	return nil
}

func(r *Repository) GetBooks(context *fiber.Ctx) error{
	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err!=nil {
		context.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "could not get books",
			},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "books fetched",
			"data": bookModels,
		},
	)
	return nil
}

func(r *Repository) GetBookById(context *fiber.Ctx) error{
	
	book:=&models.Book{}
	id:=context.Params("id") // body will be converted to book

	
	err:=r.DB.First(&book, id).Error
	if err!=nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "could not find book with given id",
			},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book created",
			"data": book,
		},
	)
	return nil
}

func(r *Repository) DeleteBookById(context *fiber.Ctx) error{
	

	id:=context.Params("id") 

	//fmt.Println(id);
	err:=r.DB.Delete(&models.Book{}, id).Error
	if err!=nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "could not create book",
			},
		)
		return err
	}

	context.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "book deleted",
		},
	)
	return nil
}

func(r *Repository) SetupRoutes(app *fiber.App){
	api:= app.Group("/api")
	api.Post("/book", r.CreateBook)
	api.Delete("/book/:id", r.DeleteBookById)
	api.Get("/book", r.GetBooks)
	api.Get("/book/:id", r.GetBookById)
}

func main(){
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading DB", err)
	}

	config:= &storage.Config{
		Host: os.Getenv("PGHOST"),
		Port: os.Getenv("PGPORT"),
		Password:os.Getenv("PGPASSWORD"),
		User:os.Getenv("PGUSER"),
		SSLMode: os.Getenv("PGSSLMODE"),
		DBName:os.Getenv("PGDATABASE"),
	}
	
	db, err:= storage.NewConnection(config)
	
	if err != nil {
		log.Fatal(err)
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	
	r.SetupRoutes(app)
	//models.MigrateBooks(r.DB)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hola khudse likh lia")
	})

	log.Fatal(app.Listen(":3000"))
	
	//fmt.Println("Hello Go Waasiyon!")
}