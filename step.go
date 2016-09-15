package taskrunner

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type Step struct {
	Name string `json:name`
	Cmd  string `json:cmd`
}

func NewStep(name string, cmd string) *Step {
	return &Step{Name: name, Cmd: cmd}
}

func (s *Step) String() string {
	return fmt.Sprintf("Step[Cmd=%s", s.Cmd)
}

func (step *Step) Run(jobRunPath string, stdoutFile io.Writer, stderrFile io.Writer) error {
	cmd := exec.Command("bash", "-c", step.Cmd) // todo - bash?
	stdoutPipe, err := cmd.StdoutPipe()
	if nil != err {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if nil != err {
		return err
	}

	go func() {
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			stdoutFile.Write([]byte(text + "\n"))
			fmt.Printf("%s\n", text)
		}
	}()

	go func() {
		stderrScanner := bufio.NewScanner(stderrPipe)
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			stderrFile.Write([]byte(text))
			fmt.Printf("%s\n", text)
		}
	}()

	err = cmd.Start()
	if nil != err {
		return err
	}

	err = cmd.Wait()
	if nil != err {
		return err
	}
	return nil
}
