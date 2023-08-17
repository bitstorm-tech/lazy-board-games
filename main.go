package main

import (
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

var pages = []string{"pages/profile", "pages/games", "pages/highscores"}
var score = 0

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

	registerTemplate(app, "pages/games")
	registerTemplate(app, "pages/highscores")
	registerTemplate(app, "pages/profile")

	log.Fatal(app.Listen("localhost:3000"))
}

func registerTemplate(app *fiber.App, template string) {
	app.Get("/"+template, func(c *fiber.Ctx) error {
		return c.Render(template, nil)
	})
}
