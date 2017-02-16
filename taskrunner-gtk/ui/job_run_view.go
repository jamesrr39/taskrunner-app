package gui

import (
	"log"
	"os"
	"strconv"
	"taskrunner-app/taskrunner"
	"time"

	"github.com/mattn/go-gtk/gtk"
)

type JobRunScene struct {
	*TaskrunnerGUI
	jobRun *taskrunner.JobRun
}

func (taskrunnerGUI *TaskrunnerGUI) NewJobRunScene(jobRun *taskrunner.JobRun) *JobRunScene {
	return &JobRunScene{taskrunnerGUI, jobRun}
}

func (jobRunScene *JobRunScene) IsCurrentlyRendered() bool {
	paneContentJobRunScene, ok := jobRunScene.TaskrunnerGUI.PaneContent.(*JobRunScene)
	if ok && paneContentJobRunScene.jobRun.Job.Id == jobRunScene.jobRun.Job.Id {
		return true
	}
	return false
}

func (jobRunScene *JobRunScene) Title() string {
	return "#" + strconv.Itoa(jobRunScene.jobRun.Id) + " :: " + jobRunScene.jobRun.Job.Name
}

func (jobRunScene *JobRunScene) Content() gtk.IWidget {
	jobRun := jobRunScene.jobRun

	isFinished := (jobRun.EndTimestamp != 0)

	vbox := gtk.NewVBox(false, 5)

	startDateTime := time.Unix(jobRun.StartTimestamp, 0)

	vbox.PackStart(gtk.IWidget(gtk.NewLabel("Started: "+startDateTime.String()+" ("+GetTimeAgo(startDateTime)+")")), false, false, 0)

	var durationStr string
	if !isFinished {
		durationStr = "Running for " + time.Now().Sub(startDateTime).String()
	} else {
		durationStr = "Ran for " + time.Unix(jobRun.EndTimestamp, 0).Sub(startDateTime).String()
	}

	vbox.PackStart(gtk.IWidget(gtk.NewLabel(durationStr)), false, false, 0)

	if isFinished {
		vbox.PackStart(gtk.NewLabel(jobRun.State.String()), false, false, 0)
	}

	logTextareaScrollWindow := gtk.NewScrolledWindow(nil, nil)
	logTextareaScrollWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	logTextareaScrollWindow.SetShadowType(gtk.SHADOW_IN)

	logTextarea := makeTextarea(jobRun)
	logTextareaScrollWindow.Add(logTextarea)

	vbox.PackStart(gtk.NewLabel("Console Output:"), false, false, 0)
	vbox.PackStart(gtk.IWidget(logTextareaScrollWindow), false, true, 0)

	return gtk.IWidget(vbox)

}

func makeTextarea(jobRun *taskrunner.JobRun) *gtk.TextView {
	logTextarea := gtk.NewTextView()
	logTextarea.SetEditable(false)
	logTextBuffer := logTextarea.GetBuffer()
	logFile, err := os.Open(jobRun.LogFilePath())
	if nil != err {
		logTextBuffer.InsertAtCursor("Couldn't read from " + jobRun.LogFilePath())
		log.Println("Couldn't read from " + jobRun.LogFilePath() + ". Error: " + err.Error())
	} else {
		defer logFile.Close()
		fillTextBufferFromFile(logTextBuffer, logFile, 200) //todo
	}
	return logTextarea
}
