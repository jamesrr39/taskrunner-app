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
	jobRun   *taskrunner.JobRun
	isClosed bool
}

func (taskrunnerGUI *TaskrunnerGUI) NewJobRunScene(jobRun *taskrunner.JobRun) *JobRunScene {
	return &JobRunScene{taskrunnerGUI, jobRun, true}
}

func (jobRunScene *JobRunScene) OnClose() {
	jobRunScene.isClosed = true
}

func (jobRunScene *JobRunScene) OnShow() {
	jobRunScene.isClosed = false

	go func(jobRunScene *JobRunScene) {
		renderedJobRun := jobRunScene.jobRun
		for {
			if jobRunScene.isClosed {
				return
			}
			select {
			case jobRun := <-jobRunScene.TaskrunnerGUI.JobStatusChangeChan:
				log.Printf("catching in job run view job run id: %d. Current job id: %d\n", jobRun.Job.Id, renderedJobRun.Id)

				gdk.ThreadsEnter()
				jobRunScene.TaskrunnerGUI.RenderScene(jobRunScene.TaskrunnerGUI.NewJobRunScene(renderedJobRun)) // todo check still on this screen interface CurrentSceneRendered
				gdk.ThreadsLeave()
			default:
				// non-blocking receive
			}
		}
	}(jobRunScene)
}

func (jobRunScene *JobRunScene) Title() string {
	return "#" + strconv.FormatUint(jobRunScene.jobRun.Id, 10) + " :: " + jobRunScene.jobRun.Job.Name
}

func (jobRunScene *JobRunScene) Content() gtk.IWidget {
	jobRun := jobRunScene.jobRun

	isFinished := (jobRun.EndTimestamp != 0)

	vbox := gtk.NewVBox(false, 5)

	startDateTime := time.Unix(jobRun.StartTimestamp, 0)

	vbox.PackStart(gtk.NewLabel("Started: "+startDateTime.String()+" ("+GetTimeAgo(startDateTime)+")"), false, false, 0)

	var durationStr string
	if !isFinished {
		durationStr = "Running for " + time.Now().Sub(startDateTime).String()
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
		fillTextBufferFromFile(logTextBuffer, logFile, 200) //todo
	}
	return logTextarea
}
