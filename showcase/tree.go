package showcase

import (
	"github.com/fbv/go-templ-ui/view/icon"
	"github.com/fbv/go-templ-ui/view/ui"
)

func GetDefaultTreeNodes() []*ui.TreeNode {
	return []*ui.TreeNode{
		{
			ID:       "1",
			Label:    "Documents",
			Icon:     icon.Folder,
			Expanded: true,
			Children: []*ui.TreeNode{
				{
					ID:       "1-1",
					Label:    "Work",
					Icon:     icon.Folder,
					Expanded: true,
					Children: []*ui.TreeNode{
						{ID: "1-1-1", Label: "report.pdf", Icon: icon.FileText},
						{ID: "1-1-2", Label: "budget.xlsx", Icon: icon.File},
					},
				},
				{ID: "1-2", Label: "Personal", Icon: icon.Folder},
			},
		},
		{
			ID:    "2",
			Label: "Pictures",
			Icon:  icon.Folder,
			Children: []*ui.TreeNode{
				{ID: "2-1", Label: "vacation.jpg", Icon: icon.Image},
				{ID: "2-2", Label: "family.png", Icon: icon.Image},
			},
		},
		{ID: "3", Label: "Music", Icon: icon.Folder},
		{ID: "4", Label: "Downloads", Icon: icon.Folder, Badge: "3 new"},
	}
}

func GetTreeWithBadges() []*ui.TreeNode {
	return []*ui.TreeNode{
		{
			ID:       "inbox",
			Label:    "Inbox",
			Icon:     icon.Inbox,
			Expanded: true,
			Children: []*ui.TreeNode{
				{ID: "inbox-1", Label: "New message", Icon: icon.Email, Badge: "New"},
				{ID: "inbox-2", Label: "Read later", Icon: icon.Clock, Badge: "5"},
				{ID: "inbox-3", Label: "Archived", Icon: icon.Folder},
			},
		},
		{
			ID:    "tasks",
			Label: "Tasks",
			Icon:  icon.Tasks,
			Children: []*ui.TreeNode{
				{ID: "tasks-1", Label: "In Progress", Badge: "2"},
				{ID: "tasks-2", Label: "Completed"},
				{ID: "tasks-3", Label: "Overdue", Badge: "1"},
			},
		},
	}
}

func GetLazyTreeNodes() []*ui.TreeNode {
	return []*ui.TreeNode{
		{
			ID:    "lazy-1",
			Label: "Parent Node 1",
			Icon:  icon.Folder,
			Children: []*ui.TreeNode{
				{ID: "lazy-1-1", Label: "Child 1-1"},
				{ID: "lazy-1-2", Label: "Child 1-2"},
			},
		},
		{
			ID:    "lazy-2",
			Label: "Parent Node 2",
			Icon:  icon.Folder,
			Children: []*ui.TreeNode{
				{ID: "lazy-2-1", Label: "Child 2-1"},
				{ID: "lazy-2-2", Label: "Child 2-2"},
			},
		},
		{
			ID:    "lazy-3",
			Label: "Parent Node 3 (Lazy)",
			Icon:  icon.Folder,
			Children: []*ui.TreeNode{
				{ID: "lazy-3-1", Label: "Child 3-1"},
			},
		},
	}
}

func GetDeepTreeNodes() []*ui.TreeNode {
	return []*ui.TreeNode{
		{
			ID:       "deep-1",
			Label:    "Level 1",
			Icon:     icon.Folder,
			Expanded: true,
			Children: []*ui.TreeNode{
				{
					ID:    "deep-1-1",
					Label: "Level 2",
					Icon:  icon.Folder,
					Children: []*ui.TreeNode{
						{
							ID:    "deep-1-1-1",
							Label: "Level 3",
							Icon:  icon.Folder,
							Children: []*ui.TreeNode{
								{
									ID:    "deep-1-1-1-1",
									Label: "Level 4",
									Icon:  icon.Folder,
									Children: []*ui.TreeNode{
										{ID: "deep-1-1-1-1-1", Label: "Level 5", Icon: icon.File},
										{ID: "deep-1-1-1-1-2", Label: "Level 5", Icon: icon.FileText},
									},
								},
							},
						},
					},
				},
				{ID: "deep-1-2", Label: "Level 2", Icon: icon.File},
			},
		},
		{
			ID:    "deep-2",
			Label: "Level 1",
			Icon:  icon.Folder,
			Children: []*ui.TreeNode{
				{ID: "deep-2-1", Label: "Level 2", Icon: icon.File},
				{ID: "deep-2-2", Label: "Level 2", Icon: icon.FileText},
			},
		},
	}
}
