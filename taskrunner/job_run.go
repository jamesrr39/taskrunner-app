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

/*
func (jobRun *JobRun) WriteToDisk() error {
	summaryFilePath := filepath.Join(jobRun.Path(), "summary.json")
	fileJson, err := json.Marshal(&jobRun)
	if nil != err {
		return err
	}

	log.Printf("serialising job run to disk: %s\n", fileJson)

	if err := ioutil.WriteFile(summaryFilePath, fileJson, 0644); nil != err {
		return err
	}

	return nil
}
*/
/*
// ~/.taskrunner/jobs/myjob/runs/{}
// creates the job run folder and any other folders, if necessary. Returns Id of the job, path to newly created job folder and an error.
func (job *Job) createNewJobRunFolder() (int, string, error) {
	jobId := 1
	for {
		path := filepath.Join(job.Path(), "runs", strconv.Itoa(jobId))
		if _, err := os.Stat(path); nil != err {
			if os.IsNotExist(err) {
				log.Printf("Creating new job run folder at %s\n", path)
				if err = os.MkdirAll(path, 0755); nil != err {
					return 0, "", err
				}
				log.Printf("allocating job path %s\n", path)
				return jobId, path, nil
			}
			return 0, "", err
		}

		if jobId > 1000000 {
			return 0, "", errors.New("Exceeded maximum amount of jobs allowed")
		}
		jobId++
	}
}

func (jobRun *JobRun) Path() string {
	return filepath.Join(jobRun.Job.GetRunsPath(), strconv.Itoa(jobRun.Id))
}

func (j *JobRun) SummaryPath() string {
	return filepath.Join(j.Path(), "summary.json")
}

func (j *JobRun) LogFilePath() string {
	return filepath.Join(j.Path(), "joboutput.log")
}
*/
