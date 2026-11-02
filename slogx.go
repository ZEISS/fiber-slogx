// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package slogx

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v3"
)

// Config caputes the configuration for running the goth middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	Next func(c fiber.Ctx) bool
	// LogLevel defines the log level for the middleware.
	LogLevel slog.Level
	// Filters defines a list of filters to apply to log records.
	Filters []Filter
}

// DefaultConfig returns the default configuration for the goth middleware.
var DefaultConfig = Config{
	LogLevel: slog.LevelInfo,
}

// New returns a new instance of the goth middleware with the given configuration.
func New(logger *slog.Logger, config ...Config) fiber.Handler {
	cfg := configDefault(config...)

	return func(c fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		start := time.Now()
		path := c.Path()
		query := string(c.Request().URI().QueryString())

		for _, filter := range cfg.Filters {
			if !filter(c) {
				return c.Next()
			}
		}

		commonAttributes := make([]slog.Attr, 0)
		requestAttributes := make([]slog.Attr, 0, 3)

		requestAttributes = append(requestAttributes,
			slog.Time("time", start.UTC()),
			slog.String("path", path),
			slog.String("query", query),
		)

		attributes := append(
			[]slog.Attr{
				{
					Key:   "request",
					Value: slog.GroupValue(requestAttributes...),
				},
			},
			commonAttributes...,
		)

		level := cfg.LogLevel

		logger.LogAttrs(c, level, "incoming request", attributes...)

		return c.Next()
	}
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return DefaultConfig
	}

	// Override default config
	cfg := config[0]

	if cfg.Next == nil {
		cfg.Next = DefaultConfig.Next
	}

	return cfg
}
