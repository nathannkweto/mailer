package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/nathannkweto/mailer/applog"
	"github.com/nathannkweto/mailer/concurrency"
	"github.com/nathannkweto/mailer/config"
	"github.com/nathannkweto/mailer/routes"
)

func main() {
	// load .env if present (no error if missing)
	_ = godotenv.Load()

	// load packages
	cfg := config.Load()
	applog.InitLogger(cfg.EnvLogLevel)
	concurrency.Init(cfg.MaxConcurrentSends)

	// start app
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// request logger
	app.Use(fiberlogger.New())

	// register routes
	routes.RegisterRoutes(app, cfg)

	// server configuration (testing)
	port := cfg.Port
	if port == "" {
		port = "3000"
	}

	// start server with parameters
	applog.Log.Infof("starting email-service on :%s (max_concurrent=%d)", port, cfg.MaxConcurrentSends)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		applog.Log.Fatalf("failed to start server: %v", err)
		os.Exit(1)
	}
}
