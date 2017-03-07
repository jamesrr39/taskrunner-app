package taskexecutor

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func getFormattedTime() string {
	now := time.Now()
	milliSeconds := int(float64(now.Nanosecond()) / float64(1000000))
	return fmt.Sprintf("%02d:%02d:%02d.%03d", now.Hour(), now.Minute(), now.Second(), milliSeconds)
}

func writeStringToLogFile(text string, writer io.Writer, sourceName string) {
	writer.Write([]byte(getFormattedTime() + ": " + sourceName + ": " + text + "\n"))
}

func writeToLogFile(pipe io.Reader, writer io.Writer, sourceName string) {
	pipeScanner := bufio.NewScanner(pipe)
	for pipeScanner.Scan() {
		writer.Write([]byte(getFormattedTime() + ": " + sourceName + ": " + pipeScanner.Text() + "\n"))
	}
}
