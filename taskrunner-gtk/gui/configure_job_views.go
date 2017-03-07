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
	vbox.PackStart(editJobTableEntries.ToTable(), false, false, 0)

	createButton := gtk.NewButtonWithLabel("Save!")
	createButton.Clicked(func(ctx *glib.CallbackContext) {
		entries, ok := ctx.Data().(*ConfigureJobTableEntries)
		if !ok {
			panic("couldn't convert createJobTableEntries")
		}

		job, err := entries.ToJob(0)
		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}

		err = editJobView.TaskrunnerDAL.JobDAL.Create(job)
		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}
		entries.TaskrunnerGUI.RenderScene(entries.TaskrunnerGUI.NewJobScene(job))

	}, editJobTableEntries)
	vbox.PackEnd(createButton, false, false, 0)

	return vbox
}
