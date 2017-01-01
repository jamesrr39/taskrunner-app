package taskrunner

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jamesrr39/goutil/user"
)

type TaskrunnerInstance struct {
	Basepath string
}

func NewTaskrunnerInstanceAndEnsurePaths(basepath string) (*TaskrunnerInstance, error) {
	log.Printf("taskrunner base path: '%s'\n", basepath)

	expandedBasepath, err := user.ExpandUser(basepath)
	if nil != err {
		return nil, err
	}

	t := &TaskrunnerInstance{Basepath: expandedBasepath}

	err = t.ensureDirectories()
	if nil != err {
		return nil, err
	}

	return t, nil
}

func (taskrunnerInstance *TaskrunnerInstance) jobsDir() string {
	return filepath.Join(taskrunnerInstance.Basepath, "jobs")
}

func (taskrunnerInstance *TaskrunnerInstance) nextId() (uint, error) {
	for i := 1; ; i++ {
		path := filepath.Join(taskrunnerInstance.jobsDir(), strconv.Itoa(i))
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return uint(i), nil
			} else {
				return 0, err
			}
		}
	}
	//return 0, errors.New("Unknown error generating new Job Id")
}

func (t *TaskrunnerInstance) ensureDirectories() error {

	err := os.MkdirAll(t.jobsDir(), 0700)
	return err
}

func (t *TaskrunnerInstance) GetJobByName(name string) (*Job, error) {
	// todo optimisation?
	allJobs, err := t.GetAllJobs()
	if nil != err {
		return nil, err
	}

	for _, job := range allJobs {
		if job.Name == name {
			return job, nil
		}
	}

	return nil, &ErrJobNotFound{}
}
