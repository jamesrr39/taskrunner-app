package taskrunnerdal

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
	"github.com/jamesrr39/taskrunner-app/taskexecutor"
	"github.com/jamesrr39/taskrunner-app/taskrunner"
)

type JobRunsDAL struct {
	jobDAL      *JobDAL
	nowProvider taskexecutor.NowProvider
}

func NewJobRunsDAL(jobDAL *JobDAL, nowProvider taskexecutor.NowProvider) *JobRunsDAL {
	return &JobRunsDAL{jobDAL, nowProvider}
}

// CreateAndRun creates a new job run and runs it.
// It is synchronous (the method returns when it has either completed or failed).
// jobRunStatusChangeExternalChan can be nil if there is no outside process listening for changes
func (jobRunsDAL *JobRunsDAL) CreateAndRun(jobRun *taskrunner.JobRun, jobRunStatusChangeExternalChan chan *taskrunner.JobRun) error {
	if jobRun.Id != 0 {
		return fmt.Errorf("expected the jobRun to have an Id of 0 but had %d", jobRun.Id)
	}

	jobRunStatusChanInternalChan := make(chan *taskrunner.JobRun)

	go jobRunsDAL.listenAndTriggerSendToJobRunStatusChan(jobRunStatusChanInternalChan, jobRunStatusChangeExternalChan, jobRun)

	err := jobRunsDAL.jobDAL.aquireWorkspaceLock(jobRun.Job)
	if nil != err {
		return err
	}
	defer jobRunsDAL.jobDAL.releaseWorkspaceLock(jobRun.Job)

	err = jobRunsDAL.jobDAL.ensureCleanWorkspaceDir(jobRun.Job)
	if nil != err {
		return err
	}

	jobRun.Id, err = jobRunsDAL.nextRunId(jobRun.Job)
	if nil != err {
		return err
	}

	// create dir for new job run
	runDir := filepath.Join(jobRunsDAL.jobDAL.getRunsDirForJob(jobRun.Job), strconv.FormatUint(uint64(jobRun.Id), 10))
	err = os.MkdirAll(runDir, 0700)
	if nil != err {
		return err
	}

	logFile, err := os.Create(filepath.Join(runDir, "joboutput.log"))
	if nil != err {
		return err
	}
	defer func() {
		err := logFile.Sync()
		if nil != err {
			log.Printf("ERROR: Failed to sync logfile to disk for job '%s' (id %d), run id %d. Error: %s\n",
				jobRun.Job.Name, jobRun.Job.Id, jobRun.Id, err)
		}
		logFile.Close()
	}()

	// set job run properties
	jobRun.StartTimestamp = time.Now().Unix()
	jobRun.State = taskrunner.JOB_RUN_STATE_IN_PROGRESS
	jobRunStatusChanInternalChan <- jobRun

	err = jobRunsDAL.writeJobRunSummary(jobRun)
	if nil != err {
		return err
	}

	err = taskexecutor.ExecuteJobRun(jobRun, jobRunStatusChanInternalChan, logFile, jobRunsDAL.jobDAL.getWorkspaceDir(jobRun.Job), jobRunsDAL.nowProvider)
	if nil != err {
		return err
	}

	return nil
}

func (jobRunsDAL *JobRunsDAL) listenAndTriggerSendToJobRunStatusChan(
	jobRunStatusChanInternalChan chan *taskrunner.JobRun,
	jobRunStatusChangeExternalChan chan *taskrunner.JobRun,
	jobRun *taskrunner.JobRun) {

	for {
		jobRun := <-jobRunStatusChanInternalChan
		err := jobRunsDAL.writeJobRunSummary(jobRun)
		if nil != err {
			log.Printf("ERROR: couldn't write job run %d for job '%s' (id: %d) to file. Error: %s\n", jobRun.Id, jobRun.Job.Name, jobRun.Job.Id, err)
		}
		if nil != jobRunStatusChangeExternalChan {
			jobRunStatusChangeExternalChan <- jobRun
		}
		if jobRun.State.IsFinished() {
			return
		}
	}
}

