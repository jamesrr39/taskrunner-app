package gui

import (
	"taskrunner"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

type EditJobTableEntries struct {
	NameEntry        *gtk.Entry
	DescriptionEntry *gtk.Entry
	StepsEntries     []*gtk.Entry
	ValidationLabel  *gtk.Label
}

func NewEditJobTable(numberOfSteps uint) *EditJobTableEntries {
	editJobTable := &EditJobTableEntries{
		NameEntry:        gtk.NewEntry(),
		DescriptionEntry: gtk.NewEntry(),
		ValidationLabel:  gtk.NewLabel(""),
	}

	var stepsEntries []*gtk.Entry
	for i := uint(0); i < numberOfSteps; i++ {
		stepsEntries = append(stepsEntries, gtk.NewEntry())
	}
	editJobTable.StepsEntries = stepsEntries
	return editJobTable
}

func (taskrunnerGUI *TaskrunnerGUI) buildNewJobView() gtk.IWidget {
	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(gtk.NewLabel("New Job Setup"), false, false, 0)

	editJobTableEntries := NewEditJobTable(1)

	table := gtk.NewTable(3, 2, false)
	table.AttachDefaults(gtk.NewLabel("Name"), 0, 1, 0, 1)
	table.AttachDefaults(editJobTableEntries.NameEntry, 1, 2, 0, 1)
	table.AttachDefaults(gtk.NewLabel("Description"), 0, 1, 1, 2)
	table.AttachDefaults(editJobTableEntries.DescriptionEntry, 1, 2, 1, 2)
	table.AttachDefaults(gtk.NewLabel("Steps"), 0, 1, 2, 3)
	table.AttachDefaults(editJobTableEntries.StepsEntries[0], 1, 2, 2, 3)
	vbox.PackStart(table, false, false, 0)

	createButton := gtk.NewButtonWithLabel("Create!")
	createButton.Clicked(func(ctx *glib.CallbackContext) {
		entries, ok := ctx.Data().(*EditJobTableEntries)
		if !ok {
			panic("not ok")
		}
		var steps []*taskrunner.Step

		job, err := taskrunner.NewJob(
			entries.NameEntry.GetText(),
			entries.DescriptionEntry.GetText(),
			steps,
			taskrunnerGUI.TaskrunnerInstance)

		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}

		err = job.Save()
		if nil != err {
			entries.ValidationLabel.SetLabel(err.Error())
			entries.ValidationLabel.ShowAll()
			return
		}
		taskrunnerGUI.RenderHomeScreen()

	}, editJobTableEntries)

	vbox.PackEnd(createButton, false, false, 0)
	editJobTableEntries.ValidationLabel.ModifyFG(gtk.STATE_NORMAL, gdk.NewColorRGB(255, 0, 0))
	vbox.PackEnd(editJobTableEntries.ValidationLabel, false, false, 0)
	return vbox

}
