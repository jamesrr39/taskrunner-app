package main

import (
	"fmt"
	"log"
	"os"
	"taskrunner"
	"time"

	"github.com/alecthomas/kingpin"
)

var (
	taskrunnerInstance    *taskrunner.TaskrunnerInstance
	taskrunnerApplication *kingpin.Application
)

func main() {
	taskrunnerApplication = kingpin.New("Taskrunner", "Taskrunner CLI application")

	setupApplicationFlags()
	setupListJobsCommand()
	setupRunJobCommand()
	setupPrintJobConfigCommand()
	setupPrintJobHistoryCommand()

	kingpin.MustParse(taskrunnerApplication.Parse(os.Args[1:]))

}

func setupApplicationFlags() {
	taskrunnerDir := taskrunnerApplication.
		Flag("taskrunner-dir", "Directory the taskruner uses to store job configs and logs of job runs.").
		Default("~/.taskrunner").
		String()

	taskrunnerApplication.Action(func(context *kingpin.ParseContext) error {
		var err error
		taskrunnerInstance, err = taskrunner.NewTaskrunnerInstance(*taskrunnerDir)
		return err
	})
}

func setupListJobsCommand() {
	listJobCommand := taskrunnerApplication.Command("list-jobs", "List all jobs")
	listJobCommand.Action(func(context *kingpin.ParseContext) error {
		jobs, err := taskrunnerInstance.Jobs()
		if nil != err {
			return fmt.Errorf("Error listing all jobs. Error: %s", err)
		}

		fmt.Println("NAME | DESCRIPTION | LAST RUN")
		for _, job := range jobs {
			lastRunSummary := getLastRunSummary(job)
			fmt.Printf("%s | %s | %s\n", job.Name, job.Description, lastRunSummary)
		}
		return nil
	})
}

func setupPrintJobHistoryCommand() {
	printJobHistoryCommand := taskrunnerApplication.Command("print-job-history", "Print out the history for a job")
	jobName := addJobNameToCommandArgs(printJobHistoryCommand)
	printJobHistoryCommand.Action(func(context *kingpin.ParseContext) error {
		err := printJobHistory(*jobName)
		return err
	})
}

func setupRunJobCommand() {
	runJobCommand := taskrunnerApplication.Command("run-job", "Run a job")
	jobName := addJobNameToCommandArgs(runJobCommand)
	runJobCommand.Action(func(context *kingpin.ParseContext) error {
		err := runJob(*jobName)
		return err
	})
}

func setupPrintJobConfigCommand() {
	printJobConfigCommand := taskrunnerApplication.Command("print-job-config", "Print out a job config")
	jobName := addJobNameToCommandArgs(printJobConfigCommand)
	//jobName = printJobConfigCommand.Arg("Job_Name", "Name of the job to be run").Required().String()
	printJobConfigCommand.Action(func(context *kingpin.ParseContext) error {
		err := printJobConfig(*jobName)
		return err
	})
}

func addJobNameToCommandArgs(cmdClause *kingpin.CmdClause) *string {
	return cmdClause.Arg("Job_Name", "Name of the job to be run").Required().String()
}

func runJob(jobName string) error {
	job, err := taskrunnerInstance.JobFromName(jobName)
	if nil != err {
		return err
	}

	err = job.Run("CLI")
	return err
}

func printJobConfig(jobName string) error {
	job, err := taskrunnerInstance.JobFromName(jobName)
	if nil != err {
		return fmt.Errorf("Error getting job '%s'. Error: %s\n", jobName, err)
	}
	fmt.Printf("Name: %s\nDescription: %s\n", job.Name, job.Description)
	if 0 == len(job.Steps) {
		fmt.Println("No steps in this job")
		return nil
	}

	fmt.Println("Steps:")
	for index, step := range job.Steps {
		fmt.Printf("Step %d, %s:\n> %s\n", index+1, step.Name, step.Cmd)
	}
	return nil
}

func printJobHistory(jobName string) error {
	job, err := taskrunnerInstance.JobFromName(jobName)
	if nil != err {
		return fmt.Errorf("Error getting job '%s'. Error: %s\n", jobName, err)
	}

	jobRuns, err := job.Runs()
	if nil != err {
		fmt.Printf("Error getting job runs for job '%s'. Error: %s\n", jobName, err)
		return nil
	}

	for _, jobRun := range jobRuns {
		lastRunTime := time.Unix(jobRun.StartTimestamp, 0)
		fmt.Printf("#%d: %s: %t\n", jobRun.Id, lastRunTime.Format(time.RFC1123), jobRun.Successful)
	}
	return nil

}

func getLastRunSummary(job *taskrunner.Job) string {
	lastRunId := job.GetLastRunId()
	if 0 == lastRunId {
		return "No runs of this job"
	}

	jobRun, err := job.GetRun(lastRunId)

	if nil != err {
		log.Printf("Error decoding last job run for job %s. Error: %s\n", job.Name, err)
		return fmt.Sprintf("Error decoding last job run for job %s.", job.Name)
	}

	lastRunTime := time.Unix(jobRun.StartTimestamp, 0)

	return fmt.Sprintf("#%d: %s: %t", lastRunId, lastRunTime.Format(time.RFC1123), jobRun.Successful)

}
