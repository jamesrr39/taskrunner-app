package gui

import (
	"taskrunner-app/taskrunner"

	"github.com/mattn/go-gtk/gtk"
)

type ConfigureJobTableEntries struct {
	NameEntry        *gtk.Entry
	DescriptionEntry *gtk.Entry
	ScriptEntry      *gtk.TextView
	ValidationLabel  *gtk.Label
	*TaskrunnerGUI
	*taskrunner.Job
}

func (taskrunnerGUI *TaskrunnerGUI) NewConfigureJobTableEntries(job *taskrunner.Job) *ConfigureJobTableEntries {
	editJobTable := &ConfigureJobTableEntries{
		NameEntry:        gtk.NewEntry(),
		DescriptionEntry: gtk.NewEntry(),
		ValidationLabel:  gtk.NewLabel(""),
		TaskrunnerGUI:    taskrunnerGUI,
		Job:              job,
	}

	editJobTable.ScriptEntry = gtk.NewTextView()
	editJobTable.ScriptEntry.SetBorderWidth(2)

	editJobTable.NameEntry.SetText(job.Name)
	editJobTable.DescriptionEntry.SetText(job.Description)
	editJobTable.ScriptEntry.GetBuffer().SetText(string(job.Script))

	return editJobTable
}

func (editJobTableEntries *ConfigureJobTableEntries) ToTable() *gtk.Table {

	table := gtk.NewTable(3, 2, false)
	table.AttachDefaults(gtk.NewLabel("Name"), 0, 1, 0, 1)
	table.AttachDefaults(editJobTableEntries.NameEntry, 1, 2, 0, 1)
	table.AttachDefaults(gtk.NewLabel("Description"), 0, 1, 1, 2)
	table.AttachDefaults(editJobTableEntries.DescriptionEntry, 1, 2, 1, 2)
	table.AttachDefaults(gtk.NewLabel("Script"), 0, 1, 2, 3)
	table.AttachDefaults(editJobTableEntries.ScriptEntry, 1, 2, 2, 3)
	return table
}

func (editJobTableEntries *ConfigureJobTableEntries) ToJob() (*taskrunner.Job, error) {

	job, err := taskrunner.NewJob(
		editJobTableEntries.NameEntry.GetText(),
		editJobTableEntries.DescriptionEntry.GetText(),
		editJobTableEntries.GetScriptEntryText(),
		editJobTableEntries.TaskrunnerGUI.TaskrunnerInstance)

	if nil != err {
		return nil, err
	}

	if nil != editJobTableEntries.Job {
		job.Id = editJobTableEntries.Job.Id
	}

	return job, nil
}

func (editJobTableEntries *ConfigureJobTableEntries) GetScriptEntryText() taskrunner.Script {
	var startIter, endIter gtk.TextIter

	scriptEntryBuffer := editJobTableEntries.ScriptEntry.GetBuffer()
	scriptEntryBuffer.GetStartIter(&startIter)
	scriptEntryBuffer.GetEndIter(&endIter)
	return taskrunner.Script(scriptEntryBuffer.GetText(&startIter, &endIter, false))
}
