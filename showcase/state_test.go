package showcase

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The state page is CAP-5's demonstration: table, list, select and panel each
// appear in all five states, with consumer-supplied text. This test pins that
// matrix so a dropped variant fails the suite rather than the review: both the
// variant titles and the rendered surfaces each preview must produce.
func TestStateDemoPageShowsFiveStatesPerComponent(t *testing.T) {
	var page Demo
	for _, d := range Demos() {
		if d.Name == "state" {
			page = d
			break
		}
	}
	require.NotEmpty(t, page.Render, "the registry has no state row")

	var buf strings.Builder
	require.NoError(t, page.Render(page).Render(context.Background(), &buf))
	html := buf.String()

	components := []string{"Panel", "Table", "List", "Select"}
	states := []string{"Loading", "Empty", "Error", "Permission denied", "Partial data"}
	for _, component := range components {
		for _, state := range states {
			assert.Contains(t, html, component+" — "+state+" — ",
				"the state page must show %s in the %s state", component, state)
		}
	}

	// Consumer-supplied text travels with each carrier; the per-surface no-
	// retry rule is pinned at the carrier level (view/ui/*_test.go).
	assert.Contains(t, html, "No accounts yet.")
	assert.Contains(t, html, "No countries available.")
	assert.Contains(t, html, "Owner and ARR are unavailable.")
	assert.Contains(t, html, "France is unavailable.")

	// Probe the rendered surfaces themselves. The ui.Code snippets escape
	// their content, so the raw markup below can only come from the live
	// previews — a preview that renders nothing fails these counts.
	assert.Equal(t, 1, strings.Count(html, "<select"),
		"only the Select partial preview keeps its control")
	assert.Equal(t, 14, strings.Count(html, "space-x-4 py-3 border-b"),
		"the four loading previews paint their skeleton rows: Panel 5 + Table 4 + List 3 + Select 2")
	assert.Equal(t, 35, strings.Count(html, "flex-1 h-3 bg-gray-200"),
		"the four loading previews paint SkeletonRows x SkeletonCols cells: 20+8+3+4")
	assert.Equal(t, 3, strings.Count(html, "<table"),
		"the props table plus the Panel and Table partial previews render a table")
	assert.Equal(t, 4, strings.Count(html, "px-3 py-2 mt-4"),
		"each carrier's partial preview renders the inline notice")
}
