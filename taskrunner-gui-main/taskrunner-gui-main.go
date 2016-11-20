package main

import (
	"os"
	"taskrunner"
	"taskrunner/gui"

	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-gtk/gtk"

	"github.com/mattn/go-gtk/glib"

	"github.com/mattn/go-gtk/gdk"
)

var (
	taskrunnerInstance    *taskrunner.TaskrunnerInstance
	taskrunnerApplication *kingpin.Application
)

func main() {

	taskrunnerApplication = kingpin.New("Taskrunner GUI", "gtk gui for taskrunner")
	setupApplicationFlags()
	kingpin.MustParse(taskrunnerApplication.Parse(os.Args[1:]))

	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	taskrunnerGUI := gui.NewTaskrunnerGUI(taskrunnerInstance)
	taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())

	gtk.Main()

}

func setupApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.local/share/taskrunner").
		String()

	taskrunnerApplication.Action(func(context *kingpin.ParseContext) error {
		var err error
		taskrunnerInstance, err = taskrunner.NewTaskrunnerInstanceAndEnsurePaths(*taskrunnerDir)
		return err
	})
}
