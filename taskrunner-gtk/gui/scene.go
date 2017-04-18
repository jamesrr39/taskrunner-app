package gui

import (
	"github.com/jamesrr39/taskrunner-app/taskrunner"

	"github.com/mattn/go-gtk/gtk"
)

type Scene interface {
	Title() string
	Content() gtk.IWidget
	OnJobRunStatusChange(jobRun *taskrunner.JobRun)
}
