package taskrunner

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bradfitz/slice"
)

const TASKRUNNER_SOURCE_NAME string = "TASKRUNNER"

type Script string

type Job struct {
	Id                 uint                `json:"-"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Script             Script              `json:"script"`
	TaskrunnerInstance *TaskrunnerInstance `json:"-"`
}

func NewJob(name string, description string, script Script, taskrunnerInstance *TaskrunnerInstance) (*Job, error) {
	if "" == name {
		return nil, errors.New("A job must have a name")
	}

	id, err := taskrunnerInstance.nextId()
	if nil != err {
		return nil, err
	}

	return &Job{Id: id, Name: name, Description: description, Script: script, TaskrunnerInstance: taskrunnerInstance}, nil
}

func (job *Job) Path() string {
	return filepath.Join(job.TaskrunnerInstance.Basepath, "jobs", strconv.Itoa(int(job.Id)))
}

// cleans previous workspace, creates folder for new workspace
func (job *Job) cleanWorkspace() error {
	workspacePath := job.workspacePath()
	err := os.RemoveAll(workspacePath)
	if nil != err {
		return err
	}
	err = os.MkdirAll(workspacePath, 0700)
	return err
}

func (job *Job) workspacePath() string {
	return filepath.Join(job.Path(), "workspace")
}

func (job *Job) prepareAndMoveToWorkspace() error {
	err := job.cleanWorkspace()
	if nil != err {
		return errors.New("Couldn't clean workspace")
	}
	err = os.Chdir(job.workspacePath())
	if nil != err {
		return errors.New("Couldn't change directory to workspace directory (" + job.workspacePath() + "). Error: " + err.Error())
	}
	return nil
}

// handles errors caused by the internal workings of taskrunner (not the script)
func handleTaskrunnerError(errorMessage string, logFile io.Writer, jobRun *JobRun) {
	jobRun.State = JOB_RUN_STATE_FAILED
	writeStringToLogFile(errorMessage, logFile, TASKRUNNER_SOURCE_NAME)
}

// change dir to workspace, write script file run
// error would be an error that prevents the job from logging
func (job *Job) Run(trigger TriggerType) (*JobRun, error) {
	jobRunId, _, err := job.createNewJobRunFolder()
	if nil != err {
		return nil, err
	}

	jobRun := &JobRun{Id: jobRunId, Job: job, StartTimestamp: time.Now().Unix(), Trigger: trigger, State: JOB_RUN_STATE_IN_PROGRESS}

	logFile, err := os.Create(jobRun.LogFilePath())
	if nil != err {
		return nil, err
	}
	defer logFile.Close()

	err = job.prepareAndMoveToWorkspace()
	if nil != err {
		handleTaskrunnerError("Couldn't prepare and move to workspace. Error: "+err.Error(), logFile, jobRun)
		return jobRun, nil
	}

	scriptFilePath := filepath.Join(job.workspacePath(), "script")
	err = ioutil.WriteFile(scriptFilePath, []byte(job.Script), 0700)
	if nil != err {
		handleTaskrunnerError("Couldn't prepare and move to workspace. Error: "+err.Error(), logFile, jobRun)
		return jobRun, nil
	}

	cmd := exec.Command(scriptFilePath)
	stdoutPipe, err := cmd.StdoutPipe()
	if nil != err {
		handleTaskrunnerError("Couldn't obtain stdoutpipe. Error: "+err.Error(), logFile, jobRun)
		return jobRun, nil
	}
	stderrPipe, err := cmd.StderrPipe()
	if nil != err {
		handleTaskrunnerError("Couldn't obtain stderrpipe. Error: "+err.Error(), logFile, jobRun)
		return jobRun, nil
	}

	go writeToLogFile(stdoutPipe, logFile, "STDOUT")
	go writeToLogFile(stderrPipe, logFile, "STDERR")

	err = cmd.Start()
	if nil != err {
		handleTaskrunnerError("Couldn't start script. Error: "+err.Error(), logFile, jobRun)
		return jobRun, nil
	}

	jobRun.Pid = cmd.Process.Pid
	writeToDiskErr := jobRun.WriteToDisk()

	err = cmd.Wait()
	if nil != err {
		switch err.(type) {
		case *exec.ExitError:
			jobRun.State = JOB_RUN_STATE_FAILED
		default:
			jobRun.State = JOB_RUN_STATE_UNKNOWN
		}
	} else {
		jobRun.State = JOB_RUN_STATE_SUCCESS
	}

	jobRun.EndTimestamp = time.Now().Unix()

	if nil != writeToDiskErr {
		return nil, writeToDiskErr
	}

	if err = jobRun.WriteToDisk(); err != nil {
		return nil, err
	}

	logFile.Sync()

	return jobRun, nil
}

func (job *Job) GetRunsPath() string {
	return filepath.Join(job.Path(), "runs")
}

func (job *Job) GetRuns() ([]*JobRun, error) {
	fileinfos, err := ioutil.ReadDir(job.GetRunsPath())
	if nil != err {
		return nil, err
	}

	var runs []*JobRun

	for _, fileinfo := range fileinfos {
		if !fileinfo.IsDir() {
			log.Printf("Found unexpected file in job runs directory: %s\n", filepath.Join(job.GetRunsPath(), fileinfo.Name()))
			continue
		}

		jobRunId, err := strconv.Atoi(fileinfo.Name())
		if nil != err {
			log.Printf("Found unexpected folder in job runs directory: %s\n", filepath.Join(job.GetRunsPath(), fileinfo.Name()))
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

	log.Printf("looking in %s\n", job.GetRunsPath())

	files, err := ioutil.ReadDir(job.GetRunsPath())
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
	runDir := filepath.Join(job.GetRunsPath(), strconv.Itoa(id))

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
