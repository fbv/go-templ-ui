package showcase

import (
	"github.com/a-h/templ"
	"github.com/fbv/go-templ-ui/view/ui"
)

type selectItem struct {
	ID    string
	Name  string
	Email string
}

type selectDialogExampleDS struct {
	Columns []ui.Column[selectItem]
	Rows    []selectItem
}

func (ds *selectDialogExampleDS) GetColumnCount() int           { return len(ds.Columns) }
func (ds *selectDialogExampleDS) GetColumnName(col int) string  { return ds.Columns[col].Name }
func (ds *selectDialogExampleDS) GetColumnStyle(col int) string { return ds.Columns[col].Style }
func (ds *selectDialogExampleDS) GetRowCount() int              { return len(ds.Rows) }
func (ds *selectDialogExampleDS) GetData(row, col int) templ.Component {
	return ds.Columns[col].GetData(ds.Rows[row])
}
func (ds *selectDialogExampleDS) GetRowMeta(row int) *ui.RowMeta { return &ui.RowMeta{} }
func (ds *selectDialogExampleDS) GetItemID(row int) string       { return ds.Rows[row].ID }

var selectDialogData = &selectDialogExampleDS{
	Columns: []ui.Column[selectItem]{
		{
			Name: "Имя",
			GetData: func(item selectItem) templ.Component {
				return ui.Text(item.Name)
			},
			Style: "w-40",
		},
		{
			Name: "Email",
			GetData: func(item selectItem) templ.Component {
				return ui.Text(item.Email)
			},
		},
	},
	Rows: []selectItem{
		{ID: "1", Name: "Иванов Иван", Email: "ivanov@example.com"},
		{ID: "2", Name: "Петров Пётр", Email: "petrov@example.com"},
		{ID: "3", Name: "Сидорова Анна", Email: "sidorova@example.com"},
		{ID: "4", Name: "Козлов Дмитрий", Email: "kozlov@example.com"},
		{ID: "5", Name: "Морозова Елена", Email: "morozova@example.com"},
		{ID: "6", Name: "Новиков Алексей", Email: "novikov@example.com"},
		{ID: "7", Name: "Волкова Мария", Email: "volkova@example.com"},
	},
}

var selectDialogProps = &ui.SelectDialogProps{
	ModalID:      "select-items-dlg",
	Title:        "Выберите сотрудников",
	SubmitURL:    "/showcase/api/selected",
	ConfirmLabel: "Отправить",
	CancelLabel:  "Отмена",
	Data:         selectDialogData,
}
