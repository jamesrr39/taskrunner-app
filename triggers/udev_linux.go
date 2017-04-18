//+build linux
package triggers

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
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

func (u *UdevRulesDAL) GetRules() ([]*UdevRule, error) {
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

		rules = append(rules, rulesFromFile(file, path)...)
		return nil
	})

	if nil != err {
		return nil, err
	}

	return rules, nil
}

const (
	idVendorKey  = "ATTRS{idVendor}=="
	idProductKey = "ATTRS{idProduct}=="
	runKey       = "RUN="
)

func rulesFromFile(file io.Reader, filePath string) []*UdevRule {
	b := bufio.NewScanner(file)
	for b.Scan() {
		line := strings.TrimSpace(b.Text())
		if strings.HasPrefix(line, "#") {
			continue
		}

	}
}
