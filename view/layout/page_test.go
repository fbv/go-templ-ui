package layout

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderPage(t *testing.T) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, Page(&PageProps{Title: "test"}).Render(context.Background(), &buf))
	return buf.String()
}

func TestPageAppliesStoredThemeBeforePaint(t *testing.T) {
	html := renderPage(t)

	init := strings.Index(html, `localStorage.getItem("theme")`)
	require.GreaterOrEqual(t, init, 0, "pre-paint theme init script missing")
	headEnd := strings.Index(html, "</head>")
	require.GreaterOrEqual(t, headEnd, 0)
	assert.Less(t, init, headEnd, "theme must be applied in <head>, before paint")

	assert.Contains(t, html, `setAttribute("data-theme", theme)`, "must set data-theme")
	assert.Contains(t, html, `classList.add("dark")`, "must set the dark class for dark")
	assert.Contains(t, html, "--g-canvas", "must inject the token stylesheet")
	assert.Contains(t, html, "--color-surface", "must inject semantic token aliases")
}

func TestPageThemeInitFallsBackToLight(t *testing.T) {
	html := renderPage(t)

	assert.Contains(t, html, `theme !== "dark" && theme !== "light"`,
		"must treat anything but dark/light as unset")
	assert.Contains(t, html, `theme = "light"`, "must fall back to light")
	assert.Contains(t, html, "catch (e)", "must survive localStorage being unavailable")
}

func TestPageRegistersThemeToggleScript(t *testing.T) {
	html := renderPage(t)
	assert.Contains(t, html, "toggleTheme", "page must ship the toggle script")
}

func TestPageCanvasFollowsTokens(t *testing.T) {
	html := renderPage(t)
	assert.Contains(t, html, "bg-canvas", "page canvas must be a token utility")
	assert.NotContains(t, html, `class="bg-white"`, "page canvas must not be hardcoded")
}
