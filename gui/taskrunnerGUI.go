package gui

import (
	"taskrunner"

	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type TaskrunnerGUI struct {
	mainFrame          gtk.IBox
	PaneContent        gtk.IWidget
	Window             *gtk.Window
	TaskrunnerInstance *taskrunner.TaskrunnerInstance
}

type WidgetRenderer struct {
	renderer func() gtk.IWidget
}

func NewWidgetRenderer(contentRenderer func() gtk.IWidget) *WidgetRenderer {
	return &WidgetRenderer{
		renderer: contentRenderer,
	}
}

func NewTaskrunnerGUI(taskrunnerInstance *taskrunner.TaskrunnerInstance) *TaskrunnerGUI {

	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetSizeRequest(400, 400)
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	})
	window.SetTitle("Taskrunner (" + taskrunnerInstance.Basepath + ")")
	window.ModifyBG(gtk.STATE_NORMAL, gdk.NewColor("white"))

	mainFrame := gtk.NewVBox(false, 0)

	paneContent := gtk.NewVBox(false, 0)

	taskrunnerGUI := &TaskrunnerGUI{
		mainFrame:          gtk.IBox(mainFrame),
		PaneContent:        gtk.IWidget(paneContent),
		Window:             window,
		TaskrunnerInstance: taskrunnerInstance,
	}

	mainFrame.PackStart(taskrunnerGUI.buildToolbar(), false, false, 0)
	mainFrame.PackStart(gtk.IWidget(paneContent), true, true, 0)
	window.Add(mainFrame)

	return taskrunnerGUI
}

func (taskrunnerGUI *TaskrunnerGUI) buildToolbar() gtk.IWidget {
	homeButton := gtk.NewButton()
	homeButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_HOME, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	homeButton.Clicked(taskrunnerGUI.RenderHomeScreen)

	newJobButton := gtk.NewButton()
	newJobButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_ADD, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	newJobButton.Clicked(func() {
		taskrunnerGUI.renderNewContent(taskrunnerGUI.buildNewJobView())
	})

	hbox := gtk.NewHBox(true, 0)

	hbox.PackStart(homeButton, false, false, 0)
	hbox.PackStart(newJobButton, false, false, 0)

	return hbox
}

func (taskrunnerGUI *TaskrunnerGUI) renderNewContent(content gtk.IWidget) {

	taskrunnerGUI.PaneContent.Destroy()
	taskrunnerGUI.PaneContent = content
	taskrunnerGUI.mainFrame.Add(taskrunnerGUI.PaneContent)
	taskrunnerGUI.Window.ShowAll()
}

func (taskrunnerGUI *TaskrunnerGUI) RenderJobRuns(job *taskrunner.Job) {
	box := gtk.NewVBox(false, 5)

	box.PackStart(gtk.NewLabel("Runs for "+job.Name), false, false, 5)

	hbox := gtk.NewHBox(true, 0)
	runButton := gtk.NewButton()
	runButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_LARGE_TOOLBAR))
	runButton.Clicked(func(ctx *glib.CallbackContext) {
		job, ok := ctx.Data().(*taskrunner.Job)
		if !ok {
			panic("couldn't convert to job")
		}
		go job.Run("GUI")
	}, job)

	hbox.PackStart(runButton, false, false, 0)
	box.PackStart(hbox, false, false, 0)
	var listing gtk.IWidget

	if 0 == job.GetLastRunId() {
		listing = gtk.NewLabel("No runs yet...")
	} else {
		table, err := taskrunnerGUI.buildJobRunsTable(job)
		if nil != err {
			listing = gtk.NewLabel("Error fetching job runs for " + job.Name + ". Error: " + err.Error())
		} else {
			listing = table
		}
	}
	box.PackStart(listing, false, false, 5)

	taskrunnerGUI.renderNewContent(gtk.IWidget(box))

}

func (taskrunnerGUI *TaskrunnerGUI) RenderHomeScreen() {
	box := gtk.NewVBox(false, 0)
	var jobsTableWidget gtk.IWidget

	titleLabel := gtk.IWidget(gtk.NewLabel("Taskrunner (" + taskrunnerGUI.TaskrunnerInstance.Basepath + ")"))
	box.PackStart(titleLabel, false, true, 5)
	jobsTable, err := taskrunnerGUI.buildJobsSummaryTable()
	if nil != err {
		jobsTableWidget = gtk.NewLabel(err.Error())
	} else {
		jobsTableWidget = jobsTable
	}
	box.PackStart(jobsTableWidget, false, false, 5)

	taskrunnerGUI.renderNewContent(box)
}

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
		vbox.PackStart(gtk.IWidget(gtk.NewLabel("Finished successfully? "+strconv.FormatBool(jobRun.Successful))), false, false, 0)
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
	vbox.PackStart(gtk.IWidget(stdoutTextareaScrollWindow), true, true, 0)

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
	jobs, err := taskrunnerGUI.TaskrunnerInstance.Jobs()
	if nil != err {
		return nil, err
	}

	table := gtk.NewTable(uint(len(jobs)), 1, false)
	for index, job := range jobs {
		jobNameLabel := gtk.NewButtonWithLabel(job.Name)
		jobNameLabel.Clicked(func(ctx *glib.CallbackContext) {
			job, ok := ctx.Data().(*taskrunner.Job)
			if !ok {
				panic("casting job didn't work")
			}
			taskrunnerGUI.RenderJobRuns(job)
		}, job)

		table.AttachDefaults(jobNameLabel, 1, 2, uint(index), uint(index+1))

		lastRunId := job.GetLastRunId()
		if 0 == lastRunId {
			table.AttachDefaults(gtk.NewLabel("No runs yet..."), 2, 5, uint(index), uint(index+1))
			continue
		}

		lastRun, err := job.GetRun(lastRunId)
		if nil != err {
			table.AttachDefaults(gtk.NewLabel(err.Error()), 2, 5, uint(index), uint(index+1))
		} else {
			endDateTime := time.Unix(int64(lastRun.EndTimestamp), 0)
			table.AttachDefaults(gtk.NewLabel("#"+strconv.Itoa(lastRunId)), 2, 3, uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(endDateTime.Format(time.RFC850)), 3, 4, uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(strconv.FormatBool(lastRun.Successful)), 4, 5, uint(index), uint(index+1))
		}
	}

	return table, nil

}

func (taskrunnerGUI *TaskrunnerGUI) buildJobRunsTable(job *taskrunner.Job) (*gtk.Table, error) {
	runs, err := job.Runs()
	if nil != err {
		return nil, err
	}
	table := gtk.NewTable(3, uint(len(runs)), false)
	for index, run := range runs {

		runIdButton := gtk.NewButtonWithLabel("#" + strconv.Itoa(run.Id))
		runIdButton.SetRelief(gtk.RELIEF_NONE)
		runIdButton.Clicked(func(context *glib.CallbackContext) {
			if run, ok := context.Data().(*taskrunner.JobRun); ok {
				taskrunnerGUI.RenderJobRun(run)
			} else {
				errorMessage := "Couldn't get job clicked on"
				taskrunnerGUI.renderNewContent(gtk.IWidget(gtk.NewLabel(errorMessage)))
				log.Println(errorMessage)
			}
		}, run)
		startDateTime := time.Unix(run.StartTimestamp, 0)

		table.AttachDefaults(runIdButton, uint(1), 2, uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(startDateTime.String()), 2, 3, uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(strconv.FormatBool(run.Successful)), 3, 4, uint(index), uint(index+1))
	}
	return table, nil
}
