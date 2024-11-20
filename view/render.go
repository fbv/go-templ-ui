package view

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

const DateTime = "2006-01-02 15:04:05"

type pathType string

var pathKey = pathType("path")

// PathPrefixIs check if current path has prefix
func PathPrefixIs(ctx context.Context, prefix string) bool {
	path := GetPath(ctx)
	return strings.HasPrefix(path, prefix)
}

func PathIs(ctx context.Context, prefix string) bool {
	path := GetPath(ctx)
	return path == prefix
}

func GetPath(ctx context.Context) string {
	if path, ok := ctx.Value(pathKey).(string); ok {
		return path
	}
	return ""
}

// Render renders page using Templ
func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	path := c.Path()
	c.Context().SetUserValue(pathKey, path)
	c.Set("Content-Type", "text/html")
	componentHandler := templ.Handler(component, options...)
	return adaptor.HTTPHandler(componentHandler)(c)
}

func FormatJSON(s string) (string, error) {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(s), "", "  "); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func IntNotZero(i int) string {
	if i > 0 {
		return strconv.Itoa(i)
	}
	return ""
}

func StaticPage(component templ.Component) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return Render(c, component)
	}
}
