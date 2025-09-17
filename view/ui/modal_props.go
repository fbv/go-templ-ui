package ui

import "github.com/a-h/templ"

// ModalDialogProps provides concrete props for the Modal templ component.
// It satisfies the ModalProps interface by exposing an id, header, and footer.
type ModalDialogProps struct {
	ModalID     string
	ModalHeader templ.Component
	ModalFooter templ.Component
}

// GetModalID returns the DOM id used to target/toggle the modal element.
func (t *ModalDialogProps) GetModalID() string {
	return t.ModalID
}

// GetModalHeader returns the component rendered in the modal header section.
func (t *ModalDialogProps) GetModalHeader() templ.Component {
	return t.ModalHeader
}

// GetModalFooter returns the component rendered in the modal footer section.
func (t *ModalDialogProps) GetModalFooter() templ.Component {
	return t.ModalFooter
}
