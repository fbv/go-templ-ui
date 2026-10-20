package ui

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StateView is the CAP-5 shared state surface: one value (StateProps), five
// treatments (EXPERIENCE.md § State Patterns), no hardcoded user-facing
// strings, and no blank surface on any path. The tests below pin the matrix
// from the ticket at the component surface.

func childMarkerComponent(marker string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, marker)
		return err
	})
}

func renderWithChildren(t *testing.T, c templ.Component, marker string) string {
	t.Helper()
	ctx := templ.WithChildren(context.Background(), childMarkerComponent(marker))
	var buf strings.Builder
	require.NoError(t, c.Render(ctx, &buf))
	return buf.String()
}

func TestStateViewLoadingRendersSkeletonRowsWithoutSpinner(t *testing.T) {
	html := renderComponent(t, StateView(&StateProps{
		State:        StateLoading,
		SkeletonRows: 5,
		SkeletonCols: 3,
	}))
	// SkeletonTable: one bordered flex row per requested row, one cell per col.
	assert.Equal(t, 5, strings.Count(html, "space-x-4 py-3 border-b"),
		"loading must render one skeleton row per SkeletonRows")
	assert.Equal(t, 15, strings.Count(html, "flex-1 h-3 bg-gray-200"),
		"loading must render SkeletonRows x SkeletonCols skeleton cells")
	assert.Contains(t, html, "<div class=\"w-full\">", "loading renders the SkeletonTable surface")
	assert.NotContains(t, html, "animate-spin",
		"loading must never render a spinner-over-blank")
}

func TestStateViewLoadingWithoutRowsFallsBackToOneSkeletonLine(t *testing.T) {
	html := renderComponent(t, StateView(&StateProps{State: StateLoading}))
	assert.Contains(t, html, "animate-pulse")
	assert.NotContains(t, html, "animate-spin")
	assert.NotContains(t, html, "space-x-4 py-3 border-b",
		"without SkeletonRows the surface is a single-line skeleton, not a row table")
}

func TestStateViewEmptyRendersSubheadingMessageAndOneAction(t *testing.T) {
	html := renderComponent(t, StateView(&StateProps{
		State:       StateEmpty,
		Message:     "No accounts yet.",
		ActionLabel: "Create account",
		ActionURL:   "/accounts/new",
	}))
	assert.Contains(t, html, "No accounts yet.")
	assert.Contains(t, html, "var(--text-subheading-size)",
		"the empty message speaks in the subheading typography role")
	assert.Contains(t, html, "var(--text-subheading-weight)")
	assert.Contains(t, html, "var(--text-subheading-leading)")
	assert.Contains(t, html, `href="/accounts/new"`)
	assert.Contains(t, html, "Create account")
	assert.Equal(t, 1, strings.Count(html, "<a "),
		"the empty surface renders exactly one primary action control")
	assert.Empty(t, hardcodedPaletteRe.FindString(html),
		"the empty surface must come from semantic tokens only")
	assert.NotContains(t, html, "dark:", "token-based markup is theme-agnostic")
}

func TestStateViewErrorRendersDangerSurfaceWithRetry(t *testing.T) {
	html := renderComponent(t, StateView(&StateProps{
		State:      StateError,
		Message:    "Accounts could not be loaded. Retry, or contact support if it persists.",
		RetryLabel: "Retry",
		RetryURL:   "/accounts",
	}))
	assert.Contains(t, html, "Accounts could not be loaded.")
	assert.Contains(t, html, "var(--color-danger-border)",
		"the error surface carries the danger accent")
	assert.Contains(t, html, "var(--color-danger)",
		"the error text is danger-coloured")
	assert.Contains(t, html, "var(--color-subtle)",
		"the error surface is a neutral fill with a danger accent, not a failure-red slab")
	assert.Contains(t, html, "Retry")
	assert.Contains(t, html, `href="/accounts"`)
	assert.Equal(t, 1, strings.Count(html, "<a "),
		"the error surface renders exactly one retry control")
	assert.Empty(t, hardcodedPaletteRe.FindString(html))
	assert.NotContains(t, html, "dark:")
}

func TestStateViewPermissionDeniedHasNoRetryControl(t *testing.T) {
	html := renderComponent(t, StateView(&StateProps{
		State:       StatePermissionDenied,
		Message:     "You do not have access to this record. Ask an administrator for access.",
		RetryLabel:  "Retry",
		RetryURL:    "/accounts",
		ActionLabel: "Create account",
		ActionURL:   "/accounts/new",
	}))
	assert.Contains(t, html, "You do not have access to this record.")
	assert.NotContains(t, html, "Retry",
		"permission-denied must render no retry control, whatever the props carry")
	assert.NotContains(t, html, "/accounts")
	assert.NotContains(t, html, "Create account")
	assert.NotContains(t, html, "/accounts/new")
	assert.NotContains(t, html, "<a ",
		"permission-denied renders no controls at all: retrying cannot help")
	assert.Contains(t, html, "var(--color-ink-muted)",
		"permission-denied speaks in muted ink on a neutral surface")
	assert.Empty(t, hardcodedPaletteRe.FindString(html))
	assert.NotContains(t, html, "dark:")
}

