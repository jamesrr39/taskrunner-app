package gui

import (
	"strconv"
	"taskrunner"
	"time"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

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

func (taskrunnerGUI *TaskrunnerGUI) buildJobsSummaryTable() (*gtk.Table, error) {
	jobs, err := taskrunnerGUI.TaskrunnerInstance.GetAllJobs()
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
			table.AttachDefaults(gtk.NewLabel(GetTimeAgo(endDateTime)), 3, 4, uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(lastRun.State.String()), 4, 5, uint(index), uint(index+1))
		}
	}

	return table, nil

}
