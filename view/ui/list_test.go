package ui

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// SimpleListWithState is the wrapper carrier of the shared CAP-5 state API:
// the list is a pure interface-taker with no props struct, so the wrapper is
// the only additive reach. A nil state renders SimpleList unchanged; a set
// state routes through StateView, which keeps the list's items only for
// StatePartial and paints its own surface for the other four treatments.

type staticListItems []*ListItemProps

func (s staticListItems) GetItemCount() int            { return len(s) }
func (s staticListItems) GetItem(i int) *ListItemProps { return s[i] }

func listStateFixture() ListDataSource {
	return staticListItems{
		{Label: "Acme Logistics"},
		{Label: "Northwind Traders"},
	}
}

func TestSimpleListWithNilStateIsByteIdenticalToSimpleList(t *testing.T) {
	ds := listStateFixture()
	assert.Equal(t, renderComponent(t, SimpleList(ds)),
		renderComponent(t, SimpleListWithState(ds, nil)),
		"a nil state must render exactly what the unwrapped component renders")
}

func TestSimpleListWithStateRendersStateSurfaceWithoutList(t *testing.T) {
	cases := []struct {
		name    string
		state   StateProps
		want    []string
		notWant []string
	}{
		{
			name:    "loading",
			state:   StateProps{State: StateLoading, SkeletonRows: 3, SkeletonCols: 1},
			want:    []string{"space-x-4 py-3 border-b", "flex-1 h-3 bg-gray-200"},
			notWant: []string{"Acme Logistics", "animate-spin"},
		},
		{
			name: "empty",
			state: StateProps{
				State:       StateEmpty,
				Message:     "No accounts yet.",
				ActionLabel: "Create account",
				ActionURL:   "/accounts/new",
			},
			want:    []string{"No accounts yet.", "var(--text-subheading-size)", `href="/accounts/new"`},
			notWant: []string{"Acme Logistics"},
		},
		{
			name: "error",
			state: StateProps{
				State:      StateError,
				Message:    "Accounts could not be loaded. Retry, or contact support if it persists.",
				RetryLabel: "Retry",
				RetryURL:   "/accounts",
			},
			want:    []string{"Accounts could not be loaded.", "var(--color-danger-border)", `href="/accounts"`},
			notWant: []string{"Acme Logistics"},
		},
		{
			name: "permission-denied",
			state: StateProps{
				State:      StatePermissionDenied,
				Message:    "You do not have access to this list. Ask an administrator for access.",
				RetryLabel: "Retry",
				RetryURL:   "/accounts",
			},
			want:    []string{"You do not have access to this list.", "var(--color-ink-muted)"},
			notWant: []string{"Acme Logistics", "Retry", `href="/accounts"`},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := renderComponent(t, SimpleListWithState(listStateFixture(), &tc.state))
			for _, w := range tc.want {
				assert.Contains(t, got, w)
			}
			for _, nw := range tc.notWant {
				assert.NotContains(t, got, nw)
			}
		})
	}
}

func TestSimpleListWithLoadingStateShowsSkeletonMatchingRequestedShape(t *testing.T) {
	got := renderComponent(t, SimpleListWithState(listStateFixture(), &StateProps{
		State: StateLoading, SkeletonRows: 3, SkeletonCols: 1,
	}))
	assert.Equal(t, 3, strings.Count(got, "space-x-4 py-3 border-b"),
		"loading must show a Skeleton matching the requested row count")
	assert.Equal(t, 3, strings.Count(got, "flex-1 h-3 bg-gray-200"),
		"loading must show SkeletonRows x SkeletonCols skeleton cells")
	assert.NotContains(t, got, "animate-spin", "never a spinner over a blank host")
}

func TestSimpleListWithPartialStateKeepsItemsAndAddsNotice(t *testing.T) {
	got := renderComponent(t, SimpleListWithState(listStateFixture(), &StateProps{
		State:   StatePartial,
		Message: "Owner is unavailable.",
	}))
	assert.Contains(t, got, "Acme Logistics", "partial data keeps the items that loaded")
	assert.Contains(t, got, "Northwind Traders")
	assert.Contains(t, got, "Owner is unavailable.",
		"partial data adds an inline notice naming what is unavailable")
	assert.Contains(t, got, "var(--color-ink-muted)")
}
