package taskrunner

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type JobRunState int

const (
	JOB_RUN_STATE_UNKNOWN JobRunState = iota
	JOB_RUN_STATE_FAILED
	JOB_RUN_STATE_SUCCESS
	JOB_RUN_STATE_IN_PROGRESS
)

var jobRunStates = [...]string{
	"Unknown",
	"Failed",
	"Success",
	"In Progress",
}

func (e JobRunState) String() string {
	return jobRunStates[e]
}

type TriggerType string

type JobRun struct {
	Id             int         `json:"id"`
	State          JobRunState `json:"status"`
	StartTimestamp int64       `json:"startTimestamp"`
	EndTimestamp   int64       `json:"endTimestamp,omitempty"`
	Trigger        TriggerType `json:"trigger"`
	Job            *Job        `json:"-"`
	Pid            int         `json:"pid"`
	ExitCode       int         `json:"exitCode"`
}

func (job *Job) NewJobRun(id int, state JobRunState, startTimestamp int64, endTimestamp int64, trigger TriggerType, pid int, exitCode int) *JobRun {
	return &JobRun{Id: id, State: state, StartTimestamp: startTimestamp, EndTimestamp: endTimestamp, Trigger: trigger, Job: job, Pid: pid, ExitCode: exitCode}
}

func (jobRun *JobRun) WriteToDisk() error {
	summaryFilePath := filepath.Join(jobRun.Path(), "summary.json")
	fileJson, err := json.Marshal(&jobRun)
	if nil != err {
		return err
	}

	log.Printf("serialising job run to disk: %s\n", fileJson)

	if err := ioutil.WriteFile(summaryFilePath, fileJson, 0644); nil != err {
		return err
	}

	return nil
}

// ~/.taskrunner/jobs/myjob/runs/{}
// creates the job run folder and any other folders, if necessary. Returns Id of the job, path to newly created job folder and an error.
func (job *Job) createNewJobRunFolder() (int, string, error) {
	jobId := 1
	for {
		path := filepath.Join(job.Path(), "runs", strconv.Itoa(jobId))
		if _, err := os.Stat(path); nil != err {
			if os.IsNotExist(err) {
				log.Printf("Creating new job run folder at %s\n", path)
				if err = os.MkdirAll(path, 0755); nil != err {
					return 0, "", err
				}
				log.Printf("allocating job path %s\n", path)
				return jobId, path, nil
			}
			return 0, "", err
		}

		if jobId > 1000000 {
			panic("Exceeded maximum amount of jobs allowed") // todo replace folder scanning with a count
		}
		jobId++
	}
}

func (jobRun *JobRun) Path() string {
	return filepath.Join(jobRun.Job.GetRunsPath(), strconv.Itoa(jobRun.Id))
}

func (j *JobRun) SummaryPath() string {
	return filepath.Join(j.Path(), "summary.json")
}

func (j *JobRun) LogFilePath() string {
	return filepath.Join(j.Path(), "joboutput.log")
}
