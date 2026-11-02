package slogx

import (
	"regexp"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// Filter is a function that filters a request based on the given context.
type Filter func(ctx fiber.Ctx) bool

// Accept is a filter that accepts all requests.
func Accept(filter Filter) Filter { return filter }

// Ignore is a filter that ignores all requests.
func Ignore(filter Filter) Filter { return func(ctx fiber.Ctx) bool { return !filter(ctx) } }

// AcceptMethod is a filter that accepts requests with the given HTTP methods.
func AcceptMethod(methods ...string) Filter {
	for i := range methods {
		methods[i] = strings.ToLower(methods[i])
	}

	return func(c fiber.Ctx) bool {
		return slices.Contains(methods, strings.ToLower(c.Method()))
	}
}

// IgnoreMethod is a filter that ignores requests with the given HTTP methods.
func IgnoreMethod(methods ...string) Filter {
	for i := range methods {
		methods[i] = strings.ToLower(methods[i])
	}

	return func(c fiber.Ctx) bool {
		return !slices.Contains(methods, strings.ToLower(c.Method()))
	}
}

// AcceptStatus is a filter that accepts requests with the given status codes.
func AcceptStatus(statuses ...int) Filter {
	return func(c fiber.Ctx) bool {
		return slices.Contains(statuses, c.Response().StatusCode())
	}
}

// IgnoreStatus is a filter that ignores requests with the given status codes.
func IgnoreStatus(statuses ...int) Filter {
	return func(c fiber.Ctx) bool {
		return !slices.Contains(statuses, c.Response().StatusCode())
	}
}

// AcceptStatusGreaterThan is a filter that accepts requests with a status code greater than the given value.
func AcceptStatusGreaterThan(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() > status
	}
}

// AcceptStatusGreaterThanOrEqual is a filter that accepts requests with a status code greater than or equal to the given value.
func AcceptStatusGreaterThanOrEqual(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() >= status
	}
}

// AcceptStatusLessThan is a filter that accepts requests with a status code less than the given value.
func AcceptStatusLessThan(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() < status
	}
}

// AcceptStatusLessThanOrEqual is a filter that accepts requests with a status code less than or equal to the given value.
func AcceptStatusLessThanOrEqual(status int) Filter {
	return func(c fiber.Ctx) bool {
		return c.Response().StatusCode() <= status
	}
}

// IgnoreStatusGreaterThan is a filter that ignores requests with a status code greater than the given value.
func IgnoreStatusGreaterThan(status int) Filter {
	return AcceptStatusLessThanOrEqual(status)
}

// IgnoreStatusGreaterThanOrEqual is a filter that ignores requests with a status code greater than or equal to the given value.
func IgnoreStatusGreaterThanOrEqual(status int) Filter {
	return AcceptStatusLessThan(status)
}

// IgnoreStatusLessThan is a filter that ignores requests with a status code less than the given value.
func IgnoreStatusLessThan(status int) Filter {
	return AcceptStatusGreaterThanOrEqual(status)
}

// IgnoreStatusLessThanOrEqual is a filter that ignores requests with a status code less than or equal to the given value.
func IgnoreStatusLessThanOrEqual(status int) Filter {
	return AcceptStatusGreaterThan(status)
}

// AcceptPath is a filter that accepts requests with a path that matches one of the given values.
func AcceptPath(urls ...string) Filter {
	return func(c fiber.Ctx) bool {
		return slices.Contains(urls, c.Path())
	}
}

// IgnorePath is a filter that ignores requests with a path that matches one of the given values.
func IgnorePath(urls ...string) Filter {
	return func(c fiber.Ctx) bool {
		return !slices.Contains(urls, c.Path())
	}
}

// AcceptPathContains is a filter that accepts requests with a path that contains one of the given values.
func AcceptPathContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Path(), part) {
				return true
			}
		}

		return false
	}
}

// IgnorePathContains is a filter that ignores requests with a path that contains one of the given values.
func IgnorePathContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, part := range parts {
			if strings.Contains(path, part) {
				return false
			}
		}

		return true
	}
}

