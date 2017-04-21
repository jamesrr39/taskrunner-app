package taskrunnerdal

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/jamesrr39/taskrunner-app/taskrunner"
	"github.com/stretchr/testify/assert"
)

func generateTestTaskrunerDAL(t *testing.T) *TaskrunnerDAL {
	tempdirPath, err := ioutil.TempDir("", "")
	assert.Nil(t, err)

	taskrunnerDAL, err := NewTaskrunnerDALAndEnsureDirectories(tempdirPath)
	assert.Nil(t, err)

	return taskrunnerDAL
}

func cleanUpTaskrunnerDAL(t *testing.T, taskrunnerDAL *TaskrunnerDAL) {
	err := os.RemoveAll(taskrunnerDAL.basePath)
	assert.Nil(t, err)
}

var testScript = taskrunner.Script(`#!/bin/bash

echo t
`)

func Test_CreateAndRun_noExternalChan(t *testing.T) {
	taskrunnerDAL := generateTestTaskrunerDAL(t)

	const testJobName = "test job"
	const testJobDesc = "my job for testing"
	testJobRunTrigger := taskrunner.TriggerType("test trigger")

	job, err := taskrunner.NewJob(0, testJobName, testJobDesc, testScript)
	assert.Nil(t, err)

	err = taskrunnerDAL.JobDAL.Create(job)
	assert.Nil(t, err)

	err = taskrunnerDAL.JobRunsDAL.CreateAndRun(job.NewJobRun(testJobRunTrigger), nil)
	assert.Nil(t, err)

	jobRuns, err := taskrunnerDAL.JobRunsDAL.GetAllForJob(job)
	assert.Nil(t, err)

	assert.Len(t, jobRuns, 1)

	assert.Equal(t, testJobName, jobRuns[0].Job.Name)
	assert.Equal(t, testJobDesc, jobRuns[0].Job.Description)

	assert.Equal(t, testJobRunTrigger, jobRuns[0].Trigger)
	assert.Equal(t, taskrunner.JOB_RUN_STATE_SUCCESS, jobRuns[0].State)

	//cleanUpTaskrunnerDAL(t, taskrunnerDAL)
}
