package gui

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/mattn/go-gtk/gtk"
)

func fillTextBufferFromFile(textBuffer *gtk.TextBuffer, fileReader io.Reader, linesToRead uint, jobRunLogLocation string) {
	fileScanner := bufio.NewScanner(fileReader)
	linesRead := uint(0)
	for fileScanner.Scan() {
		if linesRead >= linesToRead {
			textBuffer.InsertAtCursor(fmt.Sprintf("\n\nOutput truncated as it exceeds %d lines.\nFull output can be found at %s\n",
				linesRead,
				jobRunLogLocation))
			return
		}
		text := fileScanner.Text()
		textBuffer.InsertAtCursor(text + "\n")
		linesRead++
	}
}

// if < 24 hours, x hours
// if < 7 days, x days
// ignores leap years, clock changes, etc. Just a rough conversion for display purposes.
func GetTimeAgo(jobTime time.Time) string {
	nowUnix := time.Now().Unix()
	jobTimeUnix := jobTime.Unix()
	secondsAgo := nowUnix - jobTimeUnix

	var oneMinuteInSeconds float64 = 60
	var oneHourInSeconds float64 = oneMinuteInSeconds * 60
	var oneDayInSeconds float64 = oneHourInSeconds * 24
	var oneWeekInSeconds float64 = oneDayInSeconds * 7
	var oneMonthInSeconds float64 = oneDayInSeconds * 31 // approx
	var oneYearInSeconds float64 = oneDayInSeconds * 365

	yearsAgo := int(float64(secondsAgo) / oneYearInSeconds)
	if yearsAgo >= 1 {
		if yearsAgo == 1 {
			return "1 year"
		} else {
			return strconv.Itoa(yearsAgo) + " years"
		}
	}

	monthsAgo := int(float64(secondsAgo) / oneMonthInSeconds)
	if monthsAgo >= 1 {
		if monthsAgo == 1 {
			return "1 month"
		} else {
			return strconv.Itoa(monthsAgo) + " months"
		}
	}

	weeksAgo := int(float64(secondsAgo) / oneWeekInSeconds)
	if weeksAgo >= 1 {
		if weeksAgo == 1 {
			return "1 week"
		} else {
			return strconv.Itoa(weeksAgo) + " weeks"
		}
	}

	daysAgo := int(float64(secondsAgo) / oneDayInSeconds)
	if daysAgo >= 1 {
		if daysAgo == 1 {
			return "1 day"
		} else {
			return strconv.Itoa(daysAgo) + " days"
		}
	}

	hoursAgo := int(float64(secondsAgo) / oneHourInSeconds)
	if hoursAgo >= 1 {
		if hoursAgo == 1 {
			return "1 hour"
		} else {
			return strconv.Itoa(hoursAgo) + " hours"
		}
	}

	minutesAgo := int(float64(secondsAgo) / oneMinuteInSeconds)
	if minutesAgo >= 1 {
		if minutesAgo == 1 {
			return "1 minute"
		} else {
			return strconv.Itoa(minutesAgo) + " minutes"
		}
	}

	if secondsAgo == 1 {
		return "1 second"
	}
	return strconv.FormatInt(secondsAgo, 10) + " seconds"

}

func formatDuration(duration time.Duration) string {
	var fragments []string
	if 0 != duration.Hours() {
		fragments = append(fragments, fmt.Sprintf("%f hours", duration.Hours()))
	}

	if 0 != duration.Minutes() {
		fragments = append(fragments, fmt.Sprintf("%f minutes", duration.Minutes()))
	}

	if 0 != duration.Seconds() {
		fragments = append(fragments, fmt.Sprintf("%f second", duration.Minutes()))
	}

	if 0 == len(fragments) {
		return "less than 1 second"
	}

	return strings.Join(fragments, " ")
}
