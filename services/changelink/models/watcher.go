package models

import (
	"errors"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

const UNBOUNDED = -1

type LineRange struct {
	StartLine int
	EndLine   int
}

type Watcher struct {
	// DefaultModel add _id,created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Name             string      `json:"name" bson:"name"`
	Host             string      `json:"host" bson:"host"`
	FilePath         string      `json:"file_path" bson:"file_path"`
	Lines            []LineRange `json:"lines" bson:"lines"`
	Actions          []Action    `json:"actions" bson:"actions"`
}

func NewWatcher(name, host, filePath string, lines []LineRange) *Watcher {
	return &Watcher{
		Name:     name,
		Host:     host,
		FilePath: filePath,
		Lines:    lines,
	}
}

func (w *Watcher) AddAction(a Action) {
	w.Actions = append(w.Actions, a)
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
		validationErrors = append(validationErrors,  errors.New("missing file path"))
	}

	if w.FilePath == "" {
		validationErrors = append(validationErrors,  errors.New("missing file path"))
	}
	return nil
}

func findWatcherForFileLocal(filePath string) ([]Watcher, error) {
	s, err := GetLocalStore()
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