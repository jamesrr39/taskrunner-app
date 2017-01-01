package gui

import (
	"taskrunner"

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

func (editJobView *EditJobView) Title() string {
	return "Editing :: " + editJobView.Job.Name
}

func (editJobView *EditJobView) IsCurrentlyRendered() bool {
	switch editJobView.TaskrunnerGUI.PaneContent.(type) {
	case *EditJobView:
		return true
	default:
		return false
	}
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

		job, err := entries.ToJob()
		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}

		err = job.TaskrunnerInstance.SaveJob(job)
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
