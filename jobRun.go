package taskrunner

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type JobRun struct {
	Id             int    `json:"id"`
	Successful     bool   `json:"successful"`
	StartTimestamp int64  `json:"startTimestamp"`
	EndTimestamp   int64  `json:"endTimestamp"`
	Trigger        string `json:"trigger"`
	Job            *Job   `json:"-"`
}

func (job *Job) NewJobRun(id int, successful bool, startTimestamp int64, endTimestamp int64, trigger string) *JobRun {
	return &JobRun{Id: id, Successful: successful, StartTimestamp: startTimestamp, EndTimestamp: endTimestamp, Trigger: trigger, Job: job}
}

func (jobRun *JobRun) WriteToDisk(jobRunPath string) error {
	summaryFilePath := filepath.Join(jobRunPath, "summary.json")
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
func createNewJobRunFolder(jobBasePath string) (int, string, error) {
	jobId := 1
	for {
		path := filepath.Join(jobBasePath, "runs", strconv.Itoa(jobId))
		if _, err := os.Stat(path); nil != err {
			if os.IsNotExist(err) {
				if err = os.MkdirAll(path, 0755); nil != err {
					return 0, "", err
				}
				log.Printf("allocating job path %s\n", path)
				return jobId, path, nil
			}
			return 0, "", err
		}

		if jobId > 100000 {
			panic("Exceeded maximum amount of jobs allowed") // todo replace folder scanning with a count
		}
		jobId++
	}
}

func (jobRun *JobRun) StdoutLogPath() string {
	return filepath.Join(jobRun.Path(), "stdout.log")
}

func (jobRun *JobRun) StderrLogPath() string {
	return filepath.Join(jobRun.Path(), "stderr.log")
}

func (jobRun *JobRun) Path() string {
	return filepath.Join(jobRun.Job.RunsPath(), strconv.Itoa(jobRun.Id))
}
