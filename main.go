package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
)

var pages = []string{"pages/profile", "pages/games", "pages/highscores"}

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
