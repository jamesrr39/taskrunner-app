package gui

import (
	"bufio"
	"io"
	"strconv"
	"time"

	"github.com/mattn/go-gtk/gtk"
)

func fillTextBufferFromFile(textBuffer *gtk.TextBuffer, fileReader io.Reader, linesToRead uint) {
	fileScanner := bufio.NewScanner(fileReader)
	linesRead := uint(0)
	for fileScanner.Scan() && linesRead < linesToRead {
		textBuffer.InsertAtCursor(fileScanner.Text())
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
			return "1 year ago"
		} else {
			return strconv.Itoa(yearsAgo) + " years ago"
		}
	}

	monthsAgo := int(float64(secondsAgo) / oneMonthInSeconds)
	if monthsAgo >= 1 {
		if monthsAgo == 1 {
			return "1 month ago"
		} else {
			return strconv.Itoa(monthsAgo) + " months ago"
		}
	}

	weeksAgo := int(float64(secondsAgo) / oneWeekInSeconds)
	if weeksAgo >= 1 {
		if weeksAgo == 1 {
			return "1 week ago"
		} else {
			return strconv.Itoa(weeksAgo) + " weeks ago"
		}
	}

	daysAgo := int(float64(secondsAgo) / oneDayInSeconds)
	if daysAgo >= 1 {
		if daysAgo == 1 {
			return "1 day ago"
		} else {
			return strconv.Itoa(daysAgo) + " days ago"
		}
	}

	hoursAgo := int(float64(secondsAgo) / oneHourInSeconds)
	if hoursAgo >= 1 {
		if hoursAgo == 1 {
			return "1 hour ago"
		} else {
			return strconv.Itoa(hoursAgo) + " hours ago"
		}
	}

	minutesAgo := int(float64(secondsAgo) / oneMinuteInSeconds)
	if minutesAgo >= 1 {
		if minutesAgo == 1 {
			return "1 minute ago"
		} else {
			return strconv.Itoa(minutesAgo) + " minutes ago"
		}
	}

	if secondsAgo == 1 {
		return "1 second ago"
	}
	return strconv.FormatInt(secondsAgo, 10) + " seconds ago"

}
