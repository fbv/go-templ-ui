package ui

import "fmt"

type SimpleDataSource[T any] struct {
	Rows        []T
	RowModifier *RowModifier[T]
}

func (ds *SimpleDataSource[T]) GetItemCount() int {
	return len(ds.Rows)
}

func (ds *SimpleDataSource[T]) GetItem(i int) *ListItemProps {
	item := &ListItemProps{}
	d := ds.Rows[i]
	if ds.RowModifier != nil {
		if ds.RowModifier.Title != nil {
			item.Label = ds.RowModifier.Title(i, d)
		}
		if ds.RowModifier.OnClick != nil {
			item.URL = ds.RowModifier.OnClick(i, d)
		}
		if ds.RowModifier.OnDelete != nil {
			item.OnDelete = ds.RowModifier.OnDelete(i, d)
		}
	}
	if item.Label == "" {
		item.Label = fmt.Sprintf("%v", d)
	}
	return item
}
