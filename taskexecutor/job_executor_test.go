package taskexecutor

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamesrr39/taskrunner-app/taskrunner"

	"github.com/stretchr/testify/assert"
)

func Test_ExecuteJobRun(t *testing.T) {
	workspaceDir, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer func() {
		err := os.RemoveAll(workspaceDir)
		if nil != err {
			t.Errorf("Couldn't remove the tempdir at '%s'. Error: %s.\n", workspaceDir, err)
		} else {
			t.Logf("Successfully removed workspace dir at %s\n", workspaceDir)
		}
	}()

	logFile := bytes.NewBuffer(nil)

	jobRunStateChan := make(chan *taskrunner.JobRun)

	job, err := taskrunner.NewJob(0, "my job", "the big job", "#!/bin/bash\n\necho 'job failed'\nexit 1")
	assert.Nil(t, err)
	jobRun := job.NewJobRun("test trigger")

	jobRunStateGoRoutineDoneChan := make(chan bool)

	var newJobRunState *taskrunner.JobRun
	go func() {
		newJobRunState = <-jobRunStateChan
		jobRunStateGoRoutineDoneChan <- true
	}()

	err = ExecuteJobRun(jobRun, jobRunStateChan, logFile, workspaceDir, mockNowProvider)
	assert.Nil(t, err)

	<-jobRunStateGoRoutineDoneChan
	assert.Equal(t, "03:04:05.006: STDOUT: job failed\n", string(logFile.Bytes()))
	assert.Equal(t, taskrunner.JOB_RUN_STATE_FAILED, newJobRunState.State)

}

func Test_handleTaskrunnerError(t *testing.T) {
	errorMessage := "setup failed"
	logFile := bytes.NewBuffer(nil)
	jobRunStateChan := make(chan *taskrunner.JobRun)

	job, err := taskrunner.NewJob(0, "my job", "the big job", "#!/bin/bash\n\necho 'my job'")
	assert.Nil(t, err)
	jobRun := job.NewJobRun("test trigger")

	var newJobRunState *taskrunner.JobRun
	go func() {
		// listen to jobRunStateChan
		newJobRunState = <-jobRunStateChan
	}()

	err = handleTaskrunnerError(errorMessage, logFile, jobRunStateChan, jobRun, mockNowProvider)
	assert.Nil(t, err)

	assert.Equal(t, taskrunner.JOB_RUN_STATE_FAILED, newJobRunState.State)
	assert.Equal(t, int64(946782245), newJobRunState.EndTimestamp)

	assert.Equal(t, "03:04:05.006: TASKRUNNER: setup failed\n", string(logFile.Bytes()))
}
