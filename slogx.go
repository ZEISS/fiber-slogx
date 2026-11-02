// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package slogx

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// The contextKey type is unexported to prevent collisions with context keys defined in
// other packages.
type contextKey int

// The keys for the values in context.
const (
	requestIDKey contextKey = iota
)

const (
	requestIDHeaderKey = "X-Request-Id"
)

// Config caputes the configuration for running the goth middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	Next func(c fiber.Ctx) bool
	// LogLevel defines the log level for the middleware.
	LogLevel slog.Level
	// Filters defines a list of filters to apply to log records.
	Filters []Filter
	// RequestHeaders defines the list of request headers to include in log records.
	RequestHeaders map[string]struct{}
}

// DefaultConfig returns the default configuration for the goth middleware.
var DefaultConfig = Config{
	LogLevel:       slog.LevelInfo,
	RequestHeaders: RequestHeaders,
}

// RequestHeaders defines the list of request headers to include in log records.
var RequestHeaders = map[string]struct{}{
	"accept":          {},
	"accept-encoding": {},
	"accept-language": {},
	"connection":      {},
	"host":            {},
	"priority":        {},
	"user-agent":      {},
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

		rid := c.Get(requestIDHeaderKey, uuid.New().String())
		fiber.StoreInContext(c, requestIDKey, rid)
		c.Set(requestIDHeaderKey, rid)

		for _, filter := range cfg.Filters {
			if !filter(c) {
				return c.Next()
			}
		}

		commonAttributes := make([]slog.Attr, 0, 1)
		requestAttributes := make([]slog.Attr, 0, 4)

		commonAttributes = append(commonAttributes,
			slog.String("request_id", rid),
		)

		requestAttributes = append(requestAttributes,
			slog.Time("time", start.UTC()),
			slog.String("path", path),
			slog.String("query", query),
		)

		kv := []any{}
		for k, v := range c.GetReqHeaders() {
			if _, found := cfg.RequestHeaders[strings.ToLower(k)]; found {
				kv = append(kv, slog.Any(k, v))
			}
		}
		requestAttributes = append(requestAttributes, slog.Group("header", kv...))

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

	if cfg.RequestHeaders == nil {
		cfg.RequestHeaders = DefaultConfig.RequestHeaders
	}

	return cfg
}
