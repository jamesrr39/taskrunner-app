package gui

import (
	"fmt"
	"log"

	"github.com/jamesrr39/taskrunner-app/taskrunner"

	"strconv"
	"time"

	"github.com/jamesrr39/taskrunner-app/triggers"
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

func (jobScene *JobScene) OnJobRunStatusChange(jobRun *taskrunner.JobRun) {
	if jobRun.Job.Id != jobScene.Job.Id {
		return
	}
	gdk.ThreadsEnter()
	jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.NewJobScene(jobRun.Job))
	gdk.ThreadsLeave()

}

func (jobScene *JobScene) Title() string {
	return "Runs for " + jobScene.Job.Name
}

func (jobScene *JobScene) Content() gtk.IWidget {

	commandsVBox := gtk.NewVBox(false, 0)
	commandsVBox.PackStart(jobScene.buildRunButton(), false, false, 0)
	commandsVBox.PackStart(jobScene.buildConfigureButton(), false, false, 0)

	var descriptionText string
	if "" == jobScene.Job.Description {
		descriptionText = "(No description)"
	} else {
		descriptionText = jobScene.Job.Description
	}

	jobSummaryVBox := gtk.NewVBox(false, 0)
	jobSummaryVBox.PackStart(gtk.NewLabel(descriptionText), false, false, 5)
	jobSummaryVBox.PackStart(jobScene.buildTriggersVBox(), false, false, 0)

	topHbox := gtk.NewHBox(false, 5)
	topHbox.PackStart(commandsVBox, false, false, 30)
	topHbox.PackStart(jobSummaryVBox, true, true, 0)

	box := gtk.NewVBox(false, 5)
	box.PackStart(topHbox, false, false, 0)
	box.PackStart(jobScene.buildListing(), true, true, 0)
	return box

}

func (jobScene *JobScene) buildRunButton() gtk.IWidget {
	runButton := gtk.NewButton()
	runButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_MEDIA_PLAY, gtk.ICON_SIZE_LARGE_TOOLBAR))
	runButton.SetTooltipText("Run Job")
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

	hbox := gtk.NewHBox(false, 5)
	hbox.PackStart(runButton, false, false, 0)
	hbox.PackStart(gtk.NewLabel("Run"), false, false, 0)
	return hbox
}

func (jobScene *JobScene) buildConfigureButton() gtk.IWidget {
	configureButton := gtk.NewButton()
	configureButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_EDIT, gtk.ICON_SIZE_LARGE_TOOLBAR))
	configureButton.SetTooltipText("Configure job")
	configureButton.Clicked(func(ctx *glib.CallbackContext) {
		jobScene, ok := ctx.Data().(*JobScene)
		if !ok {
			panic("couldn't convert to job")
		}

		jobScene.TaskrunnerGUI.RenderScene(jobScene.TaskrunnerGUI.NewEditJobView(jobScene.Job))
	}, jobScene)

	hbox := gtk.NewHBox(false, 5)
	hbox.PackStart(configureButton, false, false, 0)
	hbox.PackStart(gtk.NewLabel("Configure"), false, false, 0)
	return hbox
}

func (jobScene *JobScene) buildTriggersVBox() *gtk.VBox {
	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(jobScene.buildCronJobsSummary(), false, false, 0)
	vbox.PackStart(jobScene.buildUdevJobsSummary(), false, false, 0)
	return vbox
}

func (jobScene *JobScene) buildUdevJobsSummary() gtk.IWidget {
	udevDAL := triggers.NewUdevRulesDAL("/etc/udev/rules.d")
	rules, err := udevDAL.GetRules(jobScene.Job)
	if nil != err {
		return gtk.NewLabel(fmt.Sprintf("Error getting Udev rules: %s", err))
	}

	if 0 == len(rules) {
		return gtk.NewLabel("No Udev rules")
	}

	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(gtk.NewLabel("This rule triggers:"), false, false, 0)
	for _, rule := range rules {
		vbox.PackStart(gtk.NewLabel(fmt.Sprintf("idVendor %s, idProduct %s", rule.IdVendor, rule.IdProduct)), false, false, 0)
	}
	return vbox
}

func (jobScene *JobScene) buildCronJobsSummary() gtk.IWidget {
	cronJobParser := triggers.NewCronParser(fmt.Sprintf("%s %s", jobScene.TaskrunnerGUI.options.CommandPrefix, jobScene.Job.Name))
	cronJobs, err := cronJobParser.SearchForJob(jobScene.Job)
	if nil != err {
		return gtk.NewLabel(fmt.Sprintf("Error getting cron jobs: %s", err))
	}

	if 0 == len(cronJobs) {
		return gtk.NewLabel("No cron jobs")
	}

	vbox := gtk.NewVBox(false, 0)
	for _, cronJob := range cronJobs {
		vbox.PackStart(gtk.NewLabel(cronJob.CronExpression+" "+cronJob.Command), false, false, 0)
	}
	return vbox
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
