package models

type ActionType string

const (
	LINK    ActionType = "changelink"
	WEBHOOK ActionType = "webhook"
	JIRA    ActionType = "jira"
	SLACK   ActionType = "slack"
	LOG     ActionType = "log"
)

type Action struct {
	ActionType ActionType `json:"action_type" bson:"action_type"`
	Message    string     `json:"message" bson:"message"`
}
