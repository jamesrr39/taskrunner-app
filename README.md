# Taskrunner-app

Taskrunner is a CLI & GUI application, with the intention of sitting somewhere between cron jobs and Jenkins.

The idea is to have a tool that looks after your repetitive tasks, and stores the output somewhere you can find the output easily.

## Install

First, you need to have Go and GTK 2 installed.

Installing GTK on ubuntu:

    sudo make install_dependencies

Project can be built as a normal Go project using the standard Go tool.

## Build & Run

    make run

### Open GUI summary of jobs

    ./taskrunner-app

### Run a job headlessly

It's also possible to run a job headlessly, and have the output saved so that it can be viewed in the UI next time it's opened. This is ideal for use by scripts and cron jobs.

The `--trigger` may be viewed as a kind of "user-agent", and is something that is surfaced later in the GUI.

    ./taskrunner-app --run-job="failing job" --trigger="my custom trigger
	
### More options

```
./taskrunner-app -h
```
