package ui

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// hardcodedPaletteRe flags a Tailwind palette class. The colour name may end the
// class token (bg-white, text-black) or continue into a shade/slash segment
// (text-gray-900, bg-white/50), so the suffix is optional and the token boundary
// is checked.
var hardcodedPaletteRe = regexp.MustCompile(`(?:^|[\s"])(?:bg|text|border|from|to|via|ring|outline|fill|stroke)-(?:red|green|blue|yellow|indigo|purple|pink|gray|grey|slate|zinc|neutral|stone|amber|orange|lime|emerald|teal|cyan|sky|violet|fuchsia|rose|white|black)(?:[-/][^\s"]*)?(?:$|[\s"])`)

func TestBadgeBaseUsesComponentTokens(t *testing.T) {
	// The pill radius must come from the L3 component token, never a hardcoded
	// rounded-full, so DESIGN.md's components.badge.radius has a live consumer.
	assert.Contains(t, badgeBase, "rounded-[var(--badge-radius)]")
	assert.NotContains(t, badgeBase, "rounded-full")
	assert.NotContains(t, badgeBase, "rounded-md")
	assert.Empty(t, hardcodedPaletteRe.FindString(badgeBase))
}

func TestBadgeStyleUsesTokenClasses(t *testing.T) {
	cases := []struct {
		style   string
		variant string
	}{
		{"gray", "neutral"},
		{"red", "danger"},
		{"green", "success"},
		{"yellow", "warning"},
		{"indigo", "neutral"},
		{"purple", "neutral"},
		{"blue", "neutral"},
		{"ping", "danger"},
		{"", "neutral"},
		{"not-a-style", "neutral"},
		{"NOT-A-STYLE", "neutral"},
	}
	for _, tc := range cases {
		got := badgeStyle(&BadgeProps{Name: "x", Style: tc.style})
		assert.Contains(t, got, "var(--badge-"+tc.variant+"-bg)", "style %q", tc.style)
		assert.Contains(t, got, "var(--badge-"+tc.variant+"-fg)", "style %q", tc.style)
		assert.Contains(t, got, "var(--badge-"+tc.variant+"-border)", "style %q", tc.style)
		assert.NotContains(t, got, "dark:", "style %q must not use dark: variants", tc.style)
		assert.Empty(t, hardcodedPaletteRe.FindString(got), "style %q has hardcoded palette class: %s", tc.style, got)
	}
}

func TestBadgeStyleFallsBackOnUnknownStyle(t *testing.T) {
	def := badgeStyle(&BadgeProps{Name: "x"})
	unknown := badgeStyle(&BadgeProps{Name: "x", Style: "nope"})
	empty := badgeStyle(&BadgeProps{Name: "x", Style: ""})
	assert.Equal(t, def, unknown)
	assert.Equal(t, def, empty)
	assert.Contains(t, def, "var(--badge-neutral-bg)")
}

func TestBadgeRendersClassWithTokens(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, Badge(&BadgeProps{Name: "Active", Style: "green"}).Render(context.Background(), &buf))
	html := buf.String()
	assert.Contains(t, html, "var(--badge-success-bg)")
	assert.Contains(t, html, "var(--badge-success-fg)")
	assert.Contains(t, html, "var(--badge-success-border)")
	assert.Contains(t, html, "Active")
	assert.Empty(t, hardcodedPaletteRe.FindString(html))
}

func TestTokenContrastMeetsFloor(t *testing.T) {
	css := renderTokenStyles(t)
	for _, dark := range []bool{false, true} {
		vars := themeVars(t, css, dark)
		name := "light"
		if dark {
			name = "dark"
		}

		texts := []string{"--color-ink", "--color-ink-muted"}
		surfaces := []string{
			"--color-canvas", "--color-surface", "--color-surface-alt",
			"--color-field", "--color-subtle", "--color-info-soft",
		}
		for _, fgName := range texts {
			for _, bgName := range surfaces {
				fg := resolveVar(t, vars, fgName)
				bg := resolveVar(t, vars, bgName)
				ratio := contrastRatio(fg, bg)
				assert.GreaterOrEqual(t, ratio, 4.5,
					"%s: %s (%s) on %s (%s) = %.2f", name, fgName, fg, bgName, bg, ratio)
			}
		}

		for _, fgName := range []string{"--color-action", "--color-success", "--color-warning", "--color-danger"} {
			fg := resolveVar(t, vars, fgName)
			bg := resolveVar(t, vars, "--color-surface")
			ratio := contrastRatio(fg, bg)
			assert.GreaterOrEqual(t, ratio, 4.5,
				"%s: %s (%s) on surface (%s) = %.2f", name, fgName, fg, bg, ratio)
		}

		for _, variant := range []string{"neutral", "success", "warning", "danger"} {
			fg := resolveVar(t, vars, "--badge-"+variant+"-fg")
			bg := resolveVar(t, vars, "--badge-"+variant+"-bg")
			ratio := contrastRatio(fg, bg)
			assert.GreaterOrEqual(t, ratio, 4.5,
				"%s: badge %s fg (%s) on bg (%s) = %.2f", name, variant, fg, bg, ratio)
		}
	}
}

