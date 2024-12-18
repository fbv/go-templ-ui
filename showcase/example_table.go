package showcase

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view/ui"
)

var exampleTable *ui.DataSource[string] = &ui.DataSource[string]{
	Columns: []ui.Column[string]{
		{
			Name: "Col 1",
			GetData: func(row string) templ.Component {
				return ui.Text(fmt.Sprintf("%s col %d", row, 1))
			},
			Style: "w-20",
		},
		{
			Name: "Col 2",
			GetData: func(row string) templ.Component {
				return ui.Text(fmt.Sprintf("%s col %d", row, 2))
			},
			Style: "w-40",
		},
		{
			Name: "Col 3",
			GetData: func(row string) templ.Component {
				return ui.Text(fmt.Sprintf("%s col %d", row, 3))
			},
			Style: "w-60 bg-yellow-50",
		},
		{
			Name: "Col 4",
			GetData: func(row string) templ.Component {
				return ui.Text(fmt.Sprintf("%s col %d", row, 4))
			},
		},
		{
			Name: "Col 5",
			GetData: func(row string) templ.Component {
				return ui.Text(fmt.Sprintf("%s col %d", row, 5))
			},
		},
	},
	Rows: []string{
		"row 1",
		"row 2",
		"row 3",
		"row 4",
		"row 5",
		"row 6",
		"row 7",
	},
}