func (jobRunsDAL *JobRunsDAL) GetAllForJob(job *taskrunner.Job) ([]*taskrunner.JobRun, error) {
	runsDir := jobRunsDAL.jobDAL.getRunsDirForJob(job)
	fileinfos, err := ioutil.ReadDir(runsDir)
	if nil != err {
		return nil, err
	}

	var runs []*taskrunner.JobRun

	for _, fileinfo := range fileinfos {
		if !fileinfo.IsDir() {
			log.Printf("Found unexpected file in job runs directory: %s\n",
				filepath.Join(runsDir, fileinfo.Name()))
			continue
		}

		jobRunId, err := strconv.ParseUint(fileinfo.Name(), 10, 64)
		if nil != err {
			log.Printf("Found unexpected folder in job runs directory: %s\n", filepath.Join(runsDir, fileinfo.Name()))
			continue
		}

		run, err := jobRunsDAL.GetRun(job, jobRunId)
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

func (jobRunsDAL *JobRunsDAL) GetRun(job *taskrunner.Job, runId uint64) (*taskrunner.JobRun, error) {
	runDir := jobRunsDAL.getDirForJobRun(job, runId)

	summaryFile := filepath.Join(runDir, "summary.json")
	fileBytes, err := ioutil.ReadFile(summaryFile)
	if nil != err {
		return nil, err
	}

	var jobRun *taskrunner.JobRun
	err = json.Unmarshal(fileBytes, &jobRun)
	if nil != err {
		return nil, err
	}

	jobRun.Job = job

	return jobRun, nil
}

// GetLastRun returns the last job run, or nil if there weren't any
// todo optimisation?
func (jobRunsDAL *JobRunsDAL) GetLastRun(job *taskrunner.Job) (*taskrunner.JobRun, error) {
	jobRuns, err := jobRunsDAL.GetAllForJob(job)
	if nil != err {
		return nil, err
	}

	if 0 == len(jobRuns) {
		return nil, nil
	}

	lastRun := &taskrunner.JobRun{Id: 0}
	for _, jobRun := range jobRuns {
		if jobRun.Id > lastRun.Id {
			lastRun = jobRun
		}
	}
	return lastRun, nil
}

func (jobRunsDAL *JobRunsDAL) GetJobRunLog(jobRun *taskrunner.JobRun) (io.ReadCloser, error) {
	return os.Open(jobRunsDAL.GetJobRunLogLocation(jobRun))
}

// exposed for UI file truncated message
func (jobRunsDAL *JobRunsDAL) GetJobRunLogLocation(jobRun *taskrunner.JobRun) string {
	return filepath.Join(jobRunsDAL.getDirForJobRun(jobRun.Job, jobRun.Id), "joboutput.log")
}

func (jobRunsDAL *JobRunsDAL) getDirForJobRun(job *taskrunner.Job, runId uint64) string {
	return filepath.Join(jobRunsDAL.jobDAL.getRunsDirForJob(job), strconv.FormatUint(runId, 10))
}

func (jobRunsDAL *JobRunsDAL) nextRunId(job *taskrunner.Job) (uint64, error) {
	for i := 1; ; i++ {
		path := filepath.Join(jobRunsDAL.jobDAL.getRunsDirForJob(job), strconv.Itoa(i))
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return uint64(i), nil
			} else {
				return 0, err
			}
		}
	}
	return 0, errors.New("Unknown error generating new Job Run Id")
}

func (jobRunsDAL *JobRunsDAL) writeJobRunSummary(jobRun *taskrunner.JobRun) error {
	summaryBytes, err := json.Marshal(jobRun)
	if nil != err {
		return err
	}
	return ioutil.WriteFile(filepath.Join(jobRunsDAL.getDirForJobRun(jobRun.Job, jobRun.Id), "summary.json"), summaryBytes, 0700)
}
