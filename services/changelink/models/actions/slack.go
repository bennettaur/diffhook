package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"log"
)

type Slack struct {
	baseAction `json:",inline" bson:",inline" yaml:",inline"`
	Channel    string `json:"channel" bson:"channel" yaml:"channel"`
	Message    string `json:"message" bson:"message" yaml:"message"`
}

func NewSlackAction(name, channel, message string) Action {
	return &Slack{
		baseAction: baseAction{
			Name: name,
			Type: SLACK,
		},
		Channel: channel,
		Message: message,
	}
}

func (s *Slack) Perform(watcherName, filePath, reason string, lines *TriggeredLines) error {
	channelId, err := findChannelId(s.Channel)
	if err != nil {
		return err
	}

	api, err := getSlackClient()
	if err != nil {
		return err
	}

	header := &slack.TextBlockObject{
		Type: slack.PlainTextType,
		Text: fmt.Sprintf("%s: %s", watcherName, s.Name),
	}

	msgSection := &slack.TextBlockObject{
		Type: slack.MarkdownType,
		Text: s.Message,
	}

	var changeTrigger string
	if lines == nil {
		changeTrigger = reason
	} else {
		changeTrigger = fmt.Sprintf("Changed lines in %s:\n\n```\n%s\n```", filePath, lines.Hunk.Body)
	}

	codeSection := &slack.TextBlockObject{
		Type: slack.MarkdownType,
		Text: changeTrigger,
	}

	postBlocks := []slack.Block{
		slack.NewHeaderBlock(header),
		slack.NewSectionBlock(msgSection, nil, nil),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(codeSection, nil, nil),
	}

	blocks, _ := json.Marshal(postBlocks)
	log.Printf("Marshaled blocks:\n%s", blocks)

	_, _, err = api.PostMessage(channelId, slack.MsgOptionBlocks(postBlocks...))

	if err != nil {
		return err
	}

	fmt.Printf("I slacked message %s to channel %s\n", s.Message, s.Channel)
	return nil
}

func findChannelId(channelName string) (string, error) {
	api, err := getSlackClient()
	if err != nil {
		return "", err
	}

	params := &slack.GetConversationsParameters{
		ExcludeArchived: "true",
		Limit:           400,
		Types:           []string{"public_channel", "private_channel"},
	}
	for {
		channels, nextCursor, err := api.GetConversations(params)
		if err != nil {
			return "", err
		}
		for _, channel := range channels {
			if channel.Name == channelName {
				return channel.ID, nil
			}
		}
		if nextCursor == "" {
			return "", errors.New("channel not found, did you /invite @changelink to the channel?")
		}
		params.Cursor = nextCursor
	}
}

func getSlackClient() (*slack.Client, error) {
	slackToken := viper.GetString("SLACK_TOKEN")
	if slackToken == "" {
		return nil, errors.New("missing slack token. Is SLACK_TOKEN set?")
	}

	return slack.New(slackToken), nil
}
