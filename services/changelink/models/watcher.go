package models

import (
	"errors"

	"github.com/bennettaur/changelink/services/changelink/models/actions"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

const UNBOUNDED = -1

var FULL_FILE = actions.LineRange{StartLine: UNBOUNDED, EndLine: UNBOUNDED}

// TODO: Watcher for new file added in a directory
type Watcher struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string              `json:"name" bson:"name" yaml:"name"`
	Host             string              `json:"host" bson:"host" yaml:"host"`
	FilePath         string              `json:"file_path" bson:"file_path" yaml:"file_path"`
	Lines            []actions.LineRange `json:"lines,omitempty" bson:"lines,omitempty" yaml:"lines,omitempty"`
	TriggerAny       bool                `json:"trigger_any" bson:"trigger_any" yaml:"trigger_any"`
	TriggerAnyLine   bool                `json:"trigger_any_line" bson:"trigger_any_line" yaml:"trigger_any_line"`
	TriggerOnRename  bool                `json:"trigger_on_rename" bson:"trigger_on_rename" yaml:"trigger_on_rename"`
	TriggerOnMove    bool                `json:"trigger_on_move" bson:"trigger_on_move" yaml:"trigger_on_move"`
	TriggerOnDelete  bool                `json:"trigger_on_delete" bson:"trigger_on_delete" yaml:"trigger_on_delete"`
	TriggerOnMode    bool                `json:"trigger_on_mode" bson:"trigger_on_mode" yaml:"trigger_on_mode"`
	Actions          *actions.Actions    `json:"actions" bson:"actions" yaml:"actions"`
}

func NewWatcher(name, host, filePath string, lines []actions.LineRange) *Watcher {
	return &Watcher{
		Name:     name,
		Host:     host,
		FilePath: filePath,
		Lines:    lines,
	}
}

func (w *Watcher) AddAction(a actions.Action) {
	*w.Actions = append(*w.Actions, a)
}

func FindWatchersForFile(filePath string) ([]Watcher, error) {
	return findWatcherForFileLocal(filePath)
}

func (w *Watcher) Validate() error {
	var validationErrors []error
	if w.Name == "" {
		validationErrors = append(validationErrors, errors.New("missing name"))
	}

	if w.FilePath == "" {
		validationErrors = append(validationErrors, errors.New("missing file path"))
	}

	if w.FilePath == "" {
		validationErrors = append(validationErrors, errors.New("missing file path"))
	}
	return nil
}

func findWatcherForFileLocal(filePath string) ([]Watcher, error) {
	s, err := GetLocalStore("")
	if err != nil {
		return nil, err
	}

	var result []Watcher

	for _, watcher := range s.Watchers {
		if watcher.FilePath == filePath {
			result = append(result, watcher)
		}
	}
	return result, nil
}

func findWatcherForFileMongo(filePath string) ([]Watcher, error) {
	var result []Watcher

	err := mgm.Coll(&Watcher{}).SimpleFind(&result, bson.M{"file_path": filePath})

	return result, err
}
