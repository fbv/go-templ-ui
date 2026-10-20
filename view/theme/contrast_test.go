package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckPassesForAThemeHoldingEveryPair(t *testing.T) {
	// A distinctive brand palette, all documented pairs above 4.5:1.
	require.NoError(t, Register(&Theme{
		Name:       "brand-ok",
		Canvas:     "#FAF7F2",
		Surface:    "#FFFFFF",
		SurfaceAlt: "#F0EBE3",
		Field:      "#FFFFFF",
		Subtle:     "#F0EBE3",
		InfoSoft:   "#E8E0D4",
		Ink:        "#1A2B4A",
		InkMuted:   "#3D4A63",
		Action:     "#0B5FA5",
		Success:    "#1F6B3A",
		Warning:    "#8A5A00",
		Danger:     "#A32B22",
	}))

	brand, ok := Get("brand-ok")
	require.True(t, ok)
	assert.Empty(t, Check(brand), "a theme whose pairs meet 4.5:1 must report zero violations")
}

func TestCheckNamesEveryViolatingPair(t *testing.T) {
	// ink too close to canvas: the check names the pair, both colours and the
	// ratio, and does not stop at the first failure.
	theme := &Theme{
		Name:   "brand-fail",
		Canvas: "#FFFFFF",
		Ink:    "#EEEEEE",
	}

	violations := Check(theme)
	require.NotEmpty(t, violations, "ink (#EEEEEE) on canvas (#FFFFFF) must fail")

	var found *Violation
	for i := range violations {
		if violations[i].Pair == "ink on canvas" {
			found = &violations[i]
		}
	}
	require.NotNil(t, found, "the violation must name the pair: %v", violations)
	assert.Equal(t, "ink", found.FGName)
	assert.Equal(t, "#EEEEEE", found.FG, "the violation must carry the foreground colour")
	assert.Equal(t, "canvas", found.BGName)
	assert.Equal(t, "#FFFFFF", found.BG, "the violation must carry the background colour")
	assert.Positive(t, found.Ratio)
	assert.Less(t, found.Ratio, 4.5, "the violation must carry the ratio")
	assert.Contains(t, found.String(), "ink")
	assert.Contains(t, found.String(), "#EEEEEE")
	assert.Contains(t, found.String(), "#FFFFFF")
}

func TestBuiltInThemesMeetTheFloor(t *testing.T) {
	// The same check over the palettes the page paints: light and dark are
	// built-ins and must not regress.
	assert.Empty(t, Check(themeFromColors("light", builtinColors(false))), "built-in light")
	assert.Empty(t, Check(themeFromColors("dark", builtinColors(true))), "built-in dark")
}

func TestBuiltInPalettesResolveFromTokenStyles(t *testing.T) {
	light := builtinColors(false)
	dark := builtinColors(true)
	assert.Equal(t, "#FFFFFF", light["canvas"])
	assert.Equal(t, "#18181B", light["ink"])
	assert.Equal(t, "#1C1B18", dark["canvas"])
	assert.Equal(t, "#EAE8E1", dark["ink"])
	for _, role := range roles {
		assert.NotEmpty(t, light[role], "light must resolve every role, missing %s", role)
		assert.NotEmpty(t, dark[role], "dark must resolve every role, missing %s", role)
	}
}

func TestCheckPartialThemeFallsBackToRoot(t *testing.T) {
	// A partial theme paints its supplied roles over the :root aliases; the
	// check must evaluate the same mix, so an unset surface is the built-in
	// one and a supplied ink is checked against it.
	assert.Empty(t, Check(&Theme{Name: "partial-ok", Ink: "#18181B", InkMuted: "#52525B"}),
		"built-in inks hold every documented pair")

	violations := Check(&Theme{Name: "partial-bad", Canvas: "#FFFFFF", Ink: "#F5F5F5"})
	require.NotEmpty(t, violations, "a partial theme must be checked against the :root fallback")
}

func TestCheckUnreadableColourFailsInsteadOfPanicking(t *testing.T) {
	violations := Check(&Theme{Name: "unreadable", Canvas: "#FFFFFF", Ink: "not-a-colour"})
	require.NotEmpty(t, violations)
	for _, v := range violations {
		assert.Zero(t, v.Ratio, "an unreadable colour must report a zero ratio: %v", v)
	}
}

func TestCheckNilThemeUsesTheBuiltinDefault(t *testing.T) {
	assert.Empty(t, Check(nil), "a nil theme paints the built-in default and must not panic")
}

func TestPairSetMirrorsBadgeContrastFloor(t *testing.T) {
	// ink/ink-muted x {canvas, surface, surface-alt, field, subtle, info-soft}
	// and action/success/warning/danger on surface.
	assert.Len(t, pairs, 16)
	assert.Contains(t, pairs, contrastPair{"ink", "canvas"})
	assert.Contains(t, pairs, contrastPair{"ink", "surface"})
	assert.Contains(t, pairs, contrastPair{"ink", "surface-alt"})
	assert.Contains(t, pairs, contrastPair{"ink", "field"})
	assert.Contains(t, pairs, contrastPair{"ink", "subtle"})
	assert.Contains(t, pairs, contrastPair{"ink", "info-soft"})
	assert.Contains(t, pairs, contrastPair{"ink-muted", "canvas"})
	assert.Contains(t, pairs, contrastPair{"ink-muted", "surface"})
	assert.Contains(t, pairs, contrastPair{"ink-muted", "surface-alt"})
	assert.Contains(t, pairs, contrastPair{"ink-muted", "field"})
	assert.Contains(t, pairs, contrastPair{"ink-muted", "subtle"})
	assert.Contains(t, pairs, contrastPair{"ink-muted", "info-soft"})
	assert.Contains(t, pairs, contrastPair{"action", "surface"})
	assert.Contains(t, pairs, contrastPair{"success", "surface"})
	assert.Contains(t, pairs, contrastPair{"warning", "surface"})
	assert.Contains(t, pairs, contrastPair{"danger", "surface"})
}
