package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Select is the props-struct carrier of the shared CAP-5 state API: an
// additive *StateProps field, the same shape Panel carries. A nil state keeps
// the pre-State markup byte-for-byte; a set state renders StateView's surface
// in place of the <select> and its label, except StatePartial which keeps the
// label and select and adds the inline notice.

func selectStateFixture() *SelectProps {
	return &SelectProps{
		Name:          "country",
		Label:         "Country",
		Placeholder:   "Choose a country",
		Items:         []SelectItem{{Name: "Germany", Value: "de"}, {Name: "France", Value: "fr"}},
		SelectedValue: "de",
	}
}

func TestSelectNilStateKeepsPreStateMarkupByteIdentical(t *testing.T) {
	// The exact markup Select rendered before the State field existed; the
	// nil-state path must not drift a single byte from it.
	want := `<div><label for="country" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Country</label> <select id="country" name="country" class="border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full py-2 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"><option hidden value="">Choose a country</option><option value="de" selected>Germany</option><option value="fr">France</option></select></div>`

	p := selectStateFixture()
	got := renderComponent(t, Select(p))
	assert.Equal(t, want, got)
	assert.Contains(t, got, `<label for="country"`,
		"the nil-state label's for= names the control it labels")
	assert.Contains(t, got, `<select id="country"`,
		"the select's id matches the label's for=")

	p.State = nil
	assert.Equal(t, want, renderComponent(t, Select(p)),
		"an explicit nil State renders exactly what the component renders without one")
}

func TestSelectNonPartialStatesRenderStateSurfaceWithoutSelect(t *testing.T) {
	cases := []struct {
		name    string
		state   StateProps
		want    []string
		notWant []string
	}{
		{
			name:    "loading",
			state:   StateProps{State: StateLoading, SkeletonRows: 3, SkeletonCols: 2},
			want:    []string{"space-x-4 py-3 border-b", "flex-1 h-3 bg-gray-200"},
			notWant: []string{"<select", "animate-spin"},
		},
		{
			name: "empty",
			state: StateProps{
				State:       StateEmpty,
				Message:     "No countries available.",
				ActionLabel: "Reload",
				ActionURL:   "/countries",
			},
			want:    []string{"No countries available.", "var(--text-subheading-size)", `href="/countries"`},
			notWant: []string{"<select"},
		},
		{
			name: "error",
			state: StateProps{
				State:      StateError,
				Message:    "Countries could not be loaded. Retry, or contact support if it persists.",
				RetryLabel: "Retry",
				RetryURL:   "/countries",
			},
			want:    []string{"Countries could not be loaded.", "var(--color-danger-border)", `href="/countries"`},
			notWant: []string{"<select"},
		},
		{
			name: "permission-denied",
			state: StateProps{
				State:      StatePermissionDenied,
				Message:    "You do not have access to this list. Ask an administrator for access.",
				RetryLabel: "Retry",
				RetryURL:   "/countries",
			},
			want:    []string{"You do not have access to this list.", "var(--color-ink-muted)"},
			notWant: []string{"<select", "Retry", `href="/countries"`},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := selectStateFixture()
			p.State = &tc.state
			got := renderComponent(t, Select(p))
			for _, w := range tc.want {
				assert.Contains(t, got, w)
			}
			for _, nw := range tc.notWant {
				assert.NotContains(t, got, nw)
			}
			assert.NotContains(t, got, `for="country"`,
				"the label is suppressed with the select, so no for= may point at the control's id")
		})
	}
}

func TestSelectLoadingStateShowsSkeletonMatchingRequestedShape(t *testing.T) {
	p := selectStateFixture()
	p.State = &StateProps{State: StateLoading, SkeletonRows: 3, SkeletonCols: 2}
	got := renderComponent(t, Select(p))
	assert.Equal(t, 3, strings.Count(got, "space-x-4 py-3 border-b"),
		"loading must show a Skeleton matching the requested row count")
	assert.Equal(t, 6, strings.Count(got, "flex-1 h-3 bg-gray-200"),
		"loading must show SkeletonRows x SkeletonCols skeleton cells")
	assert.NotContains(t, got, "animate-spin", "never a spinner over a blank host")
}

func TestSelectPartialStateKeepsSelectAndAddsNotice(t *testing.T) {
	p := selectStateFixture()
	p.State = &StateProps{State: StatePartial, Message: "France is unavailable."}
	got := renderComponent(t, Select(p))
	assert.Contains(t, got, "<select", "partial data keeps the real control")
	assert.Contains(t, got, `<label for="country"`,
		"partial data keeps the label, and its for= names the control it labels")
	assert.Contains(t, got, `<select id="country"`,
		"the select's id matches the label's for=")
	assert.Contains(t, got, ">Germany<")
	assert.Contains(t, got, ">France<")
	assert.Contains(t, got, "France is unavailable.",
		"partial data adds an inline notice naming what is unavailable")
	assert.Contains(t, got, "var(--color-ink-muted)")
}

func TestSelectPermissionDeniedHasNoRetryControl(t *testing.T) {
	p := selectStateFixture()
	p.State = &StateProps{
		State:       StatePermissionDenied,
		Message:     "You do not have access to this list. Ask an administrator for access.",
		RetryLabel:  "Retry",
		RetryURL:    "/countries",
		ActionLabel: "Reload",
		ActionURL:   "/countries",
	}
	got := renderComponent(t, Select(p))
	assert.NotContains(t, got, "Retry", "permission-denied never renders a retry")
	assert.NotContains(t, got, "<a ", "permission-denied renders no controls at all")
	assert.Contains(t, got, "You do not have access to this list.")
}