func TestStateViewPartialKeepsChildrenAndAddsNotice(t *testing.T) {
	html := renderWithChildren(t, StateView(&StateProps{
		State:   StatePartial,
		Message: "Owner and ARR are unavailable.",
	}), "partial-host-content")
	assert.Contains(t, html, "partial-host-content",
		"partial data renders what loaded — the children come through")
	assert.Contains(t, html, "Owner and ARR are unavailable.",
		"partial data adds an inline notice naming what is unavailable")
	assert.Contains(t, html, "var(--color-ink-muted)")
	assert.Empty(t, hardcodedPaletteRe.FindString(html))
	assert.NotContains(t, html, "dark:")
}

func TestStateViewNilRendersChildrenByteIdentical(t *testing.T) {
	// State: nil is the no-state path: children pass through untouched, which
	// is what every pre-State caller gets today.
	got := renderWithChildren(t, StateView(nil), "child-marker")
	assert.Equal(t, "child-marker", got)

	var buf strings.Builder
	require.NoError(t, StateView(nil).Render(context.Background(), &buf))
	assert.Equal(t, "", buf.String(), "nil state with no children renders nothing")
}

func TestStateViewUnknownStateFallsBackToEmptyTreatment(t *testing.T) {
	bogus := renderComponent(t, StateView(&StateProps{
		State:       DataState("bogus"),
		Message:     "No accounts yet.",
		ActionLabel: "Create account",
		ActionURL:   "/accounts/new",
	}))
	empty := renderComponent(t, StateView(&StateProps{
		State:       StateEmpty,
		Message:     "No accounts yet.",
		ActionLabel: "Create account",
		ActionURL:   "/accounts/new",
	}))
	assert.Equal(t, empty, bogus,
		"an unrecognised State value renders the empty treatment, not a blank surface")

	// An empty StateProps is an unrecognised value too: same treatment, no panic.
	blank := renderComponent(t, StateView(&StateProps{}))
	assert.Contains(t, blank, "var(--text-subheading-size)")
}

func TestStateViewNonPartialStatesIgnoreChildren(t *testing.T) {
	states := []StateProps{
		{State: StateLoading, SkeletonRows: 2},
		{State: StateEmpty, Message: "No accounts yet."},
		{State: StateError, Message: "Could not load.", RetryLabel: "Retry", RetryURL: "/"},
		{State: StatePermissionDenied, Message: "Access denied. Ask an administrator."},
		{State: ""},
		{State: DataState("bogus")},
	}
	for _, p := range states {
		p := p
		html := renderWithChildren(t, StateView(&p), "child-marker")
		assert.NotContains(t, html, "child-marker",
			"state %q replaces children with its surface (only partial keeps them)", p.State)
	}
}

func TestStateClassStringsUseSemanticTokensOnly(t *testing.T) {
	// StateView's own markup: the four non-loading states (and the empty
	// fallback) must be token-only, with no dark: variants.
	cases := []StateProps{
		{State: StateEmpty, Message: "No accounts yet."},
		{State: StateError, Message: "Could not load.", RetryLabel: "Retry", RetryURL: "/"},
		{State: StatePermissionDenied, Message: "Access denied. Ask an administrator."},
		{State: StatePartial, Message: "Owner and ARR are unavailable."},
		{State: DataState("bogus"), Message: "No accounts yet."},
	}
	for _, props := range cases {
		props := props
		html := renderComponent(t, StateView(&props))
		assert.Empty(t, hardcodedPaletteRe.FindString(html),
			"StateView(%q) has a hardcoded palette class: %s", props.State, html)
		assert.NotContains(t, html, "dark:", "StateView(%q) carries a dark: variant", props.State)
	}

	// Loading is the exception: StateView emits the pre-existing Skeleton
	// components (known palette debt owned by the refactor sweep), not its
	// own wrapper markup.
	withRows := renderComponent(t, StateView(&StateProps{State: StateLoading, SkeletonRows: 5, SkeletonCols: 3}))
	assert.Equal(t, renderComponent(t, SkeletonTable(5, 3)), withRows,
		"loading with SkeletonRows must emit SkeletonTable markup")

	noRows := renderComponent(t, StateView(&StateProps{State: StateLoading}))
	assert.Equal(t, renderComponent(t, Skeleton(&SkeletonProps{Lines: 1})), noRows,
		"loading without SkeletonRows must emit Skeleton markup")
}
