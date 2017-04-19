package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/jamesrr39/goutil/user"
	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/jamesrr39/taskrunner-app/taskrunner-gtk/gui"
	"github.com/jamesrr39/taskrunner-app/taskrunnerdal"

	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var (
	taskrunnerDAL          *taskrunnerdal.TaskrunnerDAL
	taskrunnerApplication  *kingpin.Application
	jobLogMaxLines         *uint
	shouldMonitor          *bool
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
	setupApplicationFlags()
	kingpin.MustParse(taskrunnerApplication.Parse(os.Args[1:]))

	if *shouldMonitor {
		go monitor()
	}

	switch applicationMode {
	case ApplicationModeGUI:
		guiMain()
	case ApplicationModeRunJobHeadless:
		runJobHeadlessMain(*headlessJobNameArg, *headlessTriggerNameArg)
	}
}

func runJobHeadlessMain(jobName string, trigger string) {
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

	taskrunnerGUI := gui.NewTaskrunnerGUI(taskrunnerDAL, gui.TaskrunnerGUIOptions{*jobLogMaxLines, "/opt/taskrunner"})
	taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())

	gtk.Main()
}

func setupApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.local/share/github.com/jamesrr39/taskrunner-app").
		String()

	shouldMonitor = taskrunnerApplication.Flag("monitor", "print information about the number of goroutines used to the log output").
		Bool()

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
	meminfo := runtime.MemStats{}

	for {
		time.Sleep(time.Second)
		runtime.ReadMemStats(&meminfo)
		log.Printf("using %d goroutines (including 1 for monitoring).Memory %d, total: %d\n", runtime.NumGoroutine(), meminfo.Alloc, meminfo.TotalAlloc)
	}
}