// AcceptPathPrefix is a filter that accepts requests with a path that starts with one of the given values.
func AcceptPathPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, prefix := range prefixs {
			if strings.HasPrefix(path, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnorePathPrefix is a filter that ignores requests with a path that starts with one of the given values.
func IgnorePathPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, prefix := range prefixs {
			if strings.HasPrefix(path, prefix) {
				return false
			}
		}

		return true
	}
}

// AcceptPathSuffix is a filter that accepts requests with a path that ends with one of the given values.
func AcceptPathSuffix(suffixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, suffix := range suffixs {
			if strings.HasSuffix(path, suffix) {
				return true
			}
		}

		return false
	}
}

// IgnorePathSuffix is a filter that ignores requests with a path that ends with one of the given values.
func IgnorePathSuffix(suffixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, suffix := range suffixs {
			if strings.HasSuffix(path, suffix) {
				return false
			}
		}

		return true
	}
}

// AcceptPathMatch is a filter that accepts requests with a path that matches one of the given regular expressions.
func AcceptPathMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, reg := range regs {
			if reg.MatchString(path) {
				return true
			}
		}

		return false
	}
}

// IgnorePathMatch is a filter that ignores requests with a path that matches one of the given regular expressions.
func IgnorePathMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		path := c.Path()
		for _, reg := range regs {
			if reg.MatchString(path) {
				return false
			}
		}

		return true
	}
}

// AcceptHost is a filter that accepts requests with a host that matches one of the given values.
func AcceptHost(hosts ...string) Filter {
	return func(c fiber.Ctx) bool {
		return slices.Contains(hosts, c.Hostname())
	}
}

// IgnoreHost is a filter that ignores requests with a host that matches one of the given values.
func IgnoreHost(hosts ...string) Filter {
	return func(c fiber.Ctx) bool {
		return !slices.Contains(hosts, c.Hostname())
	}
}

// AcceptHostContains is a filter that accepts requests with a host that contains one of the given parts.
func AcceptHostContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Hostname(), part) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostContains is a filter that ignores requests with a host that contains one of the given parts.
func IgnoreHostContains(parts ...string) Filter {
	return func(c fiber.Ctx) bool {
		for _, part := range parts {
			if strings.Contains(c.Hostname(), part) {
				return false
			}
		}

		return true
	}
}

// AcceptHostPrefix is a filter that accepts requests with a host that starts with one of the given prefixes.
func AcceptHostPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		host := c.Hostname()
		for _, prefix := range prefixs {
			if strings.HasPrefix(host, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostPrefix is a filter that ignores requests with a host that starts with one of the given prefixes.
func IgnoreHostPrefix(prefixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		host := c.Hostname()
		for _, prefix := range prefixs {
			if strings.HasPrefix(host, prefix) {
				return false
			}
		}

		return true
	}
}

// AcceptHostSuffix is a filter that accepts requests with a host that ends with one of the given suffixes.
func AcceptHostSuffix(suffixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		host := c.Hostname()
		for _, suffix := range suffixs {
			if strings.HasSuffix(host, suffix) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostSuffix is a filter that ignores requests with a host that ends with one of the given suffixes.
func IgnoreHostSuffix(suffixs ...string) Filter {
	return func(c fiber.Ctx) bool {
		host := c.Hostname()
		for _, suffix := range suffixs {
			if strings.HasSuffix(host, suffix) {
				return false
			}
		}

		return true
	}
}

// AcceptHostMatch is a filter that accepts requests with a host that matches one of the given regular expressions.
func AcceptHostMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		host := c.Hostname()
		for _, reg := range regs {
			if reg.MatchString(host) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostMatch is a filter that ignores requests with a host that matches one of the given regular expressions.
func IgnoreHostMatch(regs ...regexp.Regexp) Filter {
	return func(c fiber.Ctx) bool {
		host := c.Hostname()
		for _, reg := range regs {
			if reg.MatchString(host) {
				return false
			}
		}

		return true
	}
}
