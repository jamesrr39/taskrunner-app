package taskrunnerdal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"github.com/jamesrr39/taskrunner-app/taskrunner"
)

var ErrLockfileAlreadyExists = errors.New("lockfile already exists")

type JobDAL struct {
	jobsDir string
}

func NewJobDAL(jobsDir string) *JobDAL {
	return &JobDAL{jobsDir}
}

func (jobDAL *JobDAL) GetAll() ([]*taskrunner.Job, error) {

	var jobs []*taskrunner.Job
	folderItems, err := ioutil.ReadDir(jobDAL.jobsDir)
	if nil != err {
		return nil, err
	}

	for _, folderItem := range folderItems {
		if !folderItem.IsDir() {
			continue
		}
		jobConfigPath := filepath.Join(jobDAL.jobsDir, folderItem.Name(), "config.json")
		if _, err := os.Stat(jobConfigPath); nil == err {
			var job *taskrunner.Job
			jobFileBytes, err := ioutil.ReadFile(jobConfigPath)
			if nil != err {
				log.Printf("Error reading file: %s. Error: %s\n", jobConfigPath, err)
				continue
			}

			err = json.Unmarshal(jobFileBytes, &job)
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

			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

func (jobDAL *JobDAL) Create(job *taskrunner.Job) error {

	// check for existing Id
	if job.Id != 0 {
		return fmt.Errorf("This job already has an Id property (%d). It should use saveJob instead.", job.Id)
	}

	// generate a new Id
	var err error
	job.Id, err = jobDAL.nextId()
	if nil != err {
		return err
	}

	err = jobDAL.save(job)
	if nil != err {
		return err
	}

	return jobDAL.ensureRunsFolder(job)

}

func (jobDAL *JobDAL) Update(job *taskrunner.Job) error {

	if job.Id == 0 {
		return fmt.Errorf("This job doesn't have an Id. Use `CreateJob` to create a new job and generate an Id.")
	}

	return jobDAL.save(job)
}

func (jobDAL *JobDAL) save(job *taskrunner.Job) error {
	// check for duplicate job names
	existingJobs, err := jobDAL.GetAll()
	if nil != err {
		return err
	}

	for _, existingJob := range existingJobs {
		if existingJob.Name == job.Name && existingJob.Id != job.Id {
			return fmt.Errorf("Job name: %s already taken by job Id %d.", job.Name, existingJob.Id)
		}
	}

	// create new job folder and write config to it.
	jobFolderPath := jobDAL.getJobPath(job)
	err = os.MkdirAll(jobFolderPath, 0700)
	if nil != err {
		return err
	}

	fileBytes, err := json.MarshalIndent(job, "", "	")
	if nil != err {
		return err
	}
	log.Printf("writing to %s\n", filepath.Join(jobFolderPath, "config.json"))
	return ioutil.WriteFile(filepath.Join(jobFolderPath, "config.json"), fileBytes, 0600)
}

func (jobDAL *JobDAL) nextId() (uint, error) {
	for i := 1; ; i++ {
		path := filepath.Join(jobDAL.jobsDir, strconv.Itoa(i))
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return uint(i), nil
			} else {
				return 0, err
			}
		}
	}
	return 0, errors.New("Unknown error generating new Job Id")
}

func (jobDAL *JobDAL) getJobPath(job *taskrunner.Job) string {
	return filepath.Join(jobDAL.jobsDir, strconv.FormatUint(uint64(job.Id), 10))
}

func (jobDAL *JobDAL) getWorkspaceDir(job *taskrunner.Job) string {
	return filepath.Join(jobDAL.getJobPath(job), "workspace")
}

func (jobDAL *JobDAL) ensureCleanWorkspaceDir(job *taskrunner.Job) error {
	workspaceDir := jobDAL.getWorkspaceDir(job)
	err := os.RemoveAll(workspaceDir)
	if nil != err {
		return err
	}
	err = os.MkdirAll(workspaceDir, 0700)
	if nil != err {
		return err
	}
	return nil
}

func (jobDAL *JobDAL) aquireWorkspaceLock(job *taskrunner.Job) error {
	lockfilePath := jobDAL.getLockfilePath(job)
	if _, err := os.Stat(lockfilePath); err != nil {
		if nil == err {
			log.Println("lockfile path: " + lockfilePath)
			return ErrLockfileAlreadyExists
		} else if nil != err {
			if !os.IsNotExist(err) {
				return err
			}
			// if err is os.IsNotExist, continue
		}
	}
	file, err := os.Create(lockfilePath)
	if nil != err {
		return err
	}
	defer file.Close()
	return nil
}

func (jobDAL *JobDAL) releaseWorkspaceLock(job *taskrunner.Job) error {
	return os.Remove(jobDAL.getLockfilePath(job))
}

func (jobDAL *JobDAL) getLockfilePath(job *taskrunner.Job) string {
	return filepath.Join(jobDAL.getJobPath(job), "lockfile.txt")
}

func (jobDAL *JobDAL) getRunsDirForJob(job *taskrunner.Job) string {
	return filepath.Join(jobDAL.getJobPath(job), "runs")
}

func (jobDAL *JobDAL) ensureRunsFolder(job *taskrunner.Job) error {
	return os.MkdirAll(jobDAL.getRunsDirForJob(job), 0700)
}
