package ui

type SelectDialogDataSource interface {
	TableDataSource
	GetItemID(row int) string
}

type SelectDialogProps struct {
	ModalID       string
	Title         string
	SubmitURL     string
	ConfirmLabel  string
	CancelLabel   string
	Data          SelectDialogDataSource
	FormAttrs     map[string]string
}

func (p *SelectDialogProps) confirmLabel() string {
	if p.ConfirmLabel != "" {
		return p.ConfirmLabel
	}
	return "OK"
}

func (p *SelectDialogProps) cancelLabel() string {
	if p.CancelLabel != "" {
		return p.CancelLabel
	}
	return "Отмена"
}
