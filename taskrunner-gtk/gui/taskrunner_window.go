package gui

import (
	"taskrunner-app/taskrunner"
	"taskrunner-app/taskrunnerdal"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type TaskrunnerGUI struct {
	mainFrame   gtk.IBox
	PaneContent Scene
	paneWidget  gtk.IWidget
	Window      *gtk.Window
	*taskrunnerdal.TaskrunnerDAL
	JobStatusChangeChan chan *taskrunner.JobRun // job runs
	titleLabel          *gtk.Label
	options             TaskrunnerGUIOptions
}

func NewTaskrunnerGUI(taskrunnerDAL *taskrunnerdal.TaskrunnerDAL, options TaskrunnerGUIOptions) *TaskrunnerGUI {

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetSizeRequest(800, 600)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	window.SetTitle("(Alpha) :: Taskrunner (" + taskrunnerDAL.String() + ")")
	window.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	mainFrame := gtk.NewVBox(false, 10)

	titleLabel := gtk.NewLabel("")
	titleLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	taskrunnerGUI := &TaskrunnerGUI{
		mainFrame:           gtk.IBox(mainFrame),
		Window:              window,
		TaskrunnerDAL:       taskrunnerDAL,
		JobStatusChangeChan: make(chan *taskrunner.JobRun),
		titleLabel:          titleLabel,
		options:             options,
	}

	mainFrame.PackStart(buildToolbar(taskrunnerGUI), false, false, 0)
	window.Add(mainFrame)

	return taskrunnerGUI
}

func (taskrunnerGUI *TaskrunnerGUI) RenderScene(scene Scene) {
	if nil != taskrunnerGUI.paneWidget {
		taskrunnerGUI.PaneContent.OnClose()
		taskrunnerGUI.paneWidget.Destroy()
	}
	taskrunnerGUI.PaneContent = scene

	taskrunnerGUI.titleLabel.SetText(scene.Title())

	taskrunnerGUI.paneWidget = scene.Content()

	taskrunnerGUI.mainFrame.Add(taskrunnerGUI.paneWidget)

	taskrunnerGUI.Window.ShowAll()
	scene.OnShow()
}
