package ui

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTreeRendersNodes(t *testing.T) {
	nodes := []*TreeNode{
		{ID: "1", Label: "Root 1"},
		{ID: "2", Label: "Root 2"},
	}
	p := &TreeProps{
		ID:    "test-tree",
		Nodes: nodes,
	}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "Root 1")
	assert.Contains(t, s, "Root 2")
	assert.Contains(t, s, "id=\"test-tree\"")
	assert.Contains(t, s, "role=\"tree\"")
}

func TestTreeRendersChildren(t *testing.T) {
	nodes := []*TreeNode{
		{
			ID:       "1",
			Label:    "Parent",
			Expanded: true,
			Children: []*TreeNode{
				{ID: "1-1", Label: "Child 1"},
				{ID: "1-2", Label: "Child 2"},
			},
		},
	}
	p := &TreeProps{ID: "test-tree", Nodes: nodes}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "Parent")
	assert.Contains(t, s, "Child 1")
	assert.Contains(t, s, "Child 2")
	assert.Contains(t, s, "tree-expanded")
}

func TestTreeSelection(t *testing.T) {
	nodes := []*TreeNode{
		{ID: "1", Label: "Node 1"},
		{ID: "2", Label: "Node 2", Selected: true},
	}
	p := &TreeProps{ID: "test-tree", Nodes: nodes, Selection: TreeSelectionSingle}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "tree-selected")
}

func TestTreeLazyLoadAttributes(t *testing.T) {
	nodes := []*TreeNode{
		{
			ID:    "1",
			Label: "Lazy Parent",
			Children: []*TreeNode{
				{ID: "1-1", Label: "Child"},
			},
		},
	}
	p := &TreeProps{
		ID:          "test-tree",
		Nodes:       nodes,
		LazyLoad:    true,
		LazyLoadURL: "/api/tree",
	}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "data-tree-lazy")
	assert.Contains(t, s, "/api/tree/1")
}

func TestTreeStyleVariants(t *testing.T) {
	p := &TreeProps{ID: "t", Style: TreeStyleBordered, Nodes: []*TreeNode{}}
	classes := treeContainerClasses(p)
	assert.Contains(t, classes, "border")
	assert.Contains(t, classes, "rounded-lg")

	p2 := &TreeProps{ID: "t", Style: TreeStyleCompact, Nodes: []*TreeNode{}}
	classes2 := treeContainerClasses(p2)
	assert.Contains(t, classes2, "text-xs")
}

func TestTreeWithIcons(t *testing.T) {
	nodes := []*TreeNode{
		{ID: "1", Label: "Node 1", Icon: "M9 5l7 7-7 7"},
	}
	p := &TreeProps{ID: "test-tree", Nodes: nodes, ShowIcons: true}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "M9 5l7 7-7 7")
}

func TestTreeWithBadge(t *testing.T) {
	nodes := []*TreeNode{
		{ID: "1", Label: "Node 1", Badge: "New"},
	}
	p := &TreeProps{ID: "test-tree", Nodes: nodes}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "New")
}

func TestTreeDisabledNode(t *testing.T) {
	nodes := []*TreeNode{
		{ID: "1", Label: "Node 1", Disabled: true},
	}
	p := &TreeProps{ID: "test-tree", Nodes: nodes}
	var b bytes.Buffer
	err := Tree(p).Render(context.Background(), &b)
	assert.NoError(t, err)
	s := b.String()
	assert.Contains(t, s, "opacity-50")
	assert.Contains(t, s, "cursor-not-allowed")
}

func TestTreeScrollable(t *testing.T) {
	p := &TreeProps{
		ID:              "test-tree",
		Nodes:           []*TreeNode{},
		Scrollable:      true,
		ScrollMaxHeight: "h-64",
	}
	classes := treeContainerClasses(p)
	assert.Contains(t, classes, "h-64")
	assert.Contains(t, classes, "overflow-y-auto")
}
