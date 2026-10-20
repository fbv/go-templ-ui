package theme

import (
	"context"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fbv/go-templ-ui/view/ui"
)

// floor is the WCAG AA text contrast ratio every documented pair must meet.
const floor = 4.5

// Violation is one documented text/background pair below the 4.5:1 floor.
type Violation struct {
	// Pair is the documented pairing, e.g. "ink on canvas".
	Pair string
	// FGName and BGName are the semantic roles of the pair, e.g. "ink" and
	// "canvas".
	FGName string
	BGName string
	// FG and BG are the colours the pair resolves to.
	FG string
	BG string
	// Ratio is the WCAG contrast ratio between FG and BG. A colour the check
	// cannot read reports 0, so an unreadable value fails the check instead of
	// panicking.
	Ratio float64
}

// String renders the violation as e.g.
// "ink (#FFFFFF) on canvas (#18181B) = 17.32".
func (v Violation) String() string {
	return fmt.Sprintf("%s (%s) on %s (%s) = %.2f", v.FGName, v.FG, v.BGName, v.BG, v.Ratio)
}

// contrastPair is one documented pairing of semantic roles.
type contrastPair struct{ fg, bg string }

// pairs is the contrast floor a theme must hold: ink/ink-muted against every
// surface they can sit on, and each status colour against surface. It mirrors
// the set view/ui/badge_test.go enforces on the built-in palettes.
var pairs = []contrastPair{
	{"ink", "canvas"},
	{"ink", "surface"},
	{"ink", "surface-alt"},
	{"ink", "field"},
	{"ink", "subtle"},
	{"ink", "info-soft"},
	{"ink-muted", "canvas"},
	{"ink-muted", "surface"},
	{"ink-muted", "surface-alt"},
	{"ink-muted", "field"},
	{"ink-muted", "subtle"},
	{"ink-muted", "info-soft"},
	{"action", "surface"},
	{"success", "surface"},
	{"warning", "surface"},
	{"danger", "surface"},
}

// Check reports every documented pair of t below the 4.5:1 floor, naming the
// pair, both colours and the ratio. It is the check a consumer runs on their
// theme before shipping it; a nil or partial theme is valid and resolves the
// roles it leaves unset exactly as the page paints them — through the :root
// alias, which is the built-in light palette.
//
// A theme that holds every pair returns an empty slice (nil).
func Check(t *Theme) []Violation {
	return check(effectiveColors(t))
}

// effectiveColors resolves t's roles the way the page paints them: supplied
// values over the :root fallback.
func effectiveColors(t *Theme) map[string]string {
	colors := make(map[string]string, len(roles))
	for role, v := range builtinColors(false) {
		colors[role] = v
	}
	for role, v := range t.colors() {
		colors[role] = v
	}
	return colors
}

func check(colors map[string]string) []Violation {
	var out []Violation
	for _, p := range pairs {
		fg, bg := colors[p.fg], colors[p.bg]
		ratio := contrastRatio(fg, bg)
		if ratio >= floor {
			continue
		}
		out = append(out, Violation{
			Pair:   p.fg + " on " + p.bg,
			FGName: p.fg,
			FG:     fg,
			BGName: p.bg,
			BG:     bg,
			Ratio:  ratio,
		})
	}
	return out
}

var (
	builtinOnce  sync.Once
	builtinLight map[string]string
	builtinDark  map[string]string
)

// builtinColors resolves the built-in light (or dark) palette — the values the
// page paints for a theme that supplies nothing — from the CSS ui.TokenStyles
// emits. The built-ins are read from the stylesheet that owns them, never
// copied into this package, so the two cannot drift.
func builtinColors(dark bool) map[string]string {
	builtinOnce.Do(func() {
		var buf strings.Builder
		_ = ui.TokenStyles().Render(context.Background(), &buf)
		builtinLight = resolvedRoles(buf.String(), false)
		builtinDark = resolvedRoles(buf.String(), true)
	})
	if dark {
		return builtinDark
	}
	return builtinLight
}

var (
	cssBlockRe = regexp.MustCompile(`([^{}]+)\{([^{}]*)\}`)
	cssVarRe   = regexp.MustCompile(`^var\((--[A-Za-z0-9_-]+)\)$`)
	cssComRe   = regexp.MustCompile(`(?s)/\*.*?\*/`)
)

// resolvedRoles returns each semantic role resolved to a concrete value for
// the built-in light palette, or for dark when dark is set.
func resolvedRoles(css string, dark bool) map[string]string {
	blocks := parseCSSBlocks(css)
	vars := map[string]string{}
	for k, v := range blocks[`:root, [data-theme="light"]`] {
		vars[k] = v
	}
	if dark {
		for k, v := range blocks[`[data-theme="dark"]`] {
			vars[k] = v
		}
	}
	for k, v := range blocks[`:root`] {
		vars[k] = v
	}
	out := make(map[string]string, len(roles))
	for _, role := range roles {
		if v, ok := resolveVar(vars, "--color-"+role); ok {
			out[role] = v
		}
	}
	return out
}

func parseCSSBlocks(css string) map[string]map[string]string {
	css = cssComRe.ReplaceAllString(css, "")
	blocks := map[string]map[string]string{}
	for _, m := range cssBlockRe.FindAllStringSubmatch(css, -1) {
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

// resolveVar follows a custom property through its var() references to a
// concrete value. Cycles and missing definitions report false.
func resolveVar(vars map[string]string, name string) (string, bool) {
	seen := map[string]bool{}
	for i := 0; i < 12; i++ {
		v, ok := vars[name]
		if !ok {
			return "", false
		}
		if seen[name] {
			return "", false
		}
		seen[name] = true
		v = strings.TrimSpace(v)
		if m := cssVarRe.FindStringSubmatch(v); m != nil {
			name = m[1]
			continue
		}
		return v, true
	}
	return "", false
}

func srgbToLinear(c float64) float64 {
	if c <= 0.04045 {
		return c / 12.92
	}
	return math.Pow((c+0.055)/1.055, 2.4)
}

// relativeLuminance returns the WCAG relative luminance of a 6-digit hex
// colour, and whether the value could be read at all.
func relativeLuminance(hex string) (float64, bool) {
	h := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(hex), "#"))
	if len(h) != 6 {
		return 0, false
	}
	r, err := strconv.ParseUint(h[0:2], 16, 8)
	if err != nil {
		return 0, false
	}
	g, err := strconv.ParseUint(h[2:4], 16, 8)
	if err != nil {
		return 0, false
	}
	b, err := strconv.ParseUint(h[4:6], 16, 8)
	if err != nil {
		return 0, false
	}
	return 0.2126*srgbToLinear(float64(r)/255) +
		0.7152*srgbToLinear(float64(g)/255) +
		0.0722*srgbToLinear(float64(b)/255), true
}

// contrastRatio returns the WCAG contrast ratio between two colours. A colour
// that is not 6-digit hex yields 0, so an unreadable value fails the check
// instead of panicking.
func contrastRatio(a, b string) float64 {
	la, okA := relativeLuminance(a)
	lb, okB := relativeLuminance(b)
	if !okA || !okB {
		return 0
	}
	hi := math.Max(la, lb)
	lo := math.Min(la, lb)
	return (hi + 0.05) / (lo + 0.05)
}