func TestTokenLayersPresent(t *testing.T) {
	css := renderTokenStyles(t)
	assert.Contains(t, css, "[data-theme=\"dark\"]")
	assert.Contains(t, css, "[data-theme=\"light\"]")
	vars := themeVars(t, css, false)
	for _, name := range []string{
		"--g-canvas", "--color-canvas", "--color-surface", "--color-ink",
		"--color-action", "--color-danger", "--color-focus",
		"--badge-radius", "--badge-neutral-bg", "--badge-success-fg",
		"--rounded-full", "--text-caption-size",
	} {
		assert.NotEmpty(t, vars[name], "missing token %s", name)
	}
	assert.Equal(t, "#F2F0EA", resolveVar(t, vars, "--color-canvas"))
	dark := themeVars(t, css, true)
	assert.Equal(t, "#1C1B18", resolveVar(t, dark, "--color-canvas"))
	assert.Equal(t, "9999px", resolveVar(t, vars, "--badge-radius"))
}

func renderTokenStyles(t *testing.T) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, TokenStyles().Render(context.Background(), &buf))
	return buf.String()
}

var blockRe = regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`)
var varRefRe = regexp.MustCompile(`^var\((--[A-Za-z0-9_-]+)\)$`)
var commentRe = regexp.MustCompile(`(?s)/\*.*?\*/`)

func parseTokenBlocks(css string) map[string]map[string]string {
	css = commentRe.ReplaceAllString(css, "")
	blocks := map[string]map[string]string{}
	for _, m := range blockRe.FindAllStringSubmatch(css, -1) {
		sel := strings.TrimSpace(m[1])
		if i := strings.LastIndexAny(sel, ";}>"); i >= 0 {
			sel = strings.TrimSpace(sel[i+1:])
		}
		sel = strings.TrimSpace(strings.TrimPrefix(sel, "<style>"))
		decls := map[string]string{}
		for _, d := range strings.Split(m[2], ";") {
			d = strings.TrimSpace(d)
			if !strings.HasPrefix(d, "--") {
				continue
			}
			parts := strings.SplitN(d, ":", 2)
			if len(parts) != 2 {
				continue
			}
			decls[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
		if len(decls) == 0 {
			continue
		}
		if existing, ok := blocks[sel]; ok {
			for k, v := range decls {
				existing[k] = v
			}
		} else {
			blocks[sel] = decls
		}
	}
	return blocks
}

func themeVars(t *testing.T, css string, dark bool) map[string]string {
	t.Helper()
	blocks := parseTokenBlocks(css)
	out := map[string]string{}
	light := blocks[":root, [data-theme=\"light\"]"]
	require.NotEmpty(t, light, "light L1 block not found")
	for k, v := range light {
		out[k] = v
	}
	if dark {
		for k, v := range blocks["[data-theme=\"dark\"]"] {
			out[k] = v
		}
	}
	for k, v := range blocks[":root"] {
		out[k] = v
	}
	return out
}

func resolveVar(t *testing.T, vars map[string]string, name string) string {
	t.Helper()
	seen := map[string]bool{}
	for i := 0; i < 12; i++ {
		v, ok := vars[name]
		require.True(t, ok, "token %s not defined", name)
		require.False(t, seen[name], "cyclic token %s", name)
		seen[name] = true
		v = strings.TrimSpace(v)
		if m := varRefRe.FindStringSubmatch(v); m != nil {
			name = m[1]
			continue
		}
		return v
	}
	t.Fatalf("token %s did not resolve", name)
	return ""
}

func srgbToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

func relativeLuminance(hex string) float64 {
	h := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(hex), "#"))
	requireHexLen(h)
	r, _ := strconv.ParseUint(h[0:2], 16, 8)
	g, _ := strconv.ParseUint(h[2:4], 16, 8)
	b, _ := strconv.ParseUint(h[4:6], 16, 8)
	return 0.2126*srgbToLinear(float64(r)/255) +
		0.7152*srgbToLinear(float64(g)/255) +
		0.0722*srgbToLinear(float64(b)/255)
}

func requireHexLen(h string) {
	if len(h) != 6 {
		panic("expected 6-digit hex color, got " + h)
	}
}

func contrastRatio(a, b string) float64 {
	la := relativeLuminance(a)
	lb := relativeLuminance(b)
	hi := math.Max(la, lb)
	lo := math.Min(la, lb)
	return (hi + 0.05) / (lo + 0.05)
}
