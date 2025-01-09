package ui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/assert"
)

type TestRowModifier struct {
}

func (m *TestRowModifier) HXGet(row int, s string) string {
	return "/app/" + s
}

func TestDataSource(t *testing.T) {
	ds := &DataSource[string]{
		Columns: []Column[string]{
			{
				Name: "col1",
				GetData: func(t string) templ.Component {
					return Text(t)
				},
			},
		},
		Rows:        []string{"row1", "row2", "row3"},
		RowModifier: &TestRowModifier{},
	}
	assert.Equal(t, 1, ds.GetColumnCount())

	ctx := context.Background()
	var b bytes.Buffer
	err := DynamicTable(ds).Render(ctx, &b)
	assert.NoError(t, err)
	s := b.String()
	assert.True(t, strings.Contains(s, "hx-get=\"/app/row3\""))
}
