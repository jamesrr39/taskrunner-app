package gui

import (
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

type ValidationLabel struct {
	label    *gtk.Label
	innerBox *gtk.EventBox
	Widget   *gtk.VBox
	errColor *gdk.Color
}

func NewValidationLabel(errColor *gdk.Color) *ValidationLabel {
	outerBox := gtk.NewVBox(true, 0)

	return &ValidationLabel{Widget: outerBox, errColor: errColor}
}

func (validationLabel *ValidationLabel) SetText(text string) {
	if nil == validationLabel.innerBox {
		validationLabel.innerBox = gtk.NewEventBox()
		validationLabel.label = gtk.NewLabel(text)
		validationLabel.label.SetPadding(5, 5)
		validationLabel.innerBox.Add(validationLabel.label)
		validationLabel.innerBox.ModifyBG(gtk.STATE_NORMAL, validationLabel.errColor)
		validationLabel.Widget.PackStart(validationLabel.innerBox, false, false, 5)
	} else {
		validationLabel.label.SetText(text)
	}
	validationLabel.Widget.ShowAll()
}

func (validationLabel *ValidationLabel) Clear() {
	validationLabel.innerBox.Destroy()
}
