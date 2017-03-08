package gui

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

type ValidationLabel struct {
	label    *gtk.Label
	innerBox *gtk.EventBox
	Widget   *gtk.EventBox
	errColor *gdk.Color
}

func NewValidationLabel(errColor *gdk.Color) *ValidationLabel {
	outerBox := gtk.NewEventBox()

	return &ValidationLabel{Widget: outerBox, errColor: errColor}
}

func (validationLabel *ValidationLabel) SetText(text string) {
	if nil == validationLabel.innerBox {
		validationLabel.innerBox = gtk.NewEventBox()
		validationLabel.label = gtk.NewLabel(text)
		validationLabel.innerBox.Add(validationLabel.label)
		validationLabel.innerBox.ModifyBG(gtk.STATE_NORMAL, validationLabel.errColor)

		validationLabel.Widget.Add(validationLabel.innerBox)
	} else {
		validationLabel.label.SetText(text)
	}
	validationLabel.Widget.ShowAll()
}

func (validationLabel *ValidationLabel) Clear() {
	validationLabel.innerBox.Destroy()
}
