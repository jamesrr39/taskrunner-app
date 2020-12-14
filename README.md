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

    ./taskrunner-app --run-job="failing job" --trigger="my custom trigger
	

`./taskrunner-app -h` gives more options.
