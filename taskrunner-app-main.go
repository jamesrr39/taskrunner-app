package main

import (
	"log"
	"os"
	"time"

	"github.com/jamesrr39/goutil/userextra"
	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/jamesrr39/taskrunner-app/taskrunner-gtk/gui"
	"github.com/jamesrr39/taskrunner-app/taskrunnerdal"
	"github.com/jamesrr39/taskrunner-app/triggers"

	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var (
	taskrunnerDAL          *taskrunnerdal.TaskrunnerDAL
	taskrunnerApplication  *kingpin.Application
	jobLogMaxLines         *uint
	applicationMode        ApplicationMode
	headlessJobNameArg     *string
	headlessTriggerNameArg *string
)

type ApplicationMode int

const (
	ApplicationModeGUI ApplicationMode = iota
	ApplicationModeRunJobHeadless
)

func main() {
	taskrunnerApplication = kingpin.New("Taskrunner GUI", "gtk gui for taskrunner")
	parseApplicationFlags()
	kingpin.MustParse(taskrunnerApplication.Parse(os.Args[1:]))

	switch applicationMode {
	case ApplicationModeGUI:
		guiMain()
	case ApplicationModeRunJobHeadless:
		runJobHeadlessMain(taskrunnerDAL, *headlessJobNameArg, *headlessTriggerNameArg)
	}
}

func runJobHeadlessMain(taskrunnerDAL *taskrunnerdal.TaskrunnerDAL, jobName string, trigger string) {
	job, err := taskrunnerDAL.GetJobByName(jobName)
	if nil != err {
		log.Fatalf("Error finding job with name: '%s'. Error: %s\n", jobName, err)
	}

	jobRun := job.NewJobRun(taskrunner.TriggerType(trigger))
	err = taskrunnerDAL.JobRunsDAL.CreateAndRun(jobRun, nil)
	if nil != err {
		log.Fatalf("Error running job with name: '%s' (id %d, run Id %d).\nError: %s\n", jobName, job.Id, jobRun.Id, err)
	}
}

func guiMain() {

	glib.ThreadInit(nil)
	gdk.ThreadsInit()
	gdk.ThreadsEnter()
	gtk.Init(nil)

	udevDAL := triggers.NewUdevRulesDAL("/etc/udev/rules.d")
	taskrunnerGUI := gui.NewTaskrunnerGUI(taskrunnerDAL, udevDAL, gui.TaskrunnerGUIOptions{*jobLogMaxLines, "/opt/taskrunner"}) // todo mock
	taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())

	gtk.Main()
}

func parseApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.local/share/github.com/jamesrr39/taskrunner-app").
		String()

	jobLogMaxLines = taskrunnerApplication.Flag("job-log-max-lines", "maximum number of lines to display in the job output. Lines after that are not shown in the UI, but the UI indicates where the whole log file is instead").
		Default("10000").
		Uint()

	headlessJobNameArg = taskrunnerApplication.Flag("run-job", "job name to be run headlessly").String()
	headlessTriggerNameArg = taskrunnerApplication.Flag("trigger", "name of trigger to be recorded in the job run").String()

	taskrunnerApplication.Action(func(context *kingpin.ParseContext) error {
		if "" == *headlessJobNameArg {
			applicationMode = ApplicationModeGUI
		} else {
			applicationMode = ApplicationModeRunJobHeadless
		}

		expandedDir, err := userextra.ExpandUser(*taskrunnerDir)
		if nil != err {
			return err
		}

		taskrunnerDAL, err = taskrunnerdal.NewTaskrunnerDALAndEnsureDirectories(expandedDir, providesNow)
		if nil != err {
			return err
		}

		return nil
	})
}

func providesNow() time.Time {
	return time.Now()
}
