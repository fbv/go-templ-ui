package ui

import "github.com/a-h/templ"

type Column[T any] struct {
	Name    string
	Style   string // optional style
	GetData func(t T) templ.Component
}

type DataSource[T any] struct {
	Columns     []Column[T]
	Rows        []T
	RowModifier *RowModifier[T]
}

type RowModifier[T any] struct {
	Title    func(row int, d T) string
	HXGet    func(row int, d T) string
	OnClick  func(row int, d T) string
	OnDelete func(row int, d T) string
}

func (ds *DataSource[T]) GetColumnCount() int {
	return len(ds.Columns)
}

func (ds *DataSource[T]) GetColumnName(col int) string {
	return ds.Columns[col].Name
}

func (ds *DataSource[T]) GetColumnStyle(col int) string {
	return ds.Columns[col].Style
}

func (ds *DataSource[T]) GetRowCount() int {
	return len(ds.Rows)
}

func (ds *DataSource[T]) GetData(row, col int) templ.Component {
	return ds.Columns[col].GetData(ds.Rows[row])
}

func (ds *DataSource[T]) GetRowMeta(row int) *RowMeta {
	meta := &RowMeta{}
	if ds.RowModifier != nil {
		if ds.RowModifier.HXGet != nil {
			meta.HXGet = ds.RowModifier.HXGet(row, ds.Rows[row])
		}
		if ds.RowModifier.OnClick != nil {
			meta.HRef = ds.RowModifier.OnClick(row, ds.Rows[row])
		}
	}
	return meta
}
