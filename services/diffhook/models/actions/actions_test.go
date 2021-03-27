package actions

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/yaml.v3"
	"testing"
)

var marshalTests = []struct {
	name    string
	actions Actions
}{
	{
		name: "empty actions",
		actions: Actions{},
	},
	{
		name: "one action",
		actions: Actions{
			NewLogAction("log", "Test"),
		},
	},
	{
		name: "multiple of the same action",
		actions: Actions{
			NewLogAction("log", "Test"),
			NewLogAction("log2", "Test2"),
			NewLogAction("log3", "Test3"),
		},
	},
	{
		name: "mix of all actions, all fields",
		actions: Actions{
			NewLogAction("log", "Test"),
			NewSlackAction("slack", "#some-channel", "Slack"),
		},
	},
	{
		name: "mix of all actions, missing some fields",
		actions: Actions{
			NewLogAction("log", "Test"),
			NewSlackAction("slack", "#some-channel", ""),
		},
	},
}

func TestActions_UnmarshalBSON(t *testing.T) {
	for _, tt := range marshalTests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := bson.Marshal(tt.actions)
			require.Nil(t, err, "BSON Marshaling Actions bson failed: %s", err)

			var newActions Actions

			err = bson.Unmarshal(data, &newActions)
			require.Nil(t, err, "BSON Unmarshaling Actions bson failed: %s", err)
			assert.Equal(t, tt.actions, newActions)
		})
	}
}

func TestActions_UnmarshalJSON(t *testing.T) {
	for _, tt := range marshalTests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.actions)
			require.Nil(t, err, "JSON Marshaling Actions bson failed: %s", err)

			var newActions Actions

			err = json.Unmarshal(data, &newActions)
			require.Nil(t, err, "JSON Unmarshaling Actions bson failed: %s", err)
			assert.Equal(t, tt.actions, newActions)
		})
	}
}

func TestActions_UnmarshalYAML(t *testing.T) {
	for _, tt := range marshalTests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := yaml.Marshal(tt.actions)
			require.Nil(t, err, "YAML Marshaling Actions bson failed: %s", err)

			var newActions Actions

			err = yaml.Unmarshal(data, &newActions)
			require.Nil(t, err, "YAML Unmarshaling Actions bson failed: %s", err)
			assert.Equal(t, tt.actions, newActions)
		})
	}
}

func Test_baseAction_ActionName(t *testing.T) {
	tests := []struct {
		name   string
		action Action
		want   string
	}{
		{
			name: "slack action returns the name",
			action: NewSlackAction("slackName", "channel", "msg"),
			want: "slackName",
		},
		{
			name: "log action returns the name",
			action: NewLogAction("logName", "msg"),
			want: "logName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.ActionName(); got != tt.want {
				t.Errorf("ActionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_baseAction_ActionType(t *testing.T) {
	tests := []struct {
		name   string
		action Action
		want   ActionType
	}{
		{
			name: "slack action returns the name",
			action: NewSlackAction("slackName", "channel", "msg"),
			want: SLACK,
		},
		{
			name: "log action returns the name",
			action: NewLogAction("logName", "msg"),
			want: LOG,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.action.ActionType(); got != tt.want {
				t.Errorf("ActionType() = %v, want %v", got, tt.want)
			}
		})
	}
}
