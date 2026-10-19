package view

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view/i18n"
	"github.com/fbv/go-templ-ui/view/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// renderThroughFiber renders component through view.Render, the repo's real
// render path, and returns the response body.
func renderThroughFiber(t *testing.T, component templ.Component) string {
	t.Helper()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, component)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/", nil))
	require.NoError(t, err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, 200, resp.StatusCode)
	return string(body)
}

// A consumer selects a locale by wrapping the component; fasthttp's RequestCtx
// is the templ context, so the wrapper's locale reaches every component below.
func TestRenderCarriesLocaleThroughTheRenderPath(t *testing.T) {
	russian := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return ui.YesNo(true).Render(i18n.WithLocale(ctx, i18n.Russian), w)
	})

	body := renderThroughFiber(t, russian)
	assert.Contains(t, body, "Да", "a wrapped Russian locale must reach the rendered component")
	assert.NotContains(t, body, ">Yes<")
}

func TestRenderDefaultsToEnglishWithoutALocale(t *testing.T) {
	body := renderThroughFiber(t, ui.YesNo(true))
	assert.Contains(t, body, ">Yes<", "the default render path keeps English")
	assert.NotContains(t, body, "Да")
}
