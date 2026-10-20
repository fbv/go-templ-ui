package showcase

import (
	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view/ui"
)

// stateDemoAccountsTable is the state page's table row data: two columns, two
// rows — what the partial treatment keeps when some fields are unavailable.
func stateDemoAccountsTable() *ui.DataSource[string] {
	return &ui.DataSource[string]{
		Columns: []ui.Column[string]{
			{
				Name: "Account",
				GetData: func(s string) templ.Component {
					return ui.Text(s)
				},
			},
			{
				Name: "Stage",
				GetData: func(s string) templ.Component {
					return ui.Text("Renewal")
				},
			},
		},
		Rows: []string{"Acme Logistics", "Northwind Traders"},
	}
}

// stateDemoAccountItems is the state page's list row data, satisfying
// ui.ListDataSource.
type stateDemoAccountItems []*ui.ListItemProps

func (s stateDemoAccountItems) GetItemCount() int { return len(s) }

func (s stateDemoAccountItems) GetItem(i int) *ui.ListItemProps { return s[i] }

func stateDemoOpenAccounts() stateDemoAccountItems {
	return stateDemoAccountItems{
		{Label: "Acme Logistics"},
		{Label: "Northwind Traders"},
	}
}
