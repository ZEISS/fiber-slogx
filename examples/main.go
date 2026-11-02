package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	safe "github.com/gofiber/fiber/v3/middleware/recover"
	slogx "github.com/zeiss/fiber-slogx"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Add an attribute to all log entries made through this logger.
	logger = logger.With("env", "production")

	app := fiber.New()

	app.Use(slogx.New(logger.WithGroup("http")))
	// config := slogfiber.Config{WithRequestBody: true, WithResponseBody: true, WithRequestHeader: true, WithResponseHeader: true}
	// app.Use(slogfiber.NewWithConfig(logger, config))
	app.Use(safe.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}

	// output:
	// time=2026-11-02T22:17:33.018+01:00 level=INFO msg="incoming request" env=production http.request.time=2026-11-02T21:17:33.018Z http.request.path=/ http.request.query=""
}
