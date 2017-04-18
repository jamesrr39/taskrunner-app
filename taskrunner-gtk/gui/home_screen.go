package gui

import (
	"strconv"
	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"time"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type HomeScene struct {
	*TaskrunnerGUI
}

func (taskrunner *TaskrunnerGUI) NewHomeScene() *HomeScene {
	return &HomeScene{taskrunner}
}

func (homeScreen *HomeScene) OnJobRunStatusChange(jobRun *taskrunner.JobRun) {
	gdk.ThreadsEnter()
	homeScreen.TaskrunnerGUI.RenderScene(homeScreen.TaskrunnerGUI.NewHomeScene())
	gdk.ThreadsLeave()
}

func (homeScreen *HomeScene) Title() string {
	return "Taskrunner"
}

func (homeScreen *HomeScene) Content() gtk.IWidget {
	box := gtk.NewVBox(false, 0)
	var jobsTableWidget gtk.IWidget

	jobsTable, err := homeScreen.buildJobsSummaryTable()
	if nil != err {
		jobsTableWidget = gtk.NewLabel(err.Error())
	} else {
		swin := gtk.NewScrolledWindow(nil, nil)
		swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
		swin.SetShadowType(gtk.SHADOW_IN)

		innerVbox := gtk.NewVBox(false, 0)
		innerVbox.PackStart(jobsTable, false, false, 0)
		swin.AddWithViewPort(innerVbox)

		jobsTableWidget = swin
	}

	box.PackStart(jobsTableWidget, true, true, 0)

	return box
}

func (homeScreen *HomeScene) buildJobsSummaryTable() (*gtk.Table, error) {
	jobs, err := homeScreen.TaskrunnerGUI.TaskrunnerDAL.GetAll()
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
			homeScreen.TaskrunnerGUI.RenderScene(homeScreen.TaskrunnerGUI.NewJobScene(job))
		}, job)

		table.AttachDefaults(jobNameLabel, 1, 2, uint(index), uint(index+1))

		lastJobRun, err := homeScreen.TaskrunnerGUI.TaskrunnerDAL.JobRunsDAL.GetLastRun(job)
		if nil != err {
			table.AttachDefaults(gtk.NewLabel(err.Error()), 2, 5, uint(index), uint(index+1))
		} else if nil == lastJobRun {
			table.AttachDefaults(gtk.NewLabel("No runs yet..."), 2, 5, uint(index), uint(index+1))
			continue
		} else {
			// todo handle in progress
			endDateTime := time.Unix(lastJobRun.EndTimestamp, 0)
			table.AttachDefaults(gtk.NewLabel("#"+strconv.FormatUint(uint64(lastJobRun.Id), 10)), 2, 3, uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(GetTimeAgo(endDateTime)+" ago"), 3, 4, uint(index), uint(index+1))
			table.AttachDefaults(gtk.NewLabel(lastJobRun.State.String()), 4, 5, uint(index), uint(index+1))
		}
	}

	return table, nil

}
