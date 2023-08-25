package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

var score = 0

type GameListItem struct {
	ID              string
	Name            string
	Description     string
	ShowDescription bool
}

type Account struct {
	ID       string
	Email    string
	Nickname string
	Password string
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	hostAndPort := host + ":" + port

	log.Printf("Starting Lazy Board Games server (on %s) ...", hostAndPort)

	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgDatabase := os.Getenv("PG_DATABASE")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s", pgHost, pgPort, pgUser, pgDatabase)
	log.Print("Connecting to database: ", connectionString)

	connectionString += " password=" + pgPassword
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal("Can't open database connection", err)
	}
	defer db.Close()

	engine := django.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil)
	})

	app.Get("/menu", func(c *fiber.Ctx) error {
		return c.Render("partials/menu", nil)
	})

	app.Get("/highscores/increment", func(c *fiber.Ctx) error {
		score++
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/highscores/decrement", func(c *fiber.Ctx) error {
		score--
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/highscores/reset", func(c *fiber.Ctx) error {
		score = 0
		return c.SendString(strconv.Itoa(score))
	})

	app.Get("/wait", func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second)
		return c.SendString("1337")
	})

	app.Get("/game-list", func(c *fiber.Ctx) error {
		clicked := c.Query("clicked")

		rows, err := db.Query("select id, name, description from games_metadata")
		if err != nil {
			return c.SendString(err.Error())
		}

		var gameListItems []GameListItem

		for rows.Next() {
			var gameListItem GameListItem
			err = rows.Scan(&gameListItem.ID, &gameListItem.Name, &gameListItem.Description)

			if err != nil {
				return c.SendString(err.Error())
			}

			gameListItems = append(gameListItems, gameListItem)
		}

		for gameListItem := range gameListItems {
			if gameListItems[gameListItem].Name == clicked {
				gameListItems[gameListItem].ShowDescription = !gameListItems[gameListItem].ShowDescription
			} else {
				gameListItems[gameListItem].ShowDescription = false
			}
		}

		return c.Render("partials/game-list", fiber.Map{"games": gameListItems, "count": len(gameListItems)})
	})

	app.Post("/create-account", func(c *fiber.Ctx) error {
		newAccount := Account{}

		err := c.BodyParser(&newAccount)
		if err != nil {
			log.Println("Can't parse request body:", err)
			return c.SendString("ERROR")
		}

		log.Printf("Create new account: %v\n", newAccount)

		return c.SendString("DONE")
	})

	registerTemplate(app, "pages/games", nil)
	registerTemplate(app, "pages/highscores", nil)
	registerTemplate(app, "pages/account", nil)
	registerTemplate(app, "pages/create-account", nil)

	log.Fatal(app.Listen(hostAndPort))
}

func registerTemplate(app *fiber.App, template string, bindings fiber.Map) {
	app.Get("/"+template, func(c *fiber.Ctx) error {
		return c.Render(template, bindings)
	})
}
