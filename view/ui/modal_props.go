package ui

import "github.com/a-h/templ"

type ModalDialogProps struct {
	ModalID     string
	ModalHeader templ.Component
	ModalFooter templ.Component
}

func (t *ModalDialogProps) GetModalID() string {
	return t.ModalID
}

func (t *ModalDialogProps) GetModalHeader() templ.Component {
	return t.ModalHeader
}

func (t *ModalDialogProps) GetModalFooter() templ.Component {
	return t.ModalFooter
}
