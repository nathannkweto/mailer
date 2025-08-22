package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/nathannkweto/mailer/config"
	"github.com/nathannkweto/mailer/handlers"
)

// RegisterRoutes attaches endpoints to the Fiber app
func RegisterRoutes(app *fiber.App, cfg *config.Config) {
	app.Post("/send", handlers.SendEmailHandler(cfg))
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]any{"ok": true, "message": "healthy"})
	})
}
