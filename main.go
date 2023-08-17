package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"log"
)

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

	log.Fatal(app.Listen("localhost:3000"))
}
