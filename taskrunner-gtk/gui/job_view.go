package gui

import (
	"log"
	"taskrunner-app/taskrunner"

	"strconv"
	"time"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type JobScene struct {
	*TaskrunnerGUI
	Job      *taskrunner.Job
	isClosed bool
}

func (taskrunnerGUI *TaskrunnerGUI) NewJobScene(job *taskrunner.Job) *JobScene {
	return &JobScene{taskrunnerGUI, job, true}
}

func (jobScene *JobScene) OnClose() {
	jobScene.isClosed = true
}

func (jobScene *JobScene) OnShow() {
	jobScene.isClosed = false

	go func(jobScene *JobScene) {
		renderedJob := jobScene.Job
		for {
			if jobScene.isClosed {
				return
			}
			select {
			case <-jobScene.TaskrunnerGUI.JobStatusChangeChan:
				gdk.ThreadsEnter()
				jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.NewJobScene(renderedJob)) // todo check still on this screen interface CurrentSceneRendered
				gdk.ThreadsLeave()
			default:
				// non-blocking receive
			}
		}
	}(jobScene)
}

func (jobScene *JobScene) Title() string {
	return "Runs for " + jobScene.Job.Name
}

func (jobScene *JobScene) Content() gtk.IWidget {

	box := gtk.NewVBox(false, 5)

	hbox := gtk.NewHBox(true, 0)
	runButton := gtk.NewButton()
	runButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_LARGE_TOOLBAR))
	runButton.Clicked(func(ctx *glib.CallbackContext) {
		job, ok := ctx.Data().(*taskrunner.Job)
		if !ok {
			panic("couldn't convert to job")
		}
		go func(job *taskrunner.Job, taskrunnerGUI *TaskrunnerGUI) {
			log.Printf("about to run %s (id: %d)\n", job.Name, job.Id)
			jobRun := job.NewJobRun("GUI")

			err := taskrunnerGUI.TaskrunnerDAL.JobRunsDAL.CreateAndRun(jobRun, taskrunnerGUI.JobStatusChangeChan)
			if nil != err {
				log.Printf("ERROR: %s\n", err)
			}
		}(job, jobScene.TaskrunnerGUI)
	}, jobScene.Job)

	configureButton := gtk.NewButton()
	configureButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_EDIT, gtk.ICON_SIZE_LARGE_TOOLBAR))
	configureButton.Clicked(func(ctx *glib.CallbackContext) {
		jobScene, ok := ctx.Data().(*JobScene)
		if !ok {
			panic("couldn't convert to job")
		}

		jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.NewEditJobView(jobScene.Job))
	}, jobScene)

	hbox.PackStart(runButton, false, false, 0)
	hbox.PackEnd(configureButton, false, false, 0)

	box.PackStart(gtk.NewLabel(jobScene.Job.Description), false, false, 5)
	box.PackStart(hbox, false, false, 0)

	box.PackStart(jobScene.buildListing(), true, true, 5)

	return box

}

func (jobScene *JobScene) buildListing() gtk.IWidget {
	var listing gtk.IWidget

	lastJobRun, err := jobScene.TaskrunnerGUI.TaskrunnerDAL.JobRunsDAL.GetLastRun(jobScene.Job)
	if nil != err {
		listing = gtk.NewLabel(err.Error())
	} else if nil == lastJobRun {
		listing = gtk.NewLabel("No runs yet...")
	} else {
		table, err := jobScene.TaskrunnerGUI.buildJobRunsTable(jobScene.Job)
		if nil != err {
			listing = gtk.NewLabel("Error fetching job runs for " + jobScene.Job.Name + ". Error: " + err.Error())
		} else {
			swin := gtk.NewScrolledWindow(nil, nil)
			swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
			swin.SetShadowType(gtk.SHADOW_IN)

			innerVbox := gtk.NewVBox(false, 0)
			innerVbox.PackStart(table, false, false, 0)
			swin.AddWithViewPort(innerVbox)

			listing = swin
		}
	}

	return listing

}

func (taskrunnerGUI *TaskrunnerGUI) buildJobRunsTable(job *taskrunner.Job) (gtk.IWidget, error) {
	runs, err := taskrunnerGUI.TaskrunnerDAL.JobRunsDAL.GetAllForJob(job)
	if nil != err {
		return nil, err
	}
	table := gtk.NewTable(3, uint(len(runs)), false)
	for index, run := range runs {

		runIdButton := gtk.NewButtonWithLabel("#" + strconv.FormatUint(run.Id, 10))
		runIdButton.SetRelief(gtk.RELIEF_NONE)
		runIdButton.Clicked(func(context *glib.CallbackContext) {
			if jobRun, ok := context.Data().(*taskrunner.JobRun); ok {
				taskrunnerGUI.RenderScene(taskrunnerGUI.NewJobRunScene(jobRun))
			} else {
				panic("couldn't cast job")
			}
		}, run)
		startDateTime := time.Unix(run.StartTimestamp, 0)

		table.AttachDefaults(runIdButton, uint(1), 2, uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(startDateTime.String()), 2, 3, uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(run.State.String()), 3, 4, uint(index), uint(index+1))
	}
	swin := gtk.NewViewport(nil, nil)
	swin.Add(table)
	return swin, nil
}
