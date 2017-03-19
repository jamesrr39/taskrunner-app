package gui

import (
	"taskrunner-app/taskrunner"

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
}

func (taskrunnerGUI *TaskrunnerGUI) NewConfigureJobTableEntries(job *taskrunner.Job) *ConfigureJobTableEntries {
	editJobTable := &ConfigureJobTableEntries{
		NameEntry:        gtk.NewEntry(),
		DescriptionEntry: gtk.NewEntry(),
		ValidationLabel:  NewValidationLabel(errorRed()),
		TaskrunnerGUI:    taskrunnerGUI,
		ScriptEntry:      gtk.NewTextView(),
		Job:              job,
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
	vbox := gtk.NewVBox(false, 0)
	table := gtk.NewTable(3, 2, false)
	table.AttachDefaults(gtk.NewLabel("Name"), 0, 1, 0, 1)
	table.AttachDefaults(editJobTableEntries.NameEntry, 1, 2, 0, 1)
	table.AttachDefaults(gtk.NewLabel("Description"), 0, 1, 1, 2)
	table.AttachDefaults(editJobTableEntries.DescriptionEntry, 1, 2, 1, 2)
	vbox.PackStart(table, false, false, 0)

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
