package showcase

import (
	"context"
	"io/fs"
	"strconv"
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

// TestRegistryIconPathsAreValid guards the icon contract: every row's Icon is
// syntactically valid SVG path data, so the sidebar never draws a broken glyph.
func TestRegistryIconPathsAreValid(t *testing.T) {
	for _, d := range Demos() {
		assert.Empty(t, svgPathError(d.Icon), "row %q icon path %q", d.Name, d.Icon)
	}
	// The validator must be able to fail: a mangled cluster after a command is
	// exactly the corruption that draws a broken sidebar glyph.
	assert.NotEmpty(t, svgPathError("M12 3c2.5 2.7 4 5.7 4 9s-1.5 6.3-4 9c-2.5-2.7-4-5.7-4-9s1.5-6.3 4-9z158.1115.2778.138"))
	assert.NotEmpty(t, svgPathError("M12 3c2.5 2.7 4"))
	assert.NotEmpty(t, svgPathError("M12 3 4"))
}

// svgPathError reports the first grammar violation in SVG path data, or "".
// An explicit command takes a fixed parameter group and may repeat it without
// the letter; A/a parameters 4 and 5 are single-digit flags, which may run
// into the next number. The Badge row is the reference: it parses clean.
func svgPathError(d string) string {
	cmds := map[byte]int{
		'M': 2, 'm': 2, 'L': 2, 'l': 2, 'H': 1, 'h': 1, 'V': 1, 'v': 1,
		'C': 6, 'c': 6, 'S': 4, 's': 4, 'Q': 4, 'q': 4, 'T': 2, 't': 2,
		'A': 7, 'a': 7, 'Z': 0, 'z': 0,
	}
	i := 0
	var cmd byte
	for i < len(d) {
		c := d[i]
		if c == ' ' || c == ',' || c == '\t' || c == '\n' || c == '\r' {
			i++
			continue
		}
		if isPathCommand(c) {
			n, ok := cmds[c]
			if !ok {
				return "unknown command " + string(c) + " at offset " + strconv.Itoa(i)
			}
			cmd = c
			i++
			if n == 0 {
				cmd = 0
			}
			continue
		}
		a, ok := cmds[cmd]
		if !ok {
			return "parameter outside a command at offset " + strconv.Itoa(i)
		}
		for n := 0; n < a; n++ {
			if n > 0 {
				for i < len(d) && (d[i] == ' ' || d[i] == ',' || d[i] == '\t' || d[i] == '\n' || d[i] == '\r') {
					i++
				}
			}
			// A/a parameters 4 and 5 are single-digit flags that may run into
			// the following number ("a2 2 0 010 2.828").
			if (cmd == 'A' || cmd == 'a') && (n == 3 || n == 4) {
				if i < len(d) && (d[i] == '0' || d[i] == '1') {
					i++
					continue
				}
				return "arc flag must be 0 or 1 at offset " + strconv.Itoa(i)
			}
			j, err := scanPathNumber(d, i)
			if err != "" {
				return err
			}
			i = j
		}
	}
	return ""
}

func isPathCommand(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

// scanPathNumber consumes one SVG number at d[i:] and returns the next index.
func scanPathNumber(d string, i int) (int, string) {
	start := i
	if i < len(d) && (d[i] == '+' || d[i] == '-') {
		i++
	}
	digits := 0
	for i < len(d) && d[i] >= '0' && d[i] <= '9' {
		i++
		digits++
	}
	if i < len(d) && d[i] == '.' {
		i++
		for i < len(d) && d[i] >= '0' && d[i] <= '9' {
			i++
			digits++
		}
	}
	if digits == 0 {
		return i, "malformed number at offset " + strconv.Itoa(start)
	}
	if i < len(d) && (d[i] == 'e' || d[i] == 'E') {
		i++
		if i < len(d) && (d[i] == '+' || d[i] == '-') {
			i++
		}
		expDigits := 0
		for i < len(d) && d[i] >= '0' && d[i] <= '9' {
			i++
			expDigits++
		}
		if expDigits == 0 {
			return i, "malformed exponent at offset " + strconv.Itoa(start)
		}
	}
	return i, ""
}
