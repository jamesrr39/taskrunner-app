package gui

import (
	"fmt"
	"log"
	"strconv"
	"taskrunner-app/taskrunner"
	"time"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gtk"
)

type JobRunScene struct {
	*TaskrunnerGUI
	jobRun *taskrunner.JobRun
}

func (taskrunnerGUI *TaskrunnerGUI) NewJobRunScene(jobRun *taskrunner.JobRun) *JobRunScene {
	return &JobRunScene{taskrunnerGUI, jobRun}
}

func (jobRunScene *JobRunScene) OnJobRunStatusChange(jobRun *taskrunner.JobRun) {
	if jobRun.Job.Id != jobRunScene.jobRun.Job.Id || jobRun.Id != jobRunScene.jobRun.Id {
		return
	}
	gdk.ThreadsEnter()
	jobRunScene.TaskrunnerGUI.RenderScene(jobRunScene.TaskrunnerGUI.NewJobRunScene(jobRun))
	gdk.ThreadsLeave()

}

func (jobRunScene *JobRunScene) Title() string {
	return "#" + strconv.FormatUint(jobRunScene.jobRun.Id, 10) + " :: " + jobRunScene.jobRun.Job.Name
}

func (jobRunScene *JobRunScene) Content() gtk.IWidget {
	jobRun := jobRunScene.jobRun

	isFinished := (jobRun.EndTimestamp != 0)

	vbox := gtk.NewVBox(false, 5)

	startDateTime := time.Unix(jobRun.StartTimestamp, 0)

	vbox.PackStart(gtk.NewLabel("Started: "+startDateTime.String()+" ("+GetTimeAgo(startDateTime)+" ago)"), false, false, 0)

	var durationStr string
	if !isFinished {
		durationStr = "Running for " + GetTimeAgo(startDateTime)
	} else {
		durationStr = "Ran for " + time.Unix(jobRun.EndTimestamp, 0).Sub(startDateTime).String()
	}

	vbox.PackStart(gtk.NewLabel(durationStr), false, false, 0)

	vbox.PackStart(gtk.NewLabel(jobRun.State.String()), false, false, 0)

	vbox.PackStart(gtk.NewLabel("Console Output:"), false, false, 0)
	vbox.PackStart(jobRunScene.buildTextareaScrollWindow(jobRun), true, true, 0)

	return vbox

}

func (jobRunScene *JobRunScene) buildTextareaScrollWindow(jobRun *taskrunner.JobRun) *gtk.ScrolledWindow {
	logTextareaScrollWindow := gtk.NewScrolledWindow(nil, nil)
	logTextareaScrollWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	logTextareaScrollWindow.SetShadowType(gtk.SHADOW_IN)

	logTextarea := jobRunScene.buildTextarea(jobRun)
	logTextareaScrollWindow.Add(logTextarea)

	return logTextareaScrollWindow
}

func (jobRunScene *JobRunScene) buildTextarea(jobRun *taskrunner.JobRun) *gtk.TextView {
	logTextarea := gtk.NewTextView()
	logTextarea.SetEditable(false)
	logTextBuffer := logTextarea.GetBuffer()

	logFile, err := jobRunScene.TaskrunnerGUI.TaskrunnerDAL.JobRunsDAL.GetJobRunLog(jobRun)
	if nil != err {
		errMessage := fmt.Sprintf("couldn't read job log file for %s. Error: %s", jobRun, err)
		logTextBuffer.InsertAtCursor(errMessage)
		log.Println(errMessage)
	} else {
		defer logFile.Close()
		fillTextBufferFromFile(logTextBuffer, logFile, jobRunScene.TaskrunnerGUI.options.JobLogMaxLines, jobRunScene.JobRunsDAL.GetJobRunLogLocation(jobRun))
	}
	return logTextarea
}
