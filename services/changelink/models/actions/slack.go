package actions

import (
	"fmt"
)

type Slack struct {
	baseAction `json:",inline" bson:",inline" yaml:",inline"`
	Channel string `json:"channel" bson:"channel" yaml:"channel"`
	Message string `json:"message" bson:"message" yaml:"message"`
}

func NewSlackAction(name, channel, message string) Action {
	return &Slack{
		baseAction: baseAction{
			Name:    name,
			Type:    SLACK,
		},
		Channel: channel,
		Message: message,
	}
}

func (s *Slack) Perform(name, filePath string, lines *TriggeredLines) error {
	fmt.Printf("I slacked message %s to channel %s\n", s.Message, s.Channel)
	return nil
}
