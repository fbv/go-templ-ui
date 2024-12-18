package ui

import "github.com/a-h/templ"

type Column[T any] struct {
	Name    string
	GetData func(t T) templ.Component
}

type DataSource[T any] struct {
	Columns []Column[T]
	Rows    []T
}

func (ds *DataSource[T]) GetColumnCount() int {
	return len(ds.Columns)
}

func (ds *DataSource[T]) GetColumnName(col int) string {
	return ds.Columns[col].Name
}

func (ds *DataSource[T]) GetRowCount() int {
	return len(ds.Rows)
}

func (ds *DataSource[T]) GetData(row, col int) templ.Component {
	return ds.Columns[col].GetData(ds.Rows[row])
}
