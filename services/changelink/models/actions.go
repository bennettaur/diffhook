package models

type ActionType string

const (
	LINK    ActionType = "changelink"
	WEBHOOK            = "webhook"
	JIRA               = "jira"
	SLACK              = "slack"
)

type Action struct {
	ActionType ActionType `json:"action_type" bson:"action_type"`
	Message    string     `json:"message" bson:"message"`
}

