package taskrunner

type JobRunState int

const (
	JOB_RUN_STATE_UNKNOWN JobRunState = iota
	JOB_RUN_STATE_FAILED
	JOB_RUN_STATE_SUCCESS
	JOB_RUN_STATE_IN_PROGRESS
	JOB_RUN_STATE_NOT_STARTED
)

var jobRunStates = [...]string{
	"Unknown",
	"Failed",
	"Success",
	"In Progress",
	"Not Started",
}

func (e JobRunState) String() string {
	return jobRunStates[e]
}

func (e JobRunState) IsFinished() bool {
	switch e {
	case JOB_RUN_STATE_SUCCESS, JOB_RUN_STATE_FAILED:
		return true
	default:
		return false
	}
}

type TriggerType string

type JobRun struct {
	Id             uint64      `json:"id"`
	State          JobRunState `json:"status"`
	StartTimestamp int64       `json:"startTimestamp"`
	EndTimestamp   int64       `json:"endTimestamp,omitempty"`
	Trigger        TriggerType `json:"trigger"`
	Job            *Job        `json:"-"`
	Pid            *int        `json:"pid"`      // nil for not started
	ExitCode       *int        `json:"exitCode"` // nil for not started
}

func (job *Job) NewJobRun(trigger TriggerType) *JobRun {
	return &JobRun{
		Id:             0,
		State:          JOB_RUN_STATE_NOT_STARTED,
		StartTimestamp: 0,
		EndTimestamp:   0,
		Trigger:        trigger,
		Job:            job,
		Pid:            nil,
		ExitCode:       nil,
	}
}
