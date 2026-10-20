package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Panel is the first carrier of the shared CAP-5 state API: an additive
// *StateProps field that renders the StateView surface in place of the
// panel's children, keeping the header. Without it, the markup is
// byte-identical to the pre-State Panel.

func TestPanelWithoutStateKeepsPreStateMarkup(t *testing.T) {
	// The exact markup Panel rendered before the State field existed; the
	// nil-state path must not drift a single byte from it.
	want := `<div class="bg-[var(--color-subtle)] p-4 rounded-lg border border-[var(--color-border)]"><p class="text-sm text-[var(--color-ink-muted)] mb-4">Head</p>child-marker</div>`
	got := renderWithChildren(t, Panel(&PanelProps{Header: "Head"}), "child-marker")
	require.Equal(t, want, got)
}

func TestPanelWithStateRendersStateSurfaceAndDropsChildren(t *testing.T) {
	p := &PanelProps{
		Header: "Accounts",
		State: &StateProps{
			State:       StateEmpty,
			Message:     "No accounts yet.",
			ActionLabel: "Create account",
			ActionURL:   "/accounts/new",
		},
	}
	got := renderWithChildren(t, Panel(p), "child-marker")
	assert.Contains(t, got, ">Accounts<", "the header stays when a state is set")
	assert.Contains(t, got, "No accounts yet.")
	assert.Contains(t, got, "var(--text-subheading-size)")
	assert.Contains(t, got, `href="/accounts/new"`)
	assert.NotContains(t, got, "child-marker",
		"the state surface replaces the panel's children")
}

func TestPanelLoadingStateShowsSkeletonNotSpinner(t *testing.T) {
	p := &PanelProps{
		Header: "Accounts",
		State:  &StateProps{State: StateLoading, SkeletonRows: 5, SkeletonCols: 3},
	}
	got := renderWithChildren(t, Panel(p), "child-marker")
	assert.Contains(t, got, ">Accounts<")
	assert.Equal(t, 5, strings.Count(got, "space-x-4 py-3 border-b"),
		"loading must show a Skeleton matching the requested row count")
	assert.Equal(t, 15, strings.Count(got, "flex-1 h-3 bg-gray-200"))
	assert.NotContains(t, got, "animate-spin", "never a spinner over a blank host")
	assert.NotContains(t, got, "child-marker")
}

func TestPanelPartialStateKeepsChildrenAndAddsNotice(t *testing.T) {
	p := &PanelProps{
		Header: "Accounts",
		State:  &StateProps{State: StatePartial, Message: "Owner and ARR are unavailable."},
	}
	got := renderWithChildren(t, Panel(p), "child-marker")
	assert.Contains(t, got, ">Accounts<")
	assert.Contains(t, got, "child-marker", "partial data keeps what loaded")
	assert.Contains(t, got, "Owner and ARR are unavailable.")
}

func TestPanelPermissionDeniedCarrierHasNoRetryControl(t *testing.T) {
	p := &PanelProps{
		State: &StateProps{
			State:       StatePermissionDenied,
			Message:     "You do not have access to this record. Ask an administrator for access.",
			RetryLabel:  "Retry",
			RetryURL:    "/accounts",
			ActionLabel: "Create account",
			ActionURL:   "/accounts/new",
		},
	}
	got := renderWithChildren(t, Panel(p), "child-marker")
	assert.NotContains(t, got, "Retry", "permission-denied never renders a retry")
	assert.NotContains(t, got, "/accounts")
	assert.NotContains(t, got, "Create account")
	assert.Contains(t, got, "You do not have access to this record.")
}
