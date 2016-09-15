package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"taskrunner"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var (
	taskrunnerInstance    *taskrunner.TaskrunnerInstance
	taskrunnerApplication *kingpin.Application
)

type TaskrunnerGUI struct {
	PaneContent gtk.IWidget
	Window      *gtk.Window
}

func main() {

	taskrunnerApplication = kingpin.New("Taskrunner GUI", "gtk gui for taskrunner")
	setupApplicationFlags()
	kingpin.MustParse(taskrunnerApplication.Parse(os.Args[1:]))

	gtk.Init(nil)
	makeWindow()
	gtk.Main()

}

func makeWindow() *TaskrunnerGUI {

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetSizeRequest(400, 400)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	window.SetTitle("Taskrunner (" + taskrunnerInstance.Basepath + ")")
	window.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	var paneContent gtk.IWidget = gtk.NewVBox(false, 0)

	taskrunnerGUI := &TaskrunnerGUI{PaneContent: paneContent, Window: window}
	taskrunnerGUI.RenderHomeScreen()

	window.ShowAll()

	return taskrunnerGUI
}

func (taskrunnerGUI *TaskrunnerGUI) RenderHomeScreen() {
	box := gtk.NewVBox(false, 0)
	var widget gtk.IWidget

	titleLabel := gtk.NewLabel("Taskrunner (" + taskrunnerInstance.Basepath + ")")
	box.Add(titleLabel)
	jobsTable, err := taskrunnerGUI.buildJobsSummaryTable()
	if nil != err {
		widget = gtk.NewLabel(err.Error())
	} else {
		widget = jobsTable
	}
	box.Add(widget)

	taskrunnerGUI.renderNewContent(box)
}

func setupApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.taskrunner").
		String()

	taskrunnerApplication.Action(func(context *kingpin.ParseContext) error {
		var err error
		taskrunnerInstance, err = taskrunner.NewTaskrunnerInstance(*taskrunnerDir)
		return err
	})
}

func (taskrunnerGUI *TaskrunnerGUI) renderNewContent(content gtk.IWidget) {

	taskrunnerGUI.PaneContent.Destroy()
	taskrunnerGUI.PaneContent = content
	taskrunnerGUI.Window.Add(taskrunnerGUI.PaneContent)
	taskrunnerGUI.Window.ShowAll()
}

func (taskrunnerGUI *TaskrunnerGUI) RenderJobRuns(job *taskrunner.Job) {
	var paneContent gtk.IWidget

	table, err := taskrunnerGUI.buildJobRunsTable(job)
	if nil != err {
		paneContent = gtk.NewLabel("Error fetching job runs for " + job.Name + ". Error: " + err.Error())
	} else {
		paneContent = table
	}
	taskrunnerGUI.renderNewContent(paneContent)

}

func (taskrunnerGUI *TaskrunnerGUI) RenderJobRun(jobRun *taskrunner.JobRun) {

	isFinished := (jobRun.EndTimestamp != 0)

	vbox := gtk.NewVBox(false, 5)
	vbox.Add(gtk.NewLabel("#" + strconv.Itoa(jobRun.Id) + " :: " + jobRun.Job.Name))

	startDateTime := time.Unix(jobRun.StartTimestamp, 0)
	vbox.Add(gtk.NewLabel("Started: " + startDateTime.String()))

	var durationStr string
	if !isFinished {
		durationStr = "Running for " + time.Now().Sub(startDateTime).String()
	} else {
		durationStr = "Ran for " + time.Unix(jobRun.EndTimestamp, 0).Sub(startDateTime).String()
	}

	vbox.Add(gtk.NewLabel(durationStr))

	if isFinished {
		vbox.Add(gtk.NewLabel("Finished successfully? " + strconv.FormatBool(jobRun.Successful)))
	}

	stdoutTextareaScrollWindow := gtk.NewScrolledWindow(nil, nil)
	stdoutTextareaScrollWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	stdoutTextareaScrollWindow.SetShadowType(gtk.SHADOW_IN)

	stdoutTextarea := gtk.NewTextView()
	stdoutTextBuffer := stdoutTextarea.GetBuffer()
	stdoutFile, err := os.Open(jobRun.StdoutLogPath())
	if nil != err {
		stdoutTextBuffer.InsertAtCursor("Couldn't read from " + jobRun.StdoutLogPath())
		log.Println("Couldn't read from " + jobRun.StdoutLogPath() + ". Error: " + err.Error())
	} else {
		defer stdoutFile.Close()
		fillTextBufferFromFile(stdoutTextBuffer, stdoutFile)
	}

	stdoutTextareaScrollWindow.Add(stdoutTextarea)
	vbox.Add(stdoutTextareaScrollWindow)

	var container gtk.IWidget = vbox
	taskrunnerGUI.renderNewContent(container)

}

func fillTextBufferFromFile(textBuffer *gtk.TextBuffer, fileReader io.Reader) {
	fileScanner := bufio.NewScanner(fileReader)
	linesRead := 0
	for fileScanner.Scan() && linesRead < 200 {
		textBuffer.InsertAtCursor(fileScanner.Text())
		linesRead++
	}
}

func (taskrunnerGUI *TaskrunnerGUI) buildJobsSummaryTable() (*gtk.Table, error) {
	jobs, err := taskrunnerInstance.Jobs()
	if nil != err {
		return nil, err
	}

	table := gtk.NewTable(uint(len(jobs)), 1, false)
	for index, job := range jobs {
		jobNameLabel := gtk.NewButtonWithLabel(job.Name)
		jobNameLabel.Clicked(func() {
			taskrunnerGUI.RenderJobRuns(job)
		})

		table.AttachDefaults(jobNameLabel, uint(1), uint(2), uint(index), uint(index+1))

		lastRunId := job.GetLastRunId()
		lastRun, err := job.GetRun(lastRunId)
		if nil != err {
			table.AttachDefaults(gtk.NewLabel(err.Error()), uint(2), uint(5), uint(index), uint(index+1))
		} else {
			table.AttachDefaults(gtk.NewLabel("#"+strconv.Itoa(lastRunId)), uint(2), uint(3), uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(strconv.Itoa(int(lastRun.EndTimestamp))), uint(3), uint(4), uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(strconv.FormatBool(lastRun.Successful)), uint(4), uint(5), uint(index), uint(index+1))
		}
	}

	return table, nil

}

func (taskrunnerGUI *TaskrunnerGUI) buildJobRunsTable(job *taskrunner.Job) (*gtk.Table, error) {
	runs, err := job.Runs()
	if nil != err {
		return nil, err
	}
	table := gtk.NewTable(uint(3), uint(len(runs)), false)
	for index, run := range runs {

		runIdButton := gtk.NewButtonWithLabel("#" + strconv.Itoa(run.Id))
		runIdButton.SetRelief(gtk.RELIEF_NONE)
		runIdButton.Clicked(func() {
			taskrunnerGUI.RenderJobRun(run)
		})
		startDateTime := time.Unix(run.StartTimestamp, 0)

		table.AttachDefaults(runIdButton, uint(1), uint(2), uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(startDateTime.String()), uint(2), uint(3), uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(strconv.FormatBool(run.Successful)), uint(3), uint(4), uint(index), uint(index+1))
	}
	return table, nil
}
