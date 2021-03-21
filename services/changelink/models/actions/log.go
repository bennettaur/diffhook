package actions

import (
	"fmt"
)

type Log struct {
	baseAction `json:",inline" bson:",inline" yaml:",inline"`
	Message string `json:"message" bson:"message" yaml:"message"`
}

func NewLogAction(name, message string) Action {
	return &Log{
		baseAction: baseAction{
			Name: name,
			Type: LOG,
		},
		Message: message,
	}
}

func (s *Log) Perform(watcherName, filePath, reason string, lines *TriggeredLines) error {
	fmt.Printf("I logged message %s\n", s.Message)
	return nil
}
