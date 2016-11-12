package gui

import (
	"taskrunner"

	"log"
	"strconv"
	"time"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

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
		go func(job *taskrunner.Job, taskrunnerGUI *TaskrunnerGUI) {
			job.Run("GUI")
		}(job, taskrunnerGUI)

		taskrunnerGUI.RenderJobRuns(job)
	}, job)

	configureButton := gtk.NewButton()
	configureButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_EDIT, gtk.ICON_SIZE_LARGE_TOOLBAR))
	configureButton.Clicked(func(ctx *glib.CallbackContext) {
		job, ok := ctx.Data().(*taskrunner.Job)
		if !ok {
			panic("couldn't convert to job")
		}
		log.Printf("job for configure edit job view: %v\n", job)
		taskrunnerGUI.renderNewContent(taskrunnerGUI.makeConfigureEditJobView(job))
	}, job)

	hbox.PackStart(runButton, false, false, 0)
	hbox.PackEnd(configureButton, false, false, 0)
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

	taskrunnerGUI.renderNewContent(box)

}

func (taskrunnerGUI *TaskrunnerGUI) buildJobRunsTable(job *taskrunner.Job) (*gtk.Table, error) {
	runs, err := job.GetRuns()
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
		table.AttachDefaults(gtk.NewLabel(run.State.String()), 3, 4, uint(index), uint(index+1))
	}
	return table, nil
}
