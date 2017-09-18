//+build linux
package triggers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesrr39/taskrunner-app/taskrunner"
)

const (
	idVendorKey  = "ATTRS{idVendor}=="
	idProductKey = "ATTRS{idProduct}=="
	runKey       = "RUN+="
)

type UdevRulesDAL struct {
	BaseDir string // /etc/udev/rules.d
}

func NewUdevRulesDAL(baseDir string) *UdevRulesDAL {
	return &UdevRulesDAL{baseDir}
}

type UdevRule struct {
	IdVendor      string
	IdProduct     string
	Run           string
	FileDefinedIn string
}

func NewUdevRule(idVendor, idProduct, run, fileDefinedIn string) *UdevRule {
	return &UdevRule{idVendor, idProduct, run, fileDefinedIn}
}

func (u *UdevRulesDAL) GetRules(job *taskrunner.Job) ([]*UdevRule, error) {
	var rules []*UdevRule

	err := filepath.Walk(u.BaseDir, func(path string, fileInfo os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if nil != err {
			return err
		}
		defer file.Close()

		rules = append(rules, rulesFromFile(file, path, job)...)
		return nil
	})

	if nil != err {
		return nil, err
	}

	return rules, nil
}

func rulesFromFile(file io.Reader, filePath string, job *taskrunner.Job) []*UdevRule {
	var rules []*UdevRule

	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		line := strings.TrimSpace(fileScanner.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

		idVendor := getValueForProperty(line, idVendorKey)
		idProduct := getValueForProperty(line, idProductKey)
		runCommand := getValueForProperty(line, runKey)
		if "" == idVendor || "" == idProduct || !strings.Contains(runCommand, fmt.Sprintf("--run-job='%s'", job.Name)) {
			continue
		}

		rules = append(rules, NewUdevRule(idVendor, idProduct, runCommand, filePath))
	}

	return rules
}

func getValueForProperty(line string, key string) string {
	propertyKeyIndex := strings.Index(line, key)
	if -1 == propertyKeyIndex {
		return ""
	}
	propertyValueIndex := propertyKeyIndex + len(key) + strings.Index(line[propertyKeyIndex+len(key):], "\"") + 1 // chop off leading quotation

	propertyValueEndIndex := propertyValueIndex + strings.Index(line[propertyValueIndex:], "\"")
	if 0 > propertyValueEndIndex {
		return ""
	}

	return line[propertyValueIndex:propertyValueEndIndex]
}
