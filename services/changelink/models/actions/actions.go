package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/yaml.v3"
)

type Actions []Action

type Action interface {
	ActionName() string
	ActionType() ActionType
	Perform(watcherName, filePath string, lines *TriggeredLines) error
}

type ActionType string

const (
	LOCAL_LINK ActionType = "local_link"
	WEBHOOK    ActionType = "webhook"
	JIRA       ActionType = "jira"
	SLACK      ActionType = "slack"
	LOG        ActionType = "log"
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
	Type                ActionType `json:"action_type" bson:"action_type" yaml:"type"`
	Name                string     `json:"name" bson:"name" yaml:"name"`
}

func (s LineRange) String() string {
	return fmt.Sprintf("L%d - L%d", s.StartLine, s.EndLine)
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

type partialAction struct {
	baseAction `json:",inline" bson:",inline" yaml:",inline"`
	node yaml.Node
}

func (actions *Actions) UnmarshalYAML(value *yaml.Node) error {
	var nodes []yaml.Node

	err := value.Decode(&nodes)
	if err != nil {
		return err
	}

	finalActions := make(Actions, len(nodes))

	for i, a := range nodes {
		/*
		This is optimized to not decode twice. Essentially it just searches for the Node with value "type" which is
		the key for the type field and uses the next node in the array. It's not as flexible in case the field name ever
		changes

		An alternate approach would be:

		// Outside of this method define:
		type partialAction struct {
			baseAction `json:",inline" bson:",inline" yaml:",inline"`
			node yaml.Node
		}

		Replace below with:
		action := &partialAction{}
		err = a.Decode(action)
		if err != nil {
			return err
		}

		switch action.Type {
			case SLACK:
				action := &Slack{baseAction: action.baseAction}
				err = a.Decode(action)
				if err != nil {
					return errors.New(fmt.Sprintf("Failed to decode slack action: %s", err))
				}
				finalActions[i] = action
			case LOG:
			...
		 */
		var actionType ActionType

		for i, n := range a.Content{
			if n.Value == "type" {
				actionType = ActionType(a.Content[i+1].Value)
				break
			}
		}

		switch actionType {
		case SLACK:
			action := &Slack{}
			err = a.Decode(action)
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to decode slack action: %s", err))
			}
			finalActions[i] = action
		case LOG:
			action := &Log{}
			err = a.Decode(action)
			if err != nil {
				return errors.New(fmt.Sprintf("Failed to decode slack action: %s", err))
			}
			finalActions[i] = action
		default:
			return errors.New(fmt.Sprintf("Unknown action type %s", actionType))
		}
	}

	*actions = finalActions
	return nil
}
