package showcase

import "github.com/a-h/templ"

// Demo is one entry in the shared demo registry. Later epics add rows here;
// server/server.go registers routes from it and showcase.templ builds sidebar
// links from it.
//
// Every field is read somewhere: Name identifies the page (and its DOM id),
// Title labels the sidebar link and the demo page heading, Icon draws the
// sidebar icon, Route is the registered URL, DocsPath is both the demo page's
// docs link and the served /showcase/docs target, and Render receives the row
// itself so a demo cannot restate its own name/title/docs path.
//
// Demos() builds the slice per call (no package-level var) so metadata
// consumers such as Sidebar do not create an initialization cycle through the
// demo components.
type Demo struct {
	Name     string
	Title    string
	Icon     string
	Route    string
	DocsPath string
	Render   func(d Demo) templ.Component
}

// Demos returns the registered demo pages.
func Demos() []Demo {
	return []Demo{
		{
			Name:     "badge",
			Title:    "Badge",
			Icon:     "M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z",
			Route:    "/showcase/badge",
			DocsPath: "docs/badge.md",
			Render:   BadgeDemoPage,
		},
		{
			Name:     "i18n",
			Title:    "i18n",
			Icon:     "M12 3a9 9 0 1 0 0 18 9 9 0 1 0 0-18zM3 12h18M12 3c2.5 2.7 4 5.7 4 9s-1.5 6.3-4 9c-2.5-2.7-4-5.7-4-9s1.5-6.3 4-9z",
			Route:    "/showcase/i18n",
			DocsPath: "docs/i18n.md",
			Render:   I18nDemoPage,
		},
		{
			Name:     "state",
			Title:    "State",
			Icon:     "M4 5h16v3H4V5zm0 5h16v3H4v-3zm0 5h16v3H4v-3z",
			Route:    "/showcase/state",
			DocsPath: "docs/state.md",
			Render:   StateDemoPage,
		},
	}
}
