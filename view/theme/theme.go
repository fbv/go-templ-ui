// Package theme lets a consumer author a brand theme against the library's
// semantic colour tokens.
//
// A Theme is a plain value: a name and the values of the 31 semantic --color-*
// roles it overrides. Themes live in a package registry and render as one
// [data-theme="name"] block over the aliases ui.TokenStyles declares, so a
// consumer rebrands without touching component markup and without forking the
// built-in light/dark palettes:
//
//	_ = theme.Register(&theme.Theme{
//		Name:   "brand",
//		Canvas: "#FAF7F2",
//		Ink:    "#1A2B4A",
//	})
//
// A role the theme leaves unset keeps its :root alias, so a partial theme is
// valid and fills in over time. Selection reuses data-theme: the pre-paint init
// on layout.Page accepts any registered name and adds Tailwind's dark class
// only for the built-in "dark". Check proves a theme's documented
// text/background pairs meet 4.5:1 before it ships.
//
// A theme is the colour roles only. Radius, spacing and type scales stay
// structural: rebranding is a palette change, not a layout change.
package theme

import (
	"fmt"
	"sort"
	"sync"
)

// Theme is a consumer-authored palette: semantic colour values under a name.
//
// Every field is one semantic role, spelled as the role name without the
// --color- prefix (Canvas is --color-canvas). An empty field means "keep the
// :root alias" and is omitted from the emitted stylesheet, so a partial theme
// is valid. Values are written into the stylesheet as given; they are
// developer-supplied data, like any CSS the consumer owns.
//
// A Theme is a plain, shareable value. Once registered it may be read from any
// goroutine; do not mutate it afterwards.
type Theme struct {
	// Name is the registry key and the data-theme value that selects the
	// theme. It must match [A-Za-z][A-Za-z0-9_-]* so it is usable in a CSS
	// attribute selector and in the pre-paint init's name list, and must not
	// be "light" or "dark" — those name the built-in palettes.
	Name string

	Canvas         string // --color-canvas
	Surface        string // --color-surface
	SurfaceAlt     string // --color-surface-alt
	SurfaceOverlay string // --color-surface-overlay
	Field          string // --color-field
	Subtle         string // --color-subtle
	InfoSoft       string // --color-info-soft
	Border         string // --color-border
	BorderStrong   string // --color-border-strong
	Ink            string // --color-ink
	InkMuted       string // --color-ink-muted
	InkSubtle      string // --color-ink-subtle
	InkDisabled    string // --color-ink-disabled
	Action         string // --color-action
	ActionHover    string // --color-action-hover
	ActionActive   string // --color-action-active
	ActionFg       string // --color-action-fg
	ActionSoft     string // --color-action-soft
	Highlight      string // --color-highlight
	Selected       string // --color-selected
	Success        string // --color-success
	SuccessSoft    string // --color-success-soft
	SuccessBorder  string // --color-success-border
	Warning        string // --color-warning
	WarningSoft    string // --color-warning-soft
	WarningBorder  string // --color-warning-border
	Danger         string // --color-danger
	DangerFg       string // --color-danger-fg
	DangerSoft     string // --color-danger-soft
	DangerBorder   string // --color-danger-border
	Focus          string // --color-focus
}

// roleField binds a semantic role to the Theme field that carries it. The
// slice order is the order ui.TokenStyles declares the aliases in.
var roleField = []struct {
	role string
	get  func(*Theme) *string
}{
	{"canvas", func(t *Theme) *string { return &t.Canvas }},
	{"surface", func(t *Theme) *string { return &t.Surface }},
	{"surface-alt", func(t *Theme) *string { return &t.SurfaceAlt }},
	{"surface-overlay", func(t *Theme) *string { return &t.SurfaceOverlay }},
	{"field", func(t *Theme) *string { return &t.Field }},
	{"subtle", func(t *Theme) *string { return &t.Subtle }},
	{"info-soft", func(t *Theme) *string { return &t.InfoSoft }},
	{"border", func(t *Theme) *string { return &t.Border }},
	{"border-strong", func(t *Theme) *string { return &t.BorderStrong }},
	{"ink", func(t *Theme) *string { return &t.Ink }},
	{"ink-muted", func(t *Theme) *string { return &t.InkMuted }},
	{"ink-subtle", func(t *Theme) *string { return &t.InkSubtle }},
	{"ink-disabled", func(t *Theme) *string { return &t.InkDisabled }},
	{"action", func(t *Theme) *string { return &t.Action }},
	{"action-hover", func(t *Theme) *string { return &t.ActionHover }},
	{"action-active", func(t *Theme) *string { return &t.ActionActive }},
	{"action-fg", func(t *Theme) *string { return &t.ActionFg }},
	{"action-soft", func(t *Theme) *string { return &t.ActionSoft }},
	{"highlight", func(t *Theme) *string { return &t.Highlight }},
	{"selected", func(t *Theme) *string { return &t.Selected }},
	{"success", func(t *Theme) *string { return &t.Success }},
	{"success-soft", func(t *Theme) *string { return &t.SuccessSoft }},
	{"success-border", func(t *Theme) *string { return &t.SuccessBorder }},
	{"warning", func(t *Theme) *string { return &t.Warning }},
	{"warning-soft", func(t *Theme) *string { return &t.WarningSoft }},
	{"warning-border", func(t *Theme) *string { return &t.WarningBorder }},
	{"danger", func(t *Theme) *string { return &t.Danger }},
	{"danger-fg", func(t *Theme) *string { return &t.DangerFg }},
	{"danger-soft", func(t *Theme) *string { return &t.DangerSoft }},
	{"danger-border", func(t *Theme) *string { return &t.DangerBorder }},
	{"focus", func(t *Theme) *string { return &t.Focus }},
}

