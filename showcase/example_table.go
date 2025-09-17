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
		{
			Name: "",
			GetData: func(row string) templ.Component {

				return ui.ButtonGroup(
					[]ui.ButtonGroupItem{
						{
							SVGIcon: &ui.SVGIconProps{
								SVG: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-4">
				  <path stroke-linecap="round" stroke-linejoin="round" d="m16.862 4.487 1.687-1.688a1.875 1.875 0 1 1 2.652 2.652L10.582 16.07a4.5 4.5 0 0 1-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 0 1 1.13-1.897l8.932-8.931Zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0 1 15.75 21H5.25A2.25 2.25 0 0 1 3 18.75V8.25A2.25 2.25 0 0 1 5.25 6H10" />
				</svg>`,
							},
							Title: "Изменить сущность...",
							Attrs: map[string]string{"data-twe-toggle": "modal", "data-twe-target": "#edit-conf-mdl",
								"data-twe-ripple-init": "true", " data-twe-ripple-color": "light"},
						},
						{
							SVGIcon: &ui.SVGIconProps{
								SVG: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="size-4">
  <path stroke-linecap="round" stroke-linejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
</svg>`,
							},
							Title: "Удалить сущность...",
							Attrs: map[string]string{"data-twe-toggle": "modal", "data-twe-target": "#delete-conf-mdl",
								"data-twe-ripple-init": "true", " data-twe-ripple-color": "light"},
						},
					},
				)
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
	RowModifier: &ui.RowModifier[string]{
		OnClick: func(row int, d string) string {
			return fmt.Sprintf("#row=%d", row)
		},
	},
}

var confirmationModal *ui.ModalDialogProps = &ui.ModalDialogProps{
	ModalID:     "delete-conf-mdl",
	ModalHeader: modalDlgHeader("Удаление записи"),
	ModalFooter: deleteRowConfirmationDlgFooter(),
}

var editableModal *ui.ModalDialogProps = &ui.ModalDialogProps{
	ModalID:     "edit-conf-mdl",
	ModalHeader: modalDlgHeader("Редактирование записи"),
	ModalFooter: editRowConfirmationDlgFooter(),
}
