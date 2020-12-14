package gui

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jamesrr39/taskrunner-app/taskrunner"
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

func (jobRunScene *JobRunScene) buildJobRunViewActionsBox() *gtk.VBox {
	text := "Back to Job Overview"

	goUpButton := gtk.NewButton()
	goUpButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_GO_UP, gtk.ICON_SIZE_LARGE_TOOLBAR))
	goUpButton.SetTooltipText(text)
	goUpButton.Clicked(func() {
		jobRunScene.TaskrunnerGUI.RenderScene(jobRunScene.TaskrunnerGUI.NewJobScene(jobRunScene.jobRun.Job))
	})

	goUpHbox := gtk.NewHBox(false, 5)
	goUpHbox.PackStart(goUpButton, false, false, 0)
	goUpHbox.PackStart(gtk.NewLabel(text), false, false, 0)

	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(goUpHbox, false, false, 0)

	return vbox
}

func (jobRunScene *JobRunScene) Content() gtk.IWidget {
	jobRun := jobRunScene.jobRun

	logLocationLbl := gtk.NewLabel(jobRunScene.JobRunsDAL.GetJobRunLogLocation(jobRun))
	logLocationLbl.SetSelectable(true)

	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(jobRunScene.buildTopBox(), false, false, 0)
	vbox.PackStart(gtk.NewLabel("Console Output:"), false, false, 0)
	vbox.PackStart(logLocationLbl, false, false, 0)
	vbox.PackStart(jobRunScene.buildTextareaScrollWindow(jobRun), true, true, 0)

	return vbox
}

func (jobRunScene *JobRunScene) buildTopBox() gtk.IBox {
	jobRun := jobRunScene.jobRun
	startDateTime := time.Unix(jobRun.StartTimestamp, 0)

	jobRunSummaryVbox := gtk.NewVBox(false, 0)
	triggerText := jobRun.Trigger
	if "" == triggerText {
		triggerText = "(No trigger)"
	} else {
		triggerText = "'" + triggerText + "'"
	}
	jobRunSummaryVbox.PackStart(gtk.NewLabel(fmt.Sprintf("Triggered by %s at %s (%s ago)", triggerText, startDateTime.Format(time.RFC822), GetTimeAgo(startDateTime))), false, false, 0)

	var durationStr string
	isFinished := (jobRun.EndTimestamp != 0)

	if !isFinished {
		durationStr = "Running for " + GetTimeAgo(startDateTime)
	} else {
		durationStr = "Ran for " + time.Unix(jobRun.EndTimestamp, 0).Sub(startDateTime).String()
	}

	jobRunSummaryVbox.PackStart(gtk.NewLabel(durationStr), false, false, 0)
	jobRunSummaryVbox.PackStart(gtk.NewLabel(jobRun.State.String()), false, false, 0)

	topHbox := gtk.NewHBox(false, 5)
	topHbox.PackStart(jobRunScene.buildJobRunViewActionsBox(), false, false, 30)
	topHbox.PackStart(jobRunSummaryVbox, true, true, 0)
	return topHbox
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
