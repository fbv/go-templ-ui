package ui

import (
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

// DynamicTableWithState is the wrapper carrier of the shared CAP-5 state API:
// the table is a pure interface-taker with no props struct, so the wrapper is
// the only additive reach. A nil state renders DynamicTable unchanged; a set
// state routes through StateView, which keeps the table's rows only for
// StatePartial and paints its own surface for the other four treatments.

func tableStateFixture() *DataSource[string] {
	return &DataSource[string]{
		Columns: []Column[string]{
			{
				Name: "Account",
				GetData: func(s string) templ.Component {
					return Text(s)
				},
			},
			{
				Name: "Stage",
				GetData: func(s string) templ.Component {
					return Text("Renewal")
				},
			},
		},
		Rows: []string{"Acme Logistics", "Northwind Traders"},
	}
}

func TestDynamicTableWithNilStateIsByteIdenticalToDynamicTable(t *testing.T) {
	ds := tableStateFixture()
	assert.Equal(t, renderComponent(t, DynamicTable(ds)),
		renderComponent(t, DynamicTableWithState(ds, nil)),
		"a nil state must render exactly what the unwrapped component renders")
}

func TestDynamicTableWithStateRendersStateSurfaceWithoutTable(t *testing.T) {
	cases := []struct {
		name    string
		state   StateProps
		want    []string
		notWant []string
	}{
		{
			name:    "loading",
			state:   StateProps{State: StateLoading, SkeletonRows: 4, SkeletonCols: 2},
			want:    []string{"space-x-4 py-3 border-b", "flex-1 h-3 bg-gray-200"},
			notWant: []string{"<table", "animate-spin"},
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
			notWant: []string{"<table"},
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
			notWant: []string{"<table"},
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
			notWant: []string{"<table", "Retry", `href="/accounts"`},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := renderComponent(t, DynamicTableWithState(tableStateFixture(), &tc.state))
			for _, w := range tc.want {
				assert.Contains(t, got, w)
			}
			for _, nw := range tc.notWant {
				assert.NotContains(t, got, nw)
			}
		})
	}
}

func TestDynamicTableWithLoadingStateShowsSkeletonMatchingRequestedShape(t *testing.T) {
	got := renderComponent(t, DynamicTableWithState(tableStateFixture(), &StateProps{
		State: StateLoading, SkeletonRows: 4, SkeletonCols: 2,
	}))
	assert.Equal(t, 4, strings.Count(got, "space-x-4 py-3 border-b"),
		"loading must show a Skeleton matching the requested row count")
	assert.Equal(t, 8, strings.Count(got, "flex-1 h-3 bg-gray-200"),
		"loading must show SkeletonRows x SkeletonCols skeleton cells")
	assert.NotContains(t, got, "animate-spin", "never a spinner over a blank host")
}

func TestDynamicTableWithPartialStateKeepsRowsAndAddsNotice(t *testing.T) {
	got := renderComponent(t, DynamicTableWithState(tableStateFixture(), &StateProps{
		State:   StatePartial,
		Message: "Stage is unavailable.",
	}))
	assert.Contains(t, got, "<table", "partial data keeps the rows that loaded")
	assert.Contains(t, got, "Acme Logistics")
	assert.Contains(t, got, "Northwind Traders")
	assert.Contains(t, got, "Stage is unavailable.",
		"partial data adds an inline notice naming what is unavailable")
	assert.Contains(t, got, "var(--color-ink-muted)")
}
