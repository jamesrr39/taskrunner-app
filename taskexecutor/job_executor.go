package taskexecutor

import (
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"time"
)

const TASKRUNNER_SOURCE_NAME string = "TASKRUNNER"

func ExecuteJobRun(jobRun *taskrunner.JobRun, jobRunStatusChangeChan chan *taskrunner.JobRun, logFile io.Writer, workspaceDir string) {

	scriptFilePath := filepath.Join(workspaceDir, "script")
	err := ioutil.WriteFile(scriptFilePath, []byte(jobRun.Job.Script), 0500)
	if nil != err {
		handleTaskrunnerError("Couldn't prepare and move to workspace. Error: "+err.Error(), logFile, jobRunStatusChangeChan, jobRun)
		return
	}

	cmd := exec.Command(scriptFilePath)
	stdoutPipe, err := cmd.StdoutPipe()
	if nil != err {
		handleTaskrunnerError("Couldn't obtain stdoutpipe. Error: "+err.Error(), logFile, jobRunStatusChangeChan, jobRun)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if nil != err {
		handleTaskrunnerError("Couldn't obtain stderrpipe. Error: "+err.Error(), logFile, jobRunStatusChangeChan, jobRun)
		return
	}

	go writeToLogFile(stdoutPipe, logFile, "STDOUT")
	go writeToLogFile(stderrPipe, logFile, "STDERR")

	err = cmd.Start()
	if nil != err {
		handleTaskrunnerError("Couldn't start script. Error: "+err.Error(), logFile, jobRunStatusChangeChan, jobRun)
		return
	}

	err = cmd.Wait()
	if nil != err {
		switch err.(type) {
		case *exec.ExitError:
			jobRun.State = taskrunner.JOB_RUN_STATE_FAILED
		default:
			jobRun.State = taskrunner.JOB_RUN_STATE_UNKNOWN
		}
	} else {
		jobRun.State = taskrunner.JOB_RUN_STATE_SUCCESS
	}

	jobRun.EndTimestamp = time.Now().Unix()
	jobRunStatusChangeChan <- jobRun

	return
}

func handleTaskrunnerError(errorMessage string, logFile io.Writer, jobRunStateChan chan *taskrunner.JobRun, jobRun *taskrunner.JobRun) {
	jobRun.EndTimestamp = time.Now().Unix()
	jobRun.State = taskrunner.JOB_RUN_STATE_FAILED
	jobRunStateChan <- jobRun
	writeStringToLogFile(errorMessage, logFile, TASKRUNNER_SOURCE_NAME)
}
