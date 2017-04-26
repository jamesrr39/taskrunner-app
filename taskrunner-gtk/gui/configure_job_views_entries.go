package gui

import (
	"log"

	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/jamesrr39/taskrunner-app/triggers"

	"github.com/mattn/go-gtk/gtk"
)

type ConfigureJobTableEntries struct {
	NameEntry               *gtk.Entry
	DescriptionEntry        *gtk.Entry
	ScriptEntry             *gtk.TextView
	scriptEntryScrollWindow *gtk.ScrolledWindow
	*ValidationLabel
	*TaskrunnerGUI
	*taskrunner.Job
	udevRulesDAL *triggers.UdevRulesDAL
}

func (taskrunnerGUI *TaskrunnerGUI) NewConfigureJobTableEntries(job *taskrunner.Job, udevRulesDAL *triggers.UdevRulesDAL) *ConfigureJobTableEntries {
	editJobTable := &ConfigureJobTableEntries{
		NameEntry:        gtk.NewEntry(),
		DescriptionEntry: gtk.NewEntry(),
		ValidationLabel:  NewValidationLabel(errorRed()),
		TaskrunnerGUI:    taskrunnerGUI,
		ScriptEntry:      gtk.NewTextView(),
		Job:              job,
		udevRulesDAL:     udevRulesDAL,
	}

	editJobTable.NameEntry.SetText(job.Name)
	editJobTable.DescriptionEntry.SetText(job.Description)

	editJobTable.ScriptEntry.SetBorderWidth(2)
	editJobTable.ScriptEntry.GetBuffer().SetText(string(editJobTable.Job.Script))

	editJobTable.scriptEntryScrollWindow = editJobTable.GetScriptScrollWindow()

	editJobTable.scriptEntryScrollWindow.AddWithViewPort(editJobTable.ScriptEntry)

	return editJobTable
}

func (editJobTableEntries *ConfigureJobTableEntries) ToWidget() gtk.IWidget {
	table := gtk.NewTable(3, 2, false)
	table.AttachDefaults(gtk.NewLabel("Name"), 0, 1, 0, 1)
	table.AttachDefaults(editJobTableEntries.NameEntry, 1, 2, 0, 1)
	table.AttachDefaults(gtk.NewLabel("Description"), 0, 1, 1, 2)
	table.AttachDefaults(editJobTableEntries.DescriptionEntry, 1, 2, 1, 2)

	udevTable := editJobTableEntries.buildUdevDisplay()

	vbox := gtk.NewVBox(false, 0)
	vbox.PackStart(table, false, false, 0)
	vbox.PackStart(udevTable, false, false, 0)

	scriptLabel := gtk.NewLabel("Script:")
	scriptLabel.SetAlignment(0, 0)
	vbox.PackStart(scriptLabel, false, false, 0)

	vbox.PackEnd(editJobTableEntries.scriptEntryScrollWindow, true, true, 0)
	return vbox
}

func (editJobTableEntries *ConfigureJobTableEntries) ToJob(jobId uint) (*taskrunner.Job, error) {

	job, err := taskrunner.NewJob(
		jobId,
		editJobTableEntries.NameEntry.GetText(),
		editJobTableEntries.DescriptionEntry.GetText(),
		editJobTableEntries.GetScriptEntryText())

	if nil != err {
		return nil, err
	}

	if nil != editJobTableEntries.Job {
		job.Id = editJobTableEntries.Job.Id
	}

	return job, nil
}

func (editJobTableEntries *ConfigureJobTableEntries) GetScriptScrollWindow() *gtk.ScrolledWindow {
	logTextareaScrollWindow := gtk.NewScrolledWindow(nil, nil)
	logTextareaScrollWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	logTextareaScrollWindow.SetShadowType(gtk.SHADOW_IN)

	return logTextareaScrollWindow
}

func (editJobTableEntries *ConfigureJobTableEntries) GetScriptEntryText() taskrunner.Script {
	var startIter, endIter gtk.TextIter

	scriptEntryBuffer := editJobTableEntries.ScriptEntry.GetBuffer()
	scriptEntryBuffer.GetStartIter(&startIter)
	scriptEntryBuffer.GetEndIter(&endIter)
	return taskrunner.Script(scriptEntryBuffer.GetText(&startIter, &endIter, false))
}

func (editJobTableEntries *ConfigureJobTableEntries) buildUdevDisplay() gtk.IWidget {
	rules, err := editJobTableEntries.udevRulesDAL.GetRules(editJobTableEntries.Job)
	if nil != err {
		return gtk.NewLabel("Failed to get udev rules. Error: " + err.Error())
	}

	vbox := gtk.NewVBox(false, 0)

	if 0 == len(rules) {
		vbox.PackStart(gtk.NewLabel("No rules found for this job"), false, false, 0)
	} else {
		vbox.PackStart(editJobTableEntries.buildUdevInnerTable(rules), false, false, 0)
	}

	vbox.PackStart(gtk.NewLabel("Adding rules is not currently supported through the user interface, however they can be added manually."), false, false, 0)

	return vbox
}

func (editJobTableEntries *ConfigureJobTableEntries) buildUdevInnerTable(rules []*triggers.UdevRule) gtk.IWidget {
	table := gtk.NewTable(uint(len(rules)+1), 3, true)
	table.AttachDefaults(gtk.NewLabel("idVendor"), 0, 1, 0, 1)
	table.AttachDefaults(gtk.NewLabel("idProduct"), 1, 2, 0, 1)
	for index, rule := range rules {
		deleteButton := gtk.NewButtonFromStock(gtk.STOCK_DELETE)
		ruleIndex := index
		deleteButton.Clicked(func() {
			log.Printf("TODO delete rule #%d\n", ruleIndex)
		})

		table.AttachDefaults(gtk.NewLabel(rule.IdVendor), 0, 1, uint(index+1), uint(index+2))
		table.AttachDefaults(gtk.NewLabel(rule.IdProduct), 1, 2, uint(index+1), uint(index+2))
		table.AttachDefaults(deleteButton, 2, 3, uint(index+1), uint(index+2))
	}

	return table
}
