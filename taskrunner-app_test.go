package main

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/jamesrr39/taskrunner-app/taskrunnerdal"
	"github.com/stretchr/testify/assert"
)

func Test_runJobHeadlessMain(t *testing.T) {
	// create data dir and defer cleaning up data dir
	dataDirBase, err := ioutil.TempDir("", "")
	assert.Nil(t, err)
	defer func() {
		err := os.RemoveAll(dataDirBase)
		if nil != err {
			t.Errorf("Couldn't remove the tempdir at '%s'. Error: %s.\n", dataDirBase, err)
		} else {
			t.Logf("Successfully removed data dir at %s\n", dataDirBase)
		}
	}()

	// create dal and job
	dal, err := taskrunnerdal.NewTaskrunnerDALAndEnsureDirectories(dataDirBase, mockNowProvider)
	assert.Nil(t, err)

	var script taskrunner.Script = `#!/bin/bash
echo 'test run'
`

	job, err := taskrunner.NewJob(0, "system test job", "", script)
	err = dal.JobDAL.Create(job)
	assert.Nil(t, err)

	// run job and assertions
	runJobHeadlessMain(dal, "system test job", "system test")

	runs, err := dal.JobRunsDAL.GetAllForJob(job)
	assert.Nil(t, err)
	assert.Len(t, runs, 1)
	assert.Equal(t, taskrunner.JOB_RUN_STATE_SUCCESS, runs[0].State)

	logFile, err := dal.JobRunsDAL.GetJobRunLog(runs[0])
	assert.Nil(t, err)
	defer logFile.Close()
	logFileBytes, err := ioutil.ReadAll(logFile)
	assert.Nil(t, err)

	assert.Equal(t, "03:04:05.006: STDOUT: test run\n", string(logFileBytes))

}

// cmd := exec.Command("go", "run", "taskrunner-app-main.go", "--taskrunner-dir='"+dataDirBase+"'", "--run-job='system test job'", "--trigger='system test'")
//	err = cmd.Run()

func mockNowProvider() time.Time {
	return time.Date(2000, 01, 02, 03, 04, 05, 6*1000*1000, time.UTC)
}
