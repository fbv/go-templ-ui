package ui

import "github.com/a-h/templ"

type TreeSelectionMode string

const (
	TreeSelectionNone   TreeSelectionMode = ""
	TreeSelectionSingle TreeSelectionMode = "single"
	TreeSelectionMulti  TreeSelectionMode = "multi"
)

type TreeStyleVariant string

const (
	TreeStyleDefault  TreeStyleVariant = ""
	TreeStyleBordered TreeStyleVariant = "bordered"
	TreeStyleCompact  TreeStyleVariant = "compact"
)

type TreeNode struct {
	ID       string
	Label    string
	Icon     string
	IconRaw  string
	URL      string
	Children []*TreeNode
	Disabled bool
	Expanded bool
	Selected bool
	Badge    string
	Actions  templ.Component
}

type TreeDataSource interface {
	GetChildNodes(parentID string) []*TreeNode
}

type TreeProps struct {
	ID              string
	Nodes           []*TreeNode
	DataSource      TreeDataSource
	Selection       TreeSelectionMode
	Style           TreeStyleVariant
	ShowIcons       bool
	LazyLoad        bool
	LazyLoadURL     string
	PersistState    bool
	Scrollable      bool
	ScrollMaxHeight string
	OnSelect        string
}

func treeContainerClasses(p *TreeProps) string {
	classes := "text-sm text-gray-800"
	switch p.Style {
	case TreeStyleBordered:
		classes += " border rounded-lg divide-y divide-gray-200"
	case TreeStyleCompact:
		classes += " text-xs"
	}
	if p.Scrollable {
		maxH := p.ScrollMaxHeight
		if maxH == "" {
			maxH = "h-96"
		}
		classes += " " + maxH + " overflow-y-auto"
	}
	return classes
}
