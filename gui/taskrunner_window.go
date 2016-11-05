package gui

import (
	"taskrunner"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type TaskrunnerGUI struct {
	mainFrame          gtk.IBox
	PaneContent        gtk.IWidget
	Window             *gtk.Window
	TaskrunnerInstance *taskrunner.TaskrunnerInstance
}

func NewTaskrunnerGUI(taskrunnerInstance *taskrunner.TaskrunnerInstance) *TaskrunnerGUI {

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetSizeRequest(400, 400)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	window.SetTitle("Taskrunner (" + taskrunnerInstance.Basepath + ")")
	window.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	mainFrame := gtk.NewVBox(false, 0)

	paneContent := gtk.NewVBox(false, 0)

	taskrunnerGUI := &TaskrunnerGUI{
		mainFrame:          gtk.IBox(mainFrame),
		PaneContent:        gtk.IWidget(paneContent),
		Window:             window,
		TaskrunnerInstance: taskrunnerInstance,
	}

	mainFrame.PackStart(taskrunnerGUI.buildToolbar(), false, false, 0)
	mainFrame.PackStart(gtk.IWidget(paneContent), true, true, 0)
	window.Add(mainFrame)

	return taskrunnerGUI
}

func (taskrunnerGUI *TaskrunnerGUI) buildToolbar() gtk.IWidget {
	homeButton := gtk.NewButton()
	homeButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_HOME, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	homeButton.Clicked(taskrunnerGUI.RenderHomeScreen)

	newJobButton := gtk.NewButton()
	newJobButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_ADD, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	newJobButton.Clicked(func() {
		taskrunnerGUI.renderNewContent(taskrunnerGUI.makeConfigureCreateJobView())
	})

	hbox := gtk.NewHBox(true, 0)

	hbox.PackStart(homeButton, false, false, 0)
	hbox.PackStart(newJobButton, false, false, 0)

	eventBox := gtk.NewEventBox()
	eventBox.Add(hbox)
	eventBox.ModifyBG(gtk.STATE_NORMAL, gdk.NewColorRGB(uint8(223), uint8(223), uint8(223)))

	return eventBox
}

func (taskrunnerGUI *TaskrunnerGUI) renderNewContent(content gtk.IWidget) {

	taskrunnerGUI.PaneContent.Destroy()
	taskrunnerGUI.PaneContent = content
	taskrunnerGUI.mainFrame.Add(taskrunnerGUI.PaneContent)

	taskrunnerGUI.Window.ShowAll()
}
