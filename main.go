package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

var score = 0

type GameListItem struct {
	Name            string
	Description     string
	ShowDescription bool
}

var gameListItems = []GameListItem{
	{
		Name:            "Chess",
		Description:     "A two-player strategy board game played on a checkered board with 64 squares arranged in an 8×8 grid.",
		ShowDescription: false,
	},
	{
		Name:            "Checkers",
		Description:     "A group of strategy board games for two players which involve diagonal moves of uniform game pieces and mandatory captures by jumping over opponent pieces.",
		ShowDescription: false,
	},
	{
		Name:            "Go",
		Description:     "An abstract strategy board game for two players in which the aim is to surround more territory than the opponent.",
		ShowDescription: false,
	},
}

func main() {
	log.Printf("Starting Lazy Board Games server ...")

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

		for gameListItem := range gameListItems {
			if gameListItems[gameListItem].Name == clicked {
				gameListItems[gameListItem].ShowDescription = !gameListItems[gameListItem].ShowDescription
			} else {
				gameListItems[gameListItem].ShowDescription = false
			}
		}

		return c.Render("partials/game-list", fiber.Map{"games": gameListItems, "count": len(gameListItems)})
	})

	registerTemplate(app, "pages/games", nil)
	registerTemplate(app, "pages/highscores", nil)
	registerTemplate(app, "pages/profile", nil)

	log.Fatal(app.Listen("localhost:3000"))
}

func registerTemplate(app *fiber.App, template string, bindings fiber.Map) {
	app.Get("/"+template, func(c *fiber.Ctx) error {
		return c.Render(template, bindings)
	})
}
