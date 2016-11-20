package gui

import (
	"taskrunner"

	//	"github.com/mattn/go-gtk/gdk"
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
	vbox.PackStart(gtk.NewLabel("Editing :: "+editJobView.Job.Name), false, false, 0)

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

/*
func (taskrunnerGUI *TaskrunnerGUI) makeConfigureCreateJobView() gtk.IWidget {
	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(gtk.NewLabel("New Job Setup"), false, false, 0)

	job := &taskrunner.Job{TaskrunnerInstance: taskrunnerGUI.TaskrunnerInstance}
	log.Printf("new job: %v\n", job)

	// create job entries table
	createJobTableEntries := taskrunnerGUI.NewConfigureJobTableEntries(job)
	vbox.PackStart(createJobTableEntries.ToTable(), false, false, 0)

	createButton := gtk.NewButtonWithLabel("Create!")
	createButton.Clicked(func(ctx *glib.CallbackContext) {
		entries, ok := ctx.Data().(*ConfigureJobTableEntries)
		if !ok {
			panic("couldn't convert EditJobTableEntries")
		}

		job, err := entries.ToJob()

		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}

		err = job.TaskrunnerInstance.CreateJob(job)
		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}
		entries.TaskrunnerGUI.RenderScene(entries.TaskrunnerGUI.NewJobScene(job))

	}, createJobTableEntries)

	vbox.PackEnd(createButton, false, false, 0)
	createJobTableEntries.ValidationLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColorRGB(255, 0, 0))
	vbox.PackEnd(createJobTableEntries.ValidationLabel, false, false, 0)
	return vbox

}
*/
