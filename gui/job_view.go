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

type JobScene struct {
	*TaskrunnerGUI
	Job *taskrunner.Job
}

func (taskrunnerGUI *TaskrunnerGUI) NewJobScene(job *taskrunner.Job) *JobScene {
	return &JobScene{taskrunnerGUI, job}
}

func (jobScene *JobScene) IsCurrentlyRendered() bool {
	paneContentJobScene, ok := jobScene.TaskrunnerGUI.PaneContent.(*JobScene)
	if ok && paneContentJobScene.Job.Id == jobScene.Job.Id {
		return true
	}
	return false
}

func (jobScene *JobScene) Content() gtk.IWidget {

	box := gtk.NewVBox(false, 5)

	box.PackStart(gtk.NewLabel("Runs for "+jobScene.Job.Name), false, false, 5)

	hbox := gtk.NewHBox(true, 0)
	runButton := gtk.NewButton()
	runButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_LARGE_TOOLBAR))
	runButton.Clicked(func(ctx *glib.CallbackContext) {
		job, ok := ctx.Data().(*taskrunner.Job)
		if !ok {
			panic("couldn't convert to job")
		}
		go func(job *taskrunner.Job, taskrunnerGUI *TaskrunnerGUI) {
			jobRun, _ := job.Run("GUI")
			taskrunnerGUI.JobStatusChangeChan <- jobRun
			/*gdk.ThreadsEnter()
			taskrunnerGUI.RenderJobRuns(job) // todo check still on this screen interface CurrentSceneRendered
			gdk.ThreadsLeave()*/
		}(job, jobScene.TaskrunnerGUI)

		// refresh
		jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.NewJobScene(job))
	}, jobScene.Job)

	configureButton := gtk.NewButton()
	configureButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_EDIT, gtk.ICON_SIZE_LARGE_TOOLBAR))
	configureButton.Clicked(func(ctx *glib.CallbackContext) {
		job, ok := ctx.Data().(*taskrunner.Job)
		if !ok {
			panic("couldn't convert to job")
		}
		log.Printf("job for configure edit job view: %v\n", job)

		//jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.makeConfigureEditJobView(job))
	}, jobScene.Job)

	hbox.PackStart(runButton, false, false, 0)
	hbox.PackEnd(configureButton, false, false, 0)
	box.PackStart(hbox, false, false, 0)
	var listing gtk.IWidget

	if 0 == jobScene.Job.GetLastRunId() {
		listing = gtk.NewLabel("No runs yet...")
	} else {
		table, err := jobScene.TaskrunnerGUI.buildJobRunsTable(jobScene.Job)
		if nil != err {
			listing = gtk.NewLabel("Error fetching job runs for " + jobScene.Job.Name + ". Error: " + err.Error())
		} else {
			listing = table
		}
	}
	box.PackStart(listing, false, false, 5)

	go func(renderedJob *taskrunner.Job) {
		jobRun := <-jobScene.TaskrunnerGUI.JobStatusChangeChan
		if renderedJob.Id == jobRun.Job.Id {
			log.Println("listing in channel")
			gdk.ThreadsEnter()
			jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.NewJobScene(renderedJob)) // todo check still on this screen interface CurrentSceneRendered
			gdk.ThreadsLeave()
		} else {
			log.Println("SKIPPING job re-render ------------")
		}
	}(jobScene.Job)

	return box

}

func (taskrunnerGUI *TaskrunnerGUI) buildJobRunsTable(job *taskrunner.Job) (gtk.IWidget, error) {
	runs, err := job.GetRuns()
	if nil != err {
		return nil, err
	}
	table := gtk.NewTable(3, uint(len(runs)), false)
	for index, run := range runs {

		runIdButton := gtk.NewButtonWithLabel("#" + strconv.Itoa(run.Id))
		runIdButton.SetRelief(gtk.RELIEF_NONE)
		runIdButton.Clicked(func(context *glib.CallbackContext) {
			if _, ok := context.Data().(*taskrunner.JobRun); ok {
				taskrunnerGUI.RenderScene(taskrunnerGUI.NewJobScene(job))
			} else {
				panic("couldn't cast job")
			}
		}, run)
		startDateTime := time.Unix(run.StartTimestamp, 0)

		table.AttachDefaults(runIdButton, uint(1), 2, uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(startDateTime.String()), 2, 3, uint(index), uint(index+1))
		table.AttachDefaults(gtk.NewLabel(run.State.String()), 3, 4, uint(index), uint(index+1))
	}
	return table, nil
}
