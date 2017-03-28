package taskrunner

import (
	"errors"
)

type Script string

type Job struct {
	Id          uint   `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Script      Script `json:"script"`
}

func NewJob(id uint, name string, description string, script Script) (*Job, error) {
	if "" == name {
		return nil, errors.New("A job must have a name")
	}

	return &Job{Id: id, Name: name, Description: description, Script: script}, nil
}
