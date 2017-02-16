package gui

import (
	"taskrunner-app/taskrunner"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type TaskrunnerGUI struct {
	mainFrame           gtk.IBox
	PaneContent         Scene
	paneWidget          gtk.IWidget
	Window              *gtk.Window
	TaskrunnerInstance  *taskrunner.TaskrunnerInstance
	JobStatusChangeChan chan *taskrunner.JobRun
	titleLabel          *gtk.Label
}

type Scene interface {
	Title() string
	Content() gtk.IWidget
	IsCurrentlyRendered() bool
}

func NewTaskrunnerGUI(taskrunnerInstance *taskrunner.TaskrunnerInstance) *TaskrunnerGUI {

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetSizeRequest(800, 600)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	window.SetTitle("(Alpha) :: Taskrunner (" + taskrunnerInstance.Basepath + ")")
	window.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	mainFrame := gtk.NewVBox(false, 10)

	titleLabel := gtk.NewLabel("")
	titleLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	taskrunnerGUI := &TaskrunnerGUI{
		mainFrame:           gtk.IBox(mainFrame),
		Window:              window,
		TaskrunnerInstance:  taskrunnerInstance,
		JobStatusChangeChan: make(chan *taskrunner.JobRun),
		titleLabel:          titleLabel,
	}

	mainFrame.PackStart(taskrunnerGUI.buildToolbar(), false, false, 0)
	window.Add(mainFrame)

	return taskrunnerGUI
}

func (taskrunnerGUI *TaskrunnerGUI) buildToolbar() gtk.IWidget {
	homeButton := gtk.NewButton()
	homeButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_HOME, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	homeButton.Clicked(func() {
		taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())
	})

	newJobButton := gtk.NewButton()
	newJobButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_ADD, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	newJobButton.Clicked(func() {
		newJob, err := taskrunner.NewJob("New Job", "", taskrunner.Script("#!/bin/bash\n\n"), taskrunnerGUI.TaskrunnerInstance)
		if nil != err {
			panic(err)
		}
		taskrunnerGUI.RenderScene(taskrunnerGUI.NewEditJobView(newJob))

	})

	hbox := gtk.NewHBox(false, 0)
	hbox.PackStart(taskrunnerGUI.titleLabel, true, true, 3)
	hbox.PackEnd(homeButton, false, false, 0)
	hbox.PackEnd(newJobButton, false, false, 0)

	eventBox := gtk.NewEventBox()
	eventBox.Add(hbox)
	eventBox.ModifyBG(gtk.STATE_NORMAL, titleBlue())

	return eventBox
}

func (taskrunnerGUI *TaskrunnerGUI) RenderScene(scene Scene) {
	if nil != taskrunnerGUI.paneWidget {
		taskrunnerGUI.paneWidget.Destroy()
	}
	taskrunnerGUI.PaneContent = scene

	taskrunnerGUI.titleLabel.SetText(scene.Title())

	taskrunnerGUI.paneWidget = scene.Content()

	taskrunnerGUI.mainFrame.Add(taskrunnerGUI.paneWidget)

	taskrunnerGUI.Window.ShowAll()
}
