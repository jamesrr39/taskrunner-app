package taskrunner

/*
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)
*/
/*
// creates a new job, and adds an Id property onto the job passed in.
func (t *TaskrunnerInstance) CreateJob(job *Job) error {

	// check for existing Id
	if job.Id != 0 {
		return fmt.Errorf("This job already has an Id property (%d). It should use saveJob instead.", job.Id)
	}

	// check for duplicate job names
	existingJobs, err := t.GetAllJobs()
	if nil != err {
		return err
	}

	for _, existingJob := range existingJobs {
		if existingJob.Name == job.Name {
			return fmt.Errorf("Job name: %s already taken by job Id %d.", job.Name, existingJob.Id)
		}
	}

	// generate a new Id
	job.Id, err = job.TaskrunnerInstance.nextId()
	if nil != err {
		return err
	}

	return t.SaveJob(job)

}
*/
/*
func (t *TaskrunnerInstance) SaveJob(job *Job) error {

	if job.Id == 0 {
		return fmt.Errorf("This job doesn't have an Id. Use `CreateJob` to create a new job and generate an Id.")
	}

	// create new job folder and write config to it.
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
*/
/*
func (taskrunnerInstance *TaskrunnerInstance) GetAllJobs() ([]*Job, error) {
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

			jobId, err := strconv.Atoi(folderItem.Name())
			if nil != err {
				log.Printf("Couldn't get a job Id for %s. Error: %s\n", jobConfigPath, err)
				continue
			}
			job.Id = uint(jobId)

			log.Printf("job: %v\n", job)

			job.TaskrunnerInstance = taskrunnerInstance

			jobs = append(jobs, job)
		}
	}
	return jobs, nil

}
*/
