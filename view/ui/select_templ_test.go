package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	Name  string
	Value string
}

func (a A) GetName() string {
	return a.Name
}

func (a A) GetValue() string {
	return a.Value
}

func TestListItemsProvider(t *testing.T) {
	l := []A{
		{Name: "a", Value: "1"},
		{Name: "b", Value: "2"},
	}
	w := ListItemProvider[A](l)
	assert.Equal(t, 2, w.GetItemsCount())
	assert.Equal(t, "1", w.GetItem(0).GetValue())
	assert.Equal(t, "b", w.GetItem(1).GetName())
}
