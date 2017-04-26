package gui

import (
	"github.com/jamesrr39/taskrunner-app/taskrunner"

	"github.com/mattn/go-gtk/gtk"
)

func buildToolbar(taskrunnerGUI *TaskrunnerGUI) gtk.IWidget {
	homeButton := gtk.NewButton()
	homeButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_HOME, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	homeButton.SetTooltipText("Back to home")
	homeButton.Clicked(func() {
		taskrunnerGUI.RenderScene(taskrunnerGUI.NewHomeScene())
	})

	newJobButton := gtk.NewButton()
	newJobButton.SetImage(gtk.NewImageFromStock(gtk.STOCK_ADD, gtk.IconSize(gtk.ICON_SIZE_LARGE_TOOLBAR)))
	newJobButton.SetTooltipText("Create a new job")
	newJobButton.Clicked(func() {
		newJob, err := taskrunner.NewJob(0, "New Job", "", taskrunner.Script("#!/bin/bash\n\n"))
		if nil != err {
			panic(err)
		}
		taskrunnerGUI.RenderScene(taskrunnerGUI.NewEditJobView(newJob, taskrunnerGUI.udevRulesDAL))
	})

	hbox := gtk.NewHBox(false, 0)
	hbox.PackStart(taskrunnerGUI.titleLabel, true, true, 3)
	hbox.PackEnd(homeButton, false, false, 0)
	hbox.PackEnd(newJobButton, false, false, 0)

	eventBox := gtk.NewEventBox()
	eventBox.Add(hbox)
	eventBox.ModifyBG(gtk.STATE_NORMAL, titleBlue())

	return eventBox
}
