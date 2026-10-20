package layout

import (
	"context"
	"strings"
	"testing"

	"github.com/fbv/go-templ-ui/view/i18n"
	"github.com/fbv/go-templ-ui/view/theme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderPage(t *testing.T) string {
	t.Helper()
	return renderPageWithContext(t, context.Background())
}

func renderPageWithContext(t *testing.T, ctx context.Context) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, Page(&PageProps{Title: "test"}).Render(ctx, &buf))
	return buf.String()
}

func TestPageDefaultsToEnglishDocumentLanguage(t *testing.T) {
	html := renderPage(t)
	assert.Contains(t, html, `<html lang="en">`, "default context must render an English document")
}

func TestPageTakesDocumentLanguageFromLocale(t *testing.T) {
	html := renderPageWithContext(t, i18n.WithLocale(context.Background(), i18n.Russian))
	assert.Contains(t, html, `<html lang="ru">`, "a Russian locale must render a Russian document")
	assert.NotContains(t, html, `<html lang="en">`, "the document language must not stay hardcoded")
}

func TestPageFallsBackToEnglishForATaglessLocale(t *testing.T) {
	html := renderPageWithContext(t, i18n.WithLocale(context.Background(), &i18n.Locale{Name: "No tag"}))
	assert.Contains(t, html, `<html lang="en">`, "a carried locale without a tag must not produce an empty lang")
	assert.NotContains(t, html, `<html lang="">`)
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

func TestPageThemeInitAcceptsRegisteredNamesAndFallsBackToLight(t *testing.T) {
	html := renderPage(t)

	assert.Contains(t, html, `var known = ["light", "dark"`,
		"the init must accept the built-in names")
	assert.Contains(t, html, `known.indexOf(theme) < 0`,
		"must treat any name outside the known list as unset")
	assert.Contains(t, html, `theme = "light"`, "must fall back to light")
	assert.Contains(t, html, "catch (e)", "must survive localStorage being unavailable")
}

func TestPageRendersRegisteredCustomTheme(t *testing.T) {
	require.NoError(t, theme.Register(&theme.Theme{
		Name:   "brand",
		Canvas: "#FAF7F2",
		Ink:    "#1A2B4A",
	}))

	html := renderPage(t)

	// data-theme carries the registered name with its values.
	assert.Contains(t, html, `data-theme="brand"`, "the page must carry the selected name")
	assert.Contains(t, html, "--color-canvas: #FAF7F2;", "the page must carry the supplied values")
	assert.Contains(t, html, "--color-ink: #1A2B4A;")

	// The pre-paint init accepts the registered name.
	assert.Contains(t, html, `"light", "dark", "brand"`, "the init must accept registered names")

	// dark is absent for a custom theme; only "dark" gains the class.
	assert.Contains(t, html, `if (theme === "dark") root.classList.add("dark")`,
		"the dark class must be added only for dark")
	assert.NotContains(t, html, `class="dark"`, "a custom theme must not ship the dark class")

	// light and dark are unchanged.
	assert.Contains(t, html, `:root, [data-theme="light"]`, "built-in light must be present")
	assert.Contains(t, html, `[data-theme="dark"]`, "built-in dark must be present")
	assert.Contains(t, html, "--g-canvas: #FFFFFF;", "built-in light values must be unchanged")
	assert.Contains(t, html, "--g-canvas: #1C1B18;", "built-in dark values must be unchanged")
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
