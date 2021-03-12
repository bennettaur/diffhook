package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Actions []Action

type Action interface {
	ActionName() string
	ActionType() ActionType
	Perform(name, filePath string, lines *TriggeredLines) error
}

type ActionType string

const (
	LOCAL_LINK    ActionType = "locallink"
	WEBHOOK ActionType = "webhook"
	JIRA    ActionType = "jira"
	SLACK   ActionType = "slack"
	LOG     ActionType = "log"
)

type TriggeredLines struct {
	DiffLines    LineRange
	WatchedLines LineRange
}

type LineRange struct {
	StartLine int
	EndLine   int
}

type baseAction struct {
	Type ActionType `json:"action_type" bson:"action_type"`
	Name string     `json:"name" bson:"name"`
	TriggerOnRename bool `json:"trigger_on_rename" bson:"trigger_on_rename"`
	TriggerOnMove bool `json:"trigger_on_move" bson:"trigger_on_move"`
	TriggerOnDelete bool `json:"trigger_on_delete" bson:"trigger_on_delete"`
	TriggerOnPermission bool `json:"trigger_on_permission" bson:"trigger_on_permission"`
}

func (s *baseAction) ActionType() ActionType {
	return s.Type
}

func (s *baseAction) ActionName() string {
	return s.Name
}

func (actions *Actions) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	rawData := bson.RawValue{Type: t, Value: data}
	err := rawData.Unmarshal(&actions)
	if err != nil {
		return err
	}

	var action struct {
		Actions []bson.RawValue
	}

	err = rawData.Unmarshal(&action)
	if err != nil {
		return err
	}

	for i, a := range *actions {
		switch a.ActionType() {
		case SLACK:
			slackAction := &Slack{}
			err = action.Actions[i].Unmarshal(slackAction)
			(*actions)[i] = slackAction
		case LOG:
			logAction := &Log{}
			err = action.Actions[i].Unmarshal(logAction)
			(*actions)[i] = logAction
		default:
			return errors.New(fmt.Sprintf("Unknown action type %s", a.ActionType()))
		}
	}
	return err
}

func (actions *Actions) UnmarshalJSON(data []byte) error {
	var stubActions struct {
		Actions []baseAction `json:"actions"`
	}

	err := json.Unmarshal(data, &stubActions)
	if err != nil {
		return err
	}

	decodeActions := make(Actions, len(stubActions.Actions))

	for i, a := range stubActions.Actions {
		switch a.Type {
		case SLACK:
			decodeActions[i] = &Slack{}
		case LOG:
			decodeActions[i] = &Log{}
		default:
			return errors.New(fmt.Sprintf("Unknown action type %s", a.Type))
		}
	}

	err = json.Unmarshal(data, &decodeActions)
	if err != nil {
		return err
	}
	*actions = decodeActions

	return nil
}

func (actions *Actions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stubActions []map[string]interface{}

	err := unmarshal(&stubActions)
	if err != nil {
		return err
	}

	finalActions := make(Actions, len(stubActions))

	for i, a := range stubActions {
		switch ActionType(a["type"].(string)) {
		case SLACK:
			action := &Slack{}
			err = mapstructure.Decode(a, action)
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to decode slack action: %s", err))
			}
			finalActions[i] = action
		case LOG:
			action := &Log{}
			err = mapstructure.Decode(a, action)
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to decode slack action: %s", err))
			}
			finalActions[i] = action
		default:
			return errors.New(fmt.Sprintf("Unknown action type %s", a["type"]))
		}
	}

	*actions = finalActions
	return nil
}
