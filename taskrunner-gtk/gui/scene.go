package gui

import (
	"github.com/mattn/go-gtk/gtk"
)

type Scene interface {
	Title() string
	Content() gtk.IWidget
	OnClose() // cleanup - ending listening for events
	OnShow()
}
