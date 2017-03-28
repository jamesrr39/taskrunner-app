package main

import (
	"log"
	"os"
	"runtime"
	"taskrunner-app/taskrunner-gtk/gui"
	"taskrunner-app/taskrunnerdal"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-gtk/gtk"

	"github.com/mattn/go-gtk/glib"

	"github.com/jamesrr39/goutil/user"
	"github.com/mattn/go-gtk/gdk"
)

var (
	taskrunnerDAL         *taskrunnerdal.TaskrunnerDAL
	taskrunnerApplication *kingpin.Application
	jobLogMaxLines        *uint
)

func main() {

	taskrunnerApplication = kingpin.New("Taskrunner GUI", "gtk gui for taskrunner")
	setupApplicationFlags()
	kingpin.MustParse(taskrunnerApplication.Parse(os.Args[1:]))

	go monitor()

	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	taskrunnerGUI := gui.NewTaskrunnerGUI(taskrunnerDAL, gui.TaskrunnerGUIOptions{*jobLogMaxLines})
	taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())

	gtk.Main()
}

func setupApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.local/share/github.com/jamesrr39/taskrunner-app").
		String()

	jobLogMaxLines = taskrunnerApplication.Flag("job-log-max-lines", "maximum number of lines to display in the job output. Lines after that are not shown in the UI, but the UI indicates where the whole log file is instead").
		Default("10000").
		Uint()

	taskrunnerApplication.Action(func(context *kingpin.ParseContext) error {
		expandedDir, err := user.ExpandUser(*taskrunnerDir)
		if nil != err {
			return err
		}

		taskrunnerDAL, err = taskrunnerdal.NewTaskrunnerDALAndEnsureDirectories(expandedDir)
		if nil != err {
			return err
		}

		return nil
	})
}

func monitor() {
	for {
		time.Sleep(time.Second)
		log.Printf("using %d goroutines\n", runtime.NumGoroutine())
	}
}
