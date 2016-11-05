package gui

import (
	"log"
	"os"
	"strconv"
	"taskrunner"
	"time"

	"github.com/mattn/go-gtk/gtk"
)

func (taskrunnerGUI *TaskrunnerGUI) RenderJobRun(jobRun *taskrunner.JobRun) {

	isFinished := (jobRun.EndTimestamp != 0)

	vbox := gtk.NewVBox(false, 5)
	vbox.PackStart(gtk.IWidget(gtk.NewLabel("#"+strconv.Itoa(jobRun.Id)+" :: "+jobRun.Job.Name)), false, false, 0)

	startDateTime := time.Unix(jobRun.StartTimestamp, 0)
	vbox.PackStart(gtk.IWidget(gtk.NewLabel("Started: "+startDateTime.String())), false, false, 0)

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

	logTextarea := gtk.NewTextView()
	logTextBuffer := logTextarea.GetBuffer()
	logFile, err := os.Open(jobRun.LogFilePath())
	if nil != err {
		logTextBuffer.InsertAtCursor("Couldn't read from " + jobRun.LogFilePath())
		log.Println("Couldn't read from " + jobRun.LogFilePath() + ". Error: " + err.Error())
	} else {
		defer logFile.Close()
		fillTextBufferFromFile(logTextBuffer, logFile, 200) //todo
	}

	logTextareaScrollWindow.Add(logTextarea)
	vbox.PackStart(gtk.IWidget(logTextareaScrollWindow), true, true, 0)

	var container gtk.IWidget = vbox
	taskrunnerGUI.renderNewContent(container)

}