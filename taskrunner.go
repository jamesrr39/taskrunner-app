package taskrunner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jamesrr39/goutil/user"
)

type TaskrunnerInstance struct {
	Basepath string
}

func NewTaskrunnerInstance(basepath string) (*TaskrunnerInstance, error) {
	log.Printf("taskrunner base path: '%s'\n", basepath)

	expandedBasepath, err := user.ExpandUser(basepath)
	if nil != err {
		return nil, err
	}

	return &TaskrunnerInstance{
		Basepath: expandedBasepath,
	}, nil
}

func (taskrunnerInstance *TaskrunnerInstance) jobsDir() string {
	return filepath.Join(taskrunnerInstance.Basepath, "jobs")
}

// for reading a config file
func (taskrunnerInstance *TaskrunnerInstance) JobFromName(jobName string) (*Job, error) {
	fmt.Printf("jobname: %s\n", jobName)
	path := filepath.Join(taskrunnerInstance.Basepath, "jobs", jobName, "config.json")

	fileBytes, err := ioutil.ReadFile(path)
	if nil != err {
		return nil, fmt.Errorf("Couldn't read file at '%s'. Error: %s\n", path, err)
	}

	var job *Job
	err = json.Unmarshal(fileBytes, &job)
	if nil != err {
		return nil, fmt.Errorf("Couldn't unmarshal '%s'. Error: %s\n", path, err)
	}

	if jobName != job.Name {
		return nil, fmt.Errorf("Job names don't match. Expected '%s' from job path but found '%s'\n", jobName, job.Name)
	}

	job.TaskrunnerInstance = taskrunnerInstance

	return job, nil
}

func (taskrunnerInstance *TaskrunnerInstance) Jobs() ([]*Job, error) {
	fmt.Printf("getting jobs dir\n")
	jobsDir := taskrunnerInstance.jobsDir()
	fmt.Printf("jobsDir: %s\n", jobsDir)

	var jobs []*Job
	folderItems, err := ioutil.ReadDir(jobsDir)
	if nil != err {
		return nil, err
	}

	for _, folderItem := range folderItems {
		if !folderItem.IsDir() {
			continue
		}
		jobConfigPath := filepath.Join(jobsDir, folderItem.Name(), "config.json")
		if _, err := os.Stat(jobConfigPath); nil == err {
			var job *Job
			jobFile, err := ioutil.ReadFile(jobConfigPath)
			if nil != err {
				log.Printf("Error reading file: %s. Error: %s\n", jobConfigPath, err)
				continue
			}

			err = json.Unmarshal(jobFile, &job)
			if nil != err {
				log.Printf("Error unmarshalling json from %s. Error: %s\n", jobConfigPath, err)
				continue
			}

			job.TaskrunnerInstance = taskrunnerInstance

			jobs = append(jobs, job)
		}
	}
	return jobs, nil

}
