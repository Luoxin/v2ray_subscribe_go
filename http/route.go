package http

import (
	"github.com/gofiber/fiber/v2"
)

func registerRouting(app *fiber.App) error {
	app.Get("version", Version)
	app.Post("version", Version)

	return nil
}
