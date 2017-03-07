package main

import (
	"os"
	"taskrunner-app/taskrunner-gtk/gui"
	"taskrunner-app/taskrunnerdal"

	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-gtk/gtk"

	"github.com/mattn/go-gtk/glib"

	"github.com/jamesrr39/goutil/user"
	"github.com/mattn/go-gtk/gdk"
)

var (
	taskrunnerDAL         *taskrunnerdal.TaskrunnerDAL
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

	taskrunnerGUI := gui.NewTaskrunnerGUI(taskrunnerDAL)
	taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())

	gtk.Main()

}

func setupApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.local/share/github.com/jamesrr39/taskrunner-app").
		String()

	taskrunnerApplication.Action(func(context *kingpin.ParseContext) error {
		expandedDir, err := user.ExpandUser(*taskrunnerDir)
		if nil != err {
			return err
		}

		taskrunnerDAL, err = taskrunnerdal.NewTaskrunnerDALAndEnsureDirectories(expandedDir)
		return err
	})
}