// roles lists every semantic colour role a theme may supply, in TokenStyles
// order.
var roles = func() []string {
	out := make([]string, len(roleField))
	for i, rf := range roleField {
		out[i] = rf.role
	}
	return out
}()

// colors returns the roles t supplies, keyed by role name; unset roles are
// omitted.
func (t *Theme) colors() map[string]string {
	if t == nil {
		return nil
	}
	out := make(map[string]string, len(roleField))
	for _, rf := range roleField {
		if v := *rf.get(t); v != "" {
			out[rf.role] = v
		}
	}
	return out
}

// declaration is one CSS custom property a theme supplies.
type declaration struct {
	name  string
	value string
}

// declarations returns the --color-* declarations t supplies, in role order;
// unset roles are omitted so CSS keeps the :root alias.
func (t *Theme) declarations() []declaration {
	if t == nil {
		return nil
	}
	out := make([]declaration, 0, len(roleField))
	for _, rf := range roleField {
		if v := *rf.get(t); v != "" {
			out = append(out, declaration{"--color-" + rf.role, v})
		}
	}
	return out
}

// themeFromColors builds a Theme whose roles come from colors. Used to run the
// exported check over palettes that are not registry entries, such as the
// built-in light/dark values.
func themeFromColors(name string, colors map[string]string) *Theme {
	t := &Theme{Name: name}
	for _, rf := range roleField {
		if v, ok := colors[rf.role]; ok {
			*rf.get(t) = v
		}
	}
	return t
}

// validName reports whether name is usable both in a CSS attribute selector
// and in the pre-paint init's name list: a letter, then letters, digits, _ or
// -. Anything else — empty, a leading _ or -, spaces, quotes, escapes — is
// rejected before registration so a registered name can never break the
// emitted CSS or script.
func validName(name string) bool {
	if name == "" {
		return false
	}
	if c := name[0]; (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
		return false
	}
	for i := 1; i < len(name); i++ {
		c := name[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c == '_', c == '-', c >= '0' && c <= '9':
		default:
			return false
		}
	}
	return true
}

var (
	registryMu sync.RWMutex
	registry   = map[string]*Theme{}
)

// Register adds t to the package registry, keyed by t.Name, so Styles emits it
// and layout.Page's pre-paint init accepts the name. It may be called from a
// consumer package's init.
//
// Register returns an error — and leaves the registry untouched — when t is
// nil, when t.Name is not a valid theme name, when t.Name is "light" or "dark"
// (reserved for the built-in palettes), or when the name is already taken. A
// name is taken once: two themes cannot own the same selector.
func Register(t *Theme) error {
	if t == nil {
		return fmt.Errorf("theme: cannot register a nil theme")
	}
	if !validName(t.Name) {
		return fmt.Errorf("theme: invalid name %q: must match [A-Za-z][A-Za-z0-9_-]*", t.Name)
	}
	if t.Name == "light" || t.Name == "dark" {
		return fmt.Errorf("theme: name %q is reserved for the built-in palettes", t.Name)
	}
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, ok := registry[t.Name]; ok {
		return fmt.Errorf("theme: %q is already registered", t.Name)
	}
	registry[t.Name] = t
	return nil
}

// Get returns the registered theme for name. An unknown or empty name reports
// false; the built-in light/dark palettes are not registry entries.
func Get(name string) (*Theme, bool) {
	if name == "" {
		return nil, false
	}
	registryMu.RLock()
	defer registryMu.RUnlock()
	t, ok := registry[name]
	return t, ok
}

// Names returns the registered theme names in sorted order, the list
// layout.Page's pre-paint init accepts beside the built-in light/dark.
func Names() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
