package taskrunnerdal

import (
	"os"
	"path/filepath"

	"github.com/jamesrr39/taskrunner-app/taskexecutor"
)

type TaskrunnerDAL struct {
	basePath string
	jobsPath string
	*JobDAL
	*JobRunsDAL
}

const jobsDirName string = "jobs"

func NewTaskrunnerDALAndEnsureDirectories(basePath string, nowProvider taskexecutor.NowProvider) (*TaskrunnerDAL, error) {
	jobsPath := filepath.Join(basePath, jobsDirName)
	jobDAL := NewJobDAL(jobsPath)

	dal := &TaskrunnerDAL{basePath, jobsPath, jobDAL, NewJobRunsDAL(jobDAL, nowProvider)}

	err := dal.ensureDirectories()
	if nil != err {
		return nil, err
	}

	return dal, nil
}

func (dal *TaskrunnerDAL) ensureDirectories() error {
	err := os.MkdirAll(dal.jobsPath, 0700)
	return err
}

func (dal *TaskrunnerDAL) String() string {
	return dal.basePath
}
