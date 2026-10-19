package ui

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Chrome token tests. The eight chrome components (headings, link, code, block,
// card, panel, breadcrumb, sidebar) must take every colour from the L2/L3 token
// layers, emit identical markup in both themes, and keep the text/fill pairs the
// matrix promises at >= 4.5:1.

func renderComponent(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf strings.Builder
	require.NoError(t, c.Render(context.Background(), &buf))
	return buf.String()
}

// chromeSurface is every class string the chrome slice can emit as an element
// class, so one pass can prove "no palette class, no dark: variant" for all of
// them without duplicating a colour table.
func chromeClassSources(t *testing.T) map[string][]string {
	t.Helper()
	sources := map[string][]string{
		"H1":   {renderComponent(t, H1("t"))},
		"H2":   {renderComponent(t, H2("t"))},
		"H3":   {renderComponent(t, H3("t"))},
		"H4":   {renderComponent(t, H4("t"))},
		"H5":   {renderComponent(t, H5("t"))},
		"H6":   {renderComponent(t, H6("t"))},
		"Link": {renderComponent(t, Link("https://example.com", "docs"))},
		"Code": {
			renderComponent(t, Code(&CodeProps{Content: "x := 1", Language: "go"})),
			renderComponent(t, Code(&CodeProps{Content: "x := 1", Wrap: true})),
			renderComponent(t, Code(&CodeProps{Content: "<b>", Wrap: true, Language: "html"})),
		},
		"Block": {renderComponent(t, Block())},
		"Card":  {renderComponent(t, Card())},
		"Panel": {
			panelClasses(""),
			panelHeaderClasses(""),
			panelClasses("blue"),
			panelHeaderClasses("blue"),
			panelClasses("nope"),
			panelHeaderClasses("nope"),
			renderComponent(t, Panel(&PanelProps{Header: "Head", Style: "blue"})),
		},
		"Breadcrumb": {renderComponent(t, Breadcrumb(
			&BreadcrumbItem{Title: "Home", URL: "/"},
			&BreadcrumbItem{Title: "Records", URL: "/records"},
			&BreadcrumbItem{Title: "Detail", Active: true},
		))},
		"Sidebar": {
			renderComponent(t, Sidebar(&SidebarProps{
				Items: []*SidebarItem{
					{Name: "Home", Icon: "M0 0", URL: "/", Active: true},
					{Name: "Messages", Icon: "M0 0", URL: "/m",
						Modifiers: []templ.Component{SidebarNew()}},
					{Name: "Tasks", Icon: "M0 0", URL: "/t",
						Modifiers: []templ.Component{SidebarRed("3")}},
					{Name: "Projects", URL: "/p",
						Modifiers: []templ.Component{SidebarGreen("New")}},
				},
				User:        &UserInfo{Name: "Jane", Email: "jane@example.com"},
				LogoutURL:   "/logout",
				Collapsible: true,
				ID:          "chrome-test",
			})),
			renderComponent(t, Sidebar(&SidebarProps{
				Items:         []*SidebarItem{{Name: "Home", URL: "/", Active: true}},
				LogoutOnClick: "doLogout()",
			})),
			renderComponent(t, SidebarNew()),
			renderComponent(t, SidebarRed("1")),
			renderComponent(t, SidebarGreen("1")),
			renderComponent(t, SidebarStyles()),
		},
	}
	return sources
}

func TestChromeEmitsNoPaletteClassesOrDarkVariants(t *testing.T) {
	for name, blobs := range chromeClassSources(t) {
		for _, html := range blobs {
			assert.Empty(t, hardcodedPaletteRe.FindString(html),
				"%s still emits a palette class: %s", name, hardcodedPaletteRe.FindString(html))
			assert.NotContains(t, html, "dark:", "%s must not carry dark: variants", name)
		}
	}
}

func TestHeadingsUseTypographyRoleTokens(t *testing.T) {
	cases := []struct {
		name     string
		html     string
		role     string
		tracking bool
	}{
		{"H1", renderComponent(t, H1("t")), "display", true},
		{"H2", renderComponent(t, H2("t")), "heading-lg", false},
		{"H3", renderComponent(t, H3("t")), "heading", false},
		{"H4", renderComponent(t, H4("t")), "subheading", false},
		{"H5", renderComponent(t, H5("t")), "body", false},
		{"H6", renderComponent(t, H6("t")), "body-sm", false},
	}
	for _, tc := range cases {
		assert.Contains(t, tc.html, "var(--text-"+tc.role+"-size)", "%s size", tc.name)
		assert.Contains(t, tc.html, "var(--text-"+tc.role+"-weight)", "%s weight", tc.name)
		assert.Contains(t, tc.html, "var(--text-"+tc.role+"-leading)", "%s leading", tc.name)
		assert.Contains(t, tc.html, "text-[var(--color-ink)]", "%s ink", tc.name)
		if tc.tracking {
			assert.Contains(t, tc.html, "var(--text-"+tc.role+"-tracking)", "%s tracking", tc.name)
		}
	}
}

func TestPanelStyleMatrix(t *testing.T) {
	// Neutral: subtle fill + border. Blue: info-soft + neutral ink (DESIGN.md
	// has no chromatic info hue). Unknown: same as neutral.
	neutral := panelClasses("")
	assert.Contains(t, neutral, "bg-[var(--color-subtle)]")
	assert.Contains(t, neutral, "border-[var(--color-border)]")
	assert.Contains(t, panelHeaderClasses(""), "text-[var(--color-ink-muted)]")

	blue := panelClasses("blue")
	assert.Contains(t, blue, "bg-[var(--color-info-soft)]")
	assert.NotContains(t, blue, "--color-action")
	assert.Contains(t, panelHeaderClasses("blue"), "text-[var(--color-ink)]")

	assert.Equal(t, neutral, panelClasses("not-a-style"))
	assert.Equal(t, panelHeaderClasses(""), panelHeaderClasses("not-a-style"))
}

func TestCodeSurfaceIsInverseTokenPair(t *testing.T) {
	html := renderComponent(t, Code(&CodeProps{Content: "a < b && c > d", Language: "go"}))
	assert.Contains(t, html, "bg-[var(--code-bg)]")
	assert.Contains(t, html, "text-[var(--code-fg)]")
	assert.Contains(t, html, "var(--font-mono)")
	// Escape behaviour is owned by templ and must survive the migration.
	assert.Contains(t, html, "a &lt; b &amp;&amp; c &gt; d")
	assert.NotContains(t, html, "a < b")

	// Both Wrap branches exist and only differ in pre whitespace/overflow.
	wrapped := renderComponent(t, Code(&CodeProps{Content: "x", Wrap: true}))
	unwrapped := renderComponent(t, Code(&CodeProps{Content: "x"}))
	assert.Contains(t, wrapped, "whitespace-pre-wrap")
	assert.NotContains(t, wrapped, "overflow-x-auto")
	assert.Contains(t, unwrapped, "overflow-x-auto")
	assert.NotContains(t, unwrapped, "whitespace-pre-wrap")
}

func TestSidebarKeepsClassHooks(t *testing.T) {
	html := renderComponent(t, Sidebar(&SidebarProps{
		Items: []*SidebarItem{
			{Name: "Home", URL: "/", Active: true},
			{Name: "Mail", URL: "/m", Modifiers: []templ.Component{SidebarNew()}},
		},
		User:        &UserInfo{Name: "Jane", Email: "jane@example.com"},
		LogoutURL:   "/logout",
		Collapsible: true,
		ID:          "hooks-test",
	}))
	for _, hook := range []string{
		"sidebar-root", "sidebar-item", "sidebar-item-text", "sidebar-item-icon",
		"sidebar-item-modifier", "sidebar-footer", "sidebar-user-info",
		"sidebar-logout", "sidebar-logout-icon", "sidebar-logout-text",
		"sidebar-icon-expand", "sidebar-icon-collapse", "sidebar-toggle-text",
	} {
		assert.Contains(t, html, hook, "class hook %q must survive", hook)
	}
	// SidebarScript/SidebarStyles depend on these class names in CSS rules.
	styles := renderComponent(t, SidebarStyles())
	for _, hook := range []string{
		"sidebar-collapsed", "sidebar-item-text", "sidebar-item-modifier",
		"sidebar-user-info", "sidebar-logout-text", "sidebar-logout-icon",
		"sidebar-icon-expand", "sidebar-icon-collapse", "sidebar-toggle-text",
	} {
		assert.Contains(t, styles, hook, "SidebarStyles must keep targeting %q", hook)
	}
}

func TestSidebarStatesAndPillsUseTokens(t *testing.T) {
	html := renderComponent(t, Sidebar(&SidebarProps{
		Items: []*SidebarItem{
			{Name: "Home", URL: "/", Active: true},
			{Name: "Mail", URL: "/m"},
		},
	}))

	// Active item: selected fill + primary ink. Inactive: ink-muted, hover highlight.
	assert.Contains(t, html, "bg-[var(--color-selected)] text-[var(--color-ink)]", "active item tokens")
	assert.Contains(t, html, "text-[var(--color-ink-muted)] hover:bg-[var(--color-highlight)]", "inactive item tokens")
	assert.Contains(t, html, "bg-[var(--color-surface)]", "sidebar root surface")

	// Pills reuse the L3 badge vocabulary, including the matching 1px border.
	for name, html := range map[string]string{
		"neutral": renderComponent(t, SidebarNew()),
		"danger":  renderComponent(t, SidebarRed("3")),
		"success": renderComponent(t, SidebarGreen("New")),
	} {
		assert.Contains(t, html, "var(--badge-"+name+"-fg)", "%s pill foreground", name)
		assert.Contains(t, html, "var(--badge-"+name+"-bg)", "%s pill fill", name)
		assert.Contains(t, html, "var(--badge-"+name+"-border)", "%s pill border", name)
	}

	// The collapsed-sidebar tooltip is an inverse pair too, not a literal hex.
	styles := renderComponent(t, SidebarStyles())
	assert.Contains(t, styles, "background: var(--color-ink)")
	assert.Contains(t, styles, "color: var(--color-canvas)")
	assert.NotContains(t, styles, "#1f2937")
	assert.NotContains(t, styles, "color: white")
}

func TestLinkAndBreadcrumbUseActionToken(t *testing.T) {
	link := renderComponent(t, Link("/docs", "docs"))
	assert.Contains(t, link, "text-[var(--color-action)]")

	bc := renderComponent(t, Breadcrumb(
		&BreadcrumbItem{Title: "Home", URL: "/"},
		&BreadcrumbItem{Title: "Now", Active: true},
	))
	assert.Contains(t, bc, "text-[var(--color-action)]", "breadcrumb links are affordances")
	assert.Contains(t, bc, "text-[var(--color-ink)]", "active crumb is a location marker")
	assert.Contains(t, bc, "text-[var(--color-ink-muted)]", "separator")
}

func TestContainerSurfacesUseSurfaceTokens(t *testing.T) {
	card := renderComponent(t, Card())
	assert.Contains(t, card, "bg-[var(--color-surface)]")
	assert.Contains(t, card, "border-[var(--color-border)]")
	assert.NotContains(t, card, "shadow", "DESIGN.md: page cards carry no shadow")

	block := renderComponent(t, Block())
	assert.Contains(t, block, "border-[var(--color-border-strong)]")
}

func TestChromeTokenContrastMeetsFloor(t *testing.T) {
	// Every text colour the chrome can emit, against every backdrop it can sit
	// on: canvas (pages), surface (cards, panels, sidebar) and highlight (hover
	// wash). No lowered thresholds.
	css := renderTokenStyles(t)
	texts := []string{"--color-ink", "--color-ink-muted", "--color-action"}
	backdrops := []string{"--color-canvas", "--color-surface", "--color-highlight"}

	// Text/fill pairs on other fills the chrome produces.
	pairs := [][2]string{
		// Sidebar active item and panel fills.
		{"--color-ink", "--color-selected"},
		{"--color-ink-muted", "--color-subtle"},
		{"--color-ink-muted", "--color-info-soft"},
		{"--color-ink", "--color-info-soft"},
		// Sidebar avatar.
		{"--color-ink-muted", "--color-surface-alt"},
		// Code block and sidebar tooltip (inverse pair).
		{"--code-fg", "--code-bg"},
		// Sidebar logout.
		{"--color-danger-fg", "--color-danger"},
		{"--color-focus", "--color-surface"},
		// Sidebar pills.
		{"--badge-neutral-fg", "--badge-neutral-bg"},
		{"--badge-danger-fg", "--badge-danger-bg"},
		{"--badge-success-fg", "--badge-success-bg"},
	}

	for _, dark := range []bool{false, true} {
		vars := themeVars(t, css, dark)
		name := "light"
		if dark {
			name = "dark"
		}
		for _, fgName := range texts {
			for _, bgName := range backdrops {
				fg := resolveVar(t, vars, fgName)
				bg := resolveVar(t, vars, bgName)
				ratio := contrastRatio(fg, bg)
				assert.GreaterOrEqual(t, ratio, 4.5,
					"%s: %s (%s) on %s (%s) = %.2f", name, fgName, fg, bgName, bg, ratio)
			}
		}
		for _, pair := range pairs {
			fg := resolveVar(t, vars, pair[0])
			bg := resolveVar(t, vars, pair[1])
			ratio := contrastRatio(fg, bg)
			assert.GreaterOrEqual(t, ratio, 4.5,
				"%s: %s (%s) on %s (%s) = %.2f", name, pair[0], fg, pair[1], bg, ratio)
		}
	}
}

func TestCodeTokensPresentInBothThemes(t *testing.T) {
	css := renderTokenStyles(t)
	for _, dark := range []bool{false, true} {
		vars := themeVars(t, css, dark)
		// The inverse pairing is defined once and flips with the theme, so the
		// slab is dark in light theme and light in dark theme.
		assert.Equal(t, resolveVar(t, vars, "--color-ink"), resolveVar(t, vars, "--code-bg"))
		assert.Equal(t, resolveVar(t, vars, "--color-canvas"), resolveVar(t, vars, "--code-fg"))
	}
}
