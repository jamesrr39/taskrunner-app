package gui

import (
	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/jamesrr39/taskrunner-app/triggers"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type EditJobView struct {
	*TaskrunnerGUI
	Job          *taskrunner.Job
	udevRulesDAL *triggers.UdevRulesDAL
}

func (taskrunnerGUI *TaskrunnerGUI) NewEditJobView(job *taskrunner.Job, udevRulesDAL *triggers.UdevRulesDAL) *EditJobView {
	return &EditJobView{taskrunnerGUI, job, udevRulesDAL}
}

func (editJobView *EditJobView) OnJobRunStatusChange(jobRun *taskrunner.JobRun) {
}

func (editJobView *EditJobView) Title() string {
	return "Editing :: " + editJobView.Job.Name
}

func (editJobView *EditJobView) Content() gtk.IWidget {

	topHbox := gtk.NewHBox(false, 0)
	topHbox.PackStart(editJobView.buildGoUpButton(), false, false, 30)

	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(topHbox, false, false, 0)
	// editing table
	editJobTableEntries := editJobView.TaskrunnerGUI.NewConfigureJobTableEntries(editJobView.Job, editJobView.udevRulesDAL)

	vbox.PackStart(editJobTableEntries.ValidationLabel.Widget, false, true, 0)
	vbox.PackStart(editJobTableEntries.ToWidget(), true, true, 0)

	saveButton := gtk.NewButtonWithLabel("Save!")
	saveButton.Clicked(editJobView.onSaveClicked, editJobTableEntries)
	vbox.PackEnd(saveButton, false, false, 0)

	return vbox
}

func (editJobView *EditJobView) onSaveClicked(ctx *glib.CallbackContext) {
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
}

func (editJobView *EditJobView) buildGoUpButton() gtk.IWidget {
	text := "Back to Job Overview (discard changes)"

	goUpButton := gtk.NewButton()
	goUpButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_GO_UP, gtk.ICON_SIZE_LARGE_TOOLBAR))
	goUpButton.SetTooltipText(text)
	goUpButton.Clicked(func() {
		editJobView.TaskrunnerGUI.RenderScene(editJobView.TaskrunnerGUI.NewJobScene(editJobView.Job))
	})

	goUpHbox := gtk.NewHBox(false, 5)
	goUpHbox.PackStart(goUpButton, false, false, 0)
	goUpHbox.PackStart(gtk.NewLabel(text), false, false, 0)

	return goUpHbox
}
