package ui

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThemeToggleScriptDefinesThemeSwitch(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, ThemeToggleScript().Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, "toggleTheme", "must define toggleTheme")
	assert.Contains(t, html, `setAttribute("data-theme"`, "must flip data-theme")
	assert.Contains(t, html, `classList.add("dark")`, "must add dark class")
	assert.Contains(t, html, `classList.remove("dark")`, "must remove dark class")
	assert.Contains(t, html, `localStorage.setItem("theme"`, "must persist the choice")
}

func TestThemeToggleScriptFlipsToOppositeTheme(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, ThemeToggleScript().Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, `=== "dark" ? "light" : "dark"`,
		"toggle must flip dark <-> light")
	assert.Regexp(t, `_themeToggleInitialized`, html, "must guard against double init")
}

func TestThemeToggleControlRenders(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, ThemeToggle().Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, "toggleTheme", "button must invoke toggleTheme")
	assert.Contains(t, html, `aria-label="Toggle light/dark theme"`, "must be labelled")
	assert.Empty(t, hardcodedPaletteRe.FindString(html), "control must not use palette classes")
}

func TestThemeToggleMarkupIsThemeIndependent(t *testing.T) {
	// Matrix row: theme switch flips token values only — markup unchanged.
	// Badge and the toggle take no theme argument, so their class strings are
	// identical in both themes; recoloring is entirely via CSS variables.
	for _, style := range []string{"", "green", "yellow", "red"} {
		got := badgeStyle(&BadgeProps{Name: "x", Style: style})
		assert.NotContains(t, got, "dark:", "style %q must not carry dark: variants", style)
		assert.NotContains(t, got, "data-theme", "style %q must not key off the theme", style)
	}
}
