package showcase

import (
	"context"
	"io/fs"
	"strings"
	"testing"

	"github.com/fbv/go-templ-ui/docs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegistryRowsAreComplete guards the registry contract: every field is
// read somewhere, so an empty one is a broken link, a blank page or a dead
// docs URL rather than a harmless default.
func TestRegistryRowsAreComplete(t *testing.T) {
	for _, d := range Demos() {
		assert.NotEmpty(t, d.Name, "row %+v has no Name", d)
		assert.NotEmpty(t, d.Title, "row %q has no Title", d.Name)
		assert.NotEmpty(t, d.Icon, "row %q has no sidebar Icon", d.Name)
		assert.True(t, strings.HasPrefix(d.Route, "/showcase/"), "row %q route %q must live under /showcase", d.Name, d.Route)
		assert.NotNil(t, d.Render, "row %q has no Render", d.Name)
	}
}

func TestRegistryDocsPathIsServed(t *testing.T) {
	// Demo pages link to docs/<name>.md relative to /showcase/<name>; the server
	// mounts the embedded docs tree at /showcase/docs, so the path must resolve
	// to a file the embed actually carries.
	for _, d := range Demos() {
		require.True(t, strings.HasPrefix(d.DocsPath, "docs/"), "row %q docs path %q must start with docs/", d.Name, d.DocsPath)
		file := strings.TrimPrefix(d.DocsPath, "docs/")
		_, err := fs.Stat(docs.Content, file)
		assert.NoError(t, err, "row %q links %s but the file is not embedded", d.Name, d.DocsPath)
	}
}

func TestRegistryRenderReceivesRowMetadata(t *testing.T) {
	// The demo page must take its title, docs path and DOM id from the registry
	// row, not restate them: change the row and the rendered page follows.
	for _, d := range Demos() {
		var buf strings.Builder
		require.NoError(t, d.Render(d).Render(context.Background(), &buf))
		html := buf.String()
		assert.Contains(t, html, d.Title, "row %q title missing from its page", d.Name)
		assert.Contains(t, html, d.DocsPath, "row %q docs path missing from its page", d.Name)
		assert.Contains(t, html, `id="demo-`+d.Name+`"`, "row %q Name must identify the page", d.Name)
	}
}

func TestRegistrySidebarLinksComeFromRegistry(t *testing.T) {
	var buf strings.Builder
	require.NoError(t, Sidebar().Render(context.Background(), &buf))
	html := buf.String()
	for _, d := range Demos() {
		assert.Contains(t, html, `href="`+d.Route+`"`, "sidebar must link row %q", d.Name)
		assert.Contains(t, html, d.Title, "sidebar must label row %q", d.Name)
	}
}
