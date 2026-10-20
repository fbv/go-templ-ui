package theme

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func renderStyles(t *testing.T) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, Styles().Render(context.Background(), &buf))
	return buf.String()
}

func stylesBlock(t *testing.T, css, name string) string {
	t.Helper()
	start := strings.Index(css, `[data-theme="`+name+`"]`)
	require.GreaterOrEqual(t, start, 0, "block for %s missing", name)
	end := strings.Index(css[start:], "}")
	require.GreaterOrEqual(t, end, 0)
	return css[start : start+end]
}

// valueFor is a distinctive 6-digit hex colour per role index.
func valueFor(i int) string {
	return fmt.Sprintf("#%06x", 0x101010+i*0x030303)
}

func TestRegisterRejectsInvalidNames(t *testing.T) {
	// A name unusable in a CSS attribute selector or in the pre-paint init's
	// name list must never enter the registry.
	for _, name := range []string{"", " ", "my theme", `br"and`, `br\and`, "1brand", "_brand", "-brand", "brand\ndrop", "brand;drop", "<script>"} {
		err := Register(&Theme{Name: name, Canvas: "#FFFFFF"})
		require.Error(t, err, "name %q must be rejected", name)
		_, ok := Get(name)
		assert.False(t, ok, "name %q must never enter the registry", name)
	}
	require.Error(t, Register(nil), "a nil theme must be rejected")
}

func TestRegisterRejectsReservedNames(t *testing.T) {
	// "light" and "dark" name the built-in palettes; a registered theme must
	// never override their [data-theme=...] blocks in the cascade.
	for _, name := range []string{"light", "dark"} {
		err := Register(&Theme{Name: name, Canvas: "#FFFFFF"})
		require.Error(t, err, "name %q is reserved and must be rejected", name)
		assert.Contains(t, err.Error(), "reserved", "the error must say the name is reserved")
		_, ok := Get(name)
		assert.False(t, ok, "name %q must never enter the registry", name)
	}
}

func TestRegisterValidNames(t *testing.T) {
	for _, name := range []string{"a", "brand", "Brand2", "brand-2", "brand_2", "x-y-z"} {
		require.NoError(t, Register(&Theme{Name: name}), "name %q must be accepted", name)
		got, ok := Get(name)
		require.True(t, ok, "name %q must be registered", name)
		assert.Equal(t, name, got.Name)
	}
}

func TestRegisterAndSelect(t *testing.T) {
	require.NoError(t, Register(&Theme{
		Name:   "brand-select",
		Canvas: "#FAF7F2",
		Ink:    "#1A2B4A",
	}))

	got, ok := Get("brand-select")
	require.True(t, ok, "a registered theme must be retrievable by name")
	assert.Equal(t, "#FAF7F2", got.Canvas)
	assert.Equal(t, "#1A2B4A", got.Ink)

	// Selection reuses data-theme: the page CSS must carry a block for the
	// name that declares the supplied values.
	css := renderStyles(t)
	assert.Contains(t, css, `[data-theme="brand-select"]`)
	assert.Contains(t, css, "--color-canvas: #FAF7F2;")
	assert.Contains(t, css, "--color-ink: #1A2B4A;")
}

func TestRegisterDuplicateIsRejectedAndRegistryUnchanged(t *testing.T) {
	require.NoError(t, Register(&Theme{Name: "dup", Canvas: "#111111"}))

	err := Register(&Theme{Name: "dup", Canvas: "#222222"})
	require.Error(t, err, "a second Register of a taken name must fail")
	assert.Contains(t, err.Error(), "dup", "the error must name the conflict")

	got, ok := Get("dup")
	require.True(t, ok)
	assert.Equal(t, "#111111", got.Canvas, "the registry must be unchanged")
	assert.Equal(t, 1, strings.Count(renderStyles(t), `[data-theme="dup"]`),
		"the stylesheet must carry one block for the name")
}

func TestUnknownSelectionFallsBackToRoot(t *testing.T) {
	// data-theme naming no registered theme produces no block, so the :root
	// aliases paint the built-in default: no error, no panic.
	_, ok := Get("nope")
	assert.False(t, ok)
	assert.NotContains(t, renderStyles(t), `[data-theme="nope"]`)
}

func TestPartialThemeOmitsUnsetKeys(t *testing.T) {
	require.NoError(t, Register(&Theme{Name: "partial", Canvas: "#FAF7F2"}))

	block := stylesBlock(t, renderStyles(t), "partial")
	assert.Contains(t, block, "--color-canvas: #FAF7F2;")
	assert.NotContains(t, block, "--color-ink", "an unset role must be omitted so :root wins")
	assert.NotContains(t, block, "--color-surface")
	assert.NotContains(t, block, "--rounded-", "a theme is the colour roles only")
	assert.NotContains(t, block, "--space-")
	assert.NotContains(t, block, "--text-")
	assert.NotContains(t, block, "--font-")
}

func TestStylesDeclaresEverySuppliedRole(t *testing.T) {
	full := &Theme{Name: "full"}
	for i, rf := range roleField {
		*rf.get(full) = valueFor(i)
	}
	require.NoError(t, Register(full))

	block := stylesBlock(t, renderStyles(t), "full")
	for i, rf := range roleField {
		assert.Contains(t, block, "--color-"+rf.role+": "+valueFor(i)+";",
			"role %s must be declared", rf.role)
	}
	assert.Equal(t, 31, len(roleField), "a theme is the 31 --color-* roles only")
}

func TestStylesEmitsOneBlockPerRegisteredTheme(t *testing.T) {
	require.NoError(t, Register(&Theme{Name: "one", Canvas: "#010101"}))
	require.NoError(t, Register(&Theme{Name: "two", Canvas: "#020202"}))

	css := renderStyles(t)
	assert.Equal(t, 1, strings.Count(css, `[data-theme="one"]`))
	assert.Equal(t, 1, strings.Count(css, `[data-theme="two"]`))
	assert.Less(t, strings.Index(css, `[data-theme="one"]`), strings.Index(css, `[data-theme="two"]`),
		"blocks must be emitted in sorted name order")
}

func TestNamesListsRegisteredThemesSorted(t *testing.T) {
	require.NoError(t, Register(&Theme{Name: "zeta", Canvas: "#010101"}))
	require.NoError(t, Register(&Theme{Name: "alpha", Canvas: "#020202"}))

	names := Names()
	assert.Contains(t, names, "zeta")
	assert.Contains(t, names, "alpha")
	for i := 1; i < len(names); i++ {
		assert.Less(t, names[i-1], names[i], "Names must be sorted: %v", names)
	}
}

func TestGetUnknownAndEmpty(t *testing.T) {
	_, ok := Get("")
	assert.False(t, ok)
	_, ok = Get("never-registered")
	assert.False(t, ok, "the built-in light/dark palettes are not registry entries")
}
