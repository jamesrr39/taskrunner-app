package taskexecutor

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"
)

func getFormattedTime(now time.Time) string {
	milliSeconds := int(float64(now.Nanosecond()) / float64(1000000))
	return fmt.Sprintf("%02d:%02d:%02d.%03d", now.Hour(), now.Minute(), now.Second(), milliSeconds)
}

func writeStringToLogFile(text string, writer io.Writer, sourceName string) {
	writer.Write([]byte(getFormattedTime(time.Now()) + ": " + sourceName + ": " + text + "\n"))
}

func writeToLogFile(pipe io.Reader, writer io.Writer, sourceName string) {
	pipeScanner := bufio.NewScanner(pipe)
	for pipeScanner.Scan() {
		writer.Write([]byte(getFormattedTime(time.Now()) + ": " + sourceName + ": " + pipeScanner.Text() + "\n"))
	}
}
