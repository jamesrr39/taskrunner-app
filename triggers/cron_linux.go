//+build linux
package triggers

import (
	"bufio"
	"io"
	"os/exec"
	"strings"

	"github.com/jamesrr39/taskrunner-app/taskrunner"
)

type CronParser struct {
	commandPrefix string
}

func NewCronParser(commandPrefix string) *CronParser {
	return &CronParser{commandPrefix}
}

type CronJob struct {
	CronExpression string
	Command        string
}

func (parser *CronParser) SearchForJob(job *taskrunner.Job) ([]*CronJob, error) {
	jobs, err := parser.getCronOutput()
	if nil != err {
		return nil, err
	}

	var filteredJobs []*CronJob
	for _, cronJob := range jobs {
		if strings.HasPrefix(cronJob.Command, parser.commandPrefix+" "+job.Name) {
			filteredJobs = append(filteredJobs, cronJob)
		}
	}
	return filteredJobs, nil
}

// no jobs = err
func (parser *CronParser) getCronOutput() ([]*CronJob, error) {
	cmd := exec.Command("crontab", "-l")

	stdoutPipe, err := cmd.StdoutPipe()
	if nil != err {
		return nil, err
	}
	defer stdoutPipe.Close()

	stderrPipe, err := cmd.StderrPipe()
	if nil != err {
		return nil, err
	}
	defer stderrPipe.Close()

	var cronJobs []*CronJob
	go func() {
		cronJobs = parser.parseCronOutput(stdoutPipe)
	}()

	var errText string
	go func() {
		var errTextLines []string
		b := bufio.NewScanner(stderrPipe)
		for b.Scan() {
			errTextLines = append(errTextLines, b.Text())
		}
		errText = strings.Join(errTextLines, "\n")
	}()

	err = cmd.Run()
	if nil != err {
		if "" == errText {
			// an error
			return nil, err
		}
		// probably no jobs
		return nil, nil
	}

	return cronJobs, nil

}

func (parser *CronParser) parseCronOutput(reader io.Reader) []*CronJob {
	var cronJobs []*CronJob

	cronScanner := bufio.NewScanner(reader)
	for cronScanner.Scan() {
		text := strings.TrimSpace(cronScanner.Text())
		if strings.Index(text, "#") == 0 {
			continue // comment
		}
		fragments := strings.Fields(text)
		if len(fragments) < 6 {
			continue // too few fragments for cron expression
		}
		cronExpression := strings.Join(fragments[0:5], " ")
		command := strings.Join(fragments[5:], " ")

		if !strings.HasPrefix(command, parser.commandPrefix) {
			continue
		}

		cronJobs = append(cronJobs, &CronJob{cronExpression, command})
	}

	return cronJobs
}
