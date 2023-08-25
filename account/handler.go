package account

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

type Account struct {
	ID       string
	Email    string
	Nickname string
	Password string
}

func Register(app *fiber.App) {
	app.Get("/account", func(c *fiber.Ctx) error {
		return c.Render("account", nil)
	})

	app.Get("/create-account", func(c *fiber.Ctx) error {
		return c.Render("create-account", nil)
	})

	app.Post("/create-account", handleCreateAccount)
}

func handleCreateAccount(c *fiber.Ctx) error {
	newAccount := Account{}

	err := c.BodyParser(&newAccount)
	if err != nil {
		log.Println("Can't parse request body:", err)
		return c.SendString("ERROR")
	}

	log.Printf("Create new account: %v\n", newAccount)

	return c.SendString("DONE")
}
