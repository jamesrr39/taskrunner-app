package taskrunner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bradfitz/slice"
)

type Job struct {
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Steps              []*Step             `json:"steps"`
	TaskrunnerInstance *TaskrunnerInstance `json:"-"`
}

func NewJob(name string, description string, steps []*Step, taskrunnerInstance *TaskrunnerInstance) (*Job, error) {
	if "" == name {
		return nil, errors.New("A job must have a name")
	}

	return &Job{Name: name, Description: description, Steps: steps, TaskrunnerInstance: taskrunnerInstance}, nil
}

func (job *Job) Path() string {
	return filepath.Join(job.TaskrunnerInstance.Basepath, "jobs", job.Name)
}

func (j *Job) String() string {
	return fmt.Sprintf("Job[Name=%s, Description=%s, Steps=%v]", j.Name, j.Description, j.Steps)
}

func (job *Job) Run(trigger string) error {
	jobBasepath := job.Path()
	jobId, jobRunPath, err := createNewJobRunFolder(filepath.Join(jobBasepath, job.Name))
	if nil != err {
		return err
	}

	stdoutFile, err := os.Create(filepath.Join(jobRunPath, "stdout.log"))
	if nil != err {
		return err
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(filepath.Join(jobRunPath, "stderr.log"))
	if nil != err {
		return err
	}
	defer stderrFile.Close()

	startTimestamp := time.Now().Unix()

	successful := false
	if err := job.runSteps(jobRunPath, stdoutFile, stderrFile); nil == err {
		successful = true
	}

	stdoutFile.Sync()
	stderrFile.Sync()

	endTimestamp := time.Now().Unix()

	jobRun := job.NewJobRun(jobId, successful, startTimestamp, endTimestamp, trigger)
	if err = jobRun.WriteToDisk(jobRunPath); err != nil {
		return err
	}

	return nil
}

func (job *Job) runSteps(jobRunPath string, stdoutFile io.Writer, stderrFile io.Writer) error {
	for index, step := range job.Steps {
		stdoutFile.Write([]byte("Starting step " + strconv.Itoa(index) + "\n"))
		err := step.Run(jobRunPath, stdoutFile, stderrFile)
		if nil != err {
			log.Printf("Error executing step %d: '%s'. Error: %s\n", index, step.Cmd, err)
			return err
		}
		stdoutFile.Write([]byte("Finished step " + strconv.Itoa(index) + "\n"))
	}
	return nil
}

func (job *Job) RunsPath() string {
	return filepath.Join(job.Path(), "runs")
}

func (job *Job) Runs() ([]*JobRun, error) {
	fileinfos, err := ioutil.ReadDir(job.RunsPath())
	if nil != err {
		return nil, err
	}

	var runs []*JobRun

	for _, fileinfo := range fileinfos {
		if !fileinfo.IsDir() {
			log.Printf("Found unexpected file in job runs directory: %s\n", filepath.Join(job.RunsPath(), fileinfo.Name()))
			continue
		}

		jobRunId, err := strconv.Atoi(fileinfo.Name())
		if nil != err {
			log.Printf("Found unexpected folder in job runs directory: %s\n", filepath.Join(job.RunsPath(), fileinfo.Name()))
			continue
		}

		run, err := job.GetRun(jobRunId)
		if nil != err {
			log.Println(err)
			continue
		}

		runs = append(runs, run) // todo go func

	}

	slice.Sort(runs, func(i, j int) bool {
		return runs[i].Id > runs[j].Id
	})

	return runs, nil

}

// 0 if no runs
func (job *Job) GetLastRunId() int {
	numberOfRuns := 0

	files, err := ioutil.ReadDir(job.RunsPath())
	if err != nil {
		return numberOfRuns // todo err type checking
	}

	for _, file := range files {
		runNumber, err := strconv.Atoi(file.Name())
		if nil != err {
			continue
		}
		if runNumber > numberOfRuns {
			numberOfRuns = runNumber
		}
	}
	return numberOfRuns

}

func (job *Job) GetRun(id int) (*JobRun, error) {
	runDir := filepath.Join(job.RunsPath(), strconv.Itoa(id))

	summaryFile := filepath.Join(runDir, "summary.json")
	fileBytes, err := ioutil.ReadFile(summaryFile)
	if nil != err {
		return nil, err
	}

	var jobRun *JobRun
	err = json.Unmarshal(fileBytes, &jobRun)
	if nil != err {
		return nil, err
	}

	jobRun.Job = job

	return jobRun, nil
}

func (job *Job) Save() error {
	jobFolderPath := job.Path()
	err := os.MkdirAll(jobFolderPath, 0700)
	if nil != err {
		return err
	}

	fileBytes, err := json.MarshalIndent(job, "", "	")
	if nil != err {
		return err
	}

	return ioutil.WriteFile(filepath.Join(jobFolderPath, "config.json"), fileBytes, 0600)
}
