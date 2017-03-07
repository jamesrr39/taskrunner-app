package taskrunner

/*
import (
	"log"

	"taskrunner-app/taskrunnerdal"

	"github.com/jamesrr39/goutil/user"
)

type TaskrunnerInstance struct {
	*taskrunnerdal.TaskrunnerDAL
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
*/
/*
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
*/
