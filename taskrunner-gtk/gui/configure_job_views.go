package gui

import (
	"taskrunner-app/taskrunner"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type EditJobView struct {
	*TaskrunnerGUI
	Job      *taskrunner.Job
	isClosed bool
}

func (taskrunnerGUI *TaskrunnerGUI) NewEditJobView(job *taskrunner.Job) *EditJobView {
	return &EditJobView{taskrunnerGUI, job, true}
}

func (editJobView *EditJobView) Title() string {
	return "Editing :: " + editJobView.Job.Name
}

func (editJobView *EditJobView) OnClose() {
	editJobView.isClosed = true
}

func (editJobView *EditJobView) OnShow() {
	editJobView.isClosed = false
}

func (editJobView *EditJobView) Content() gtk.IWidget {
	vbox := gtk.NewVBox(false, 0)

	// editing table
	editJobTableEntries := editJobView.TaskrunnerGUI.NewConfigureJobTableEntries(editJobView.Job)

	vbox.PackStart(editJobTableEntries.ValidationLabel.Widget, false, true, 0)

	vbox.PackStart(editJobTableEntries.ToTable(), false, false, 0)

	saveButton := gtk.NewButtonWithLabel("Save!")
	saveButton.Clicked(func(ctx *glib.CallbackContext) {
		entries, ok := ctx.Data().(*ConfigureJobTableEntries)
		if !ok {
			panic("couldn't convert createJobTableEntries")
		}

		// start test validation label
		entries.ValidationLabel.SetText("error - red")
		return
		// end test validation label

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
