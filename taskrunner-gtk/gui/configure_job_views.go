package gui

import (
	"github.com/jamesrr39/taskrunner-app/taskrunner"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type EditJobView struct {
	*TaskrunnerGUI
	Job *taskrunner.Job
}

func (taskrunnerGUI *TaskrunnerGUI) NewEditJobView(job *taskrunner.Job) *EditJobView {
	return &EditJobView{taskrunnerGUI, job}
}

func (editJobView *EditJobView) OnJobRunStatusChange(jobRun *taskrunner.JobRun) {
}

func (editJobView *EditJobView) Title() string {
	return "Editing :: " + editJobView.Job.Name
}

func (editJobView *EditJobView) Content() gtk.IWidget {
	vbox := gtk.NewVBox(false, 0)

	// editing table
	editJobTableEntries := editJobView.TaskrunnerGUI.NewConfigureJobTableEntries(editJobView.Job)

	vbox.PackStart(editJobTableEntries.ValidationLabel.Widget, false, true, 0)

	vbox.PackStart(editJobTableEntries.ToWidget(), true, true, 0)

	saveButton := gtk.NewButtonWithLabel("Save!")
	saveButton.Clicked(func(ctx *glib.CallbackContext) {
		entries, ok := ctx.Data().(*ConfigureJobTableEntries)
		if !ok {
			panic("couldn't convert createJobTableEntries")
		}

		job, err := entries.ToJob(editJobView.Job.Id)
		if nil != err {
			entries.ValidationLabel.SetText(err.Error())
			return
		}

		if 0 == job.Id {
			err = editJobView.TaskrunnerDAL.JobDAL.Create(job)
		} else {
			err = editJobView.TaskrunnerDAL.JobDAL.Update(job)
		}
		if nil != err {
			entries.ValidationLabel.SetText(err.Error())
			return
		}
		entries.TaskrunnerGUI.RenderScene(entries.TaskrunnerGUI.NewJobScene(job))

	}, editJobTableEntries)
	vbox.PackEnd(saveButton, false, false, 0)

	return vbox
}
