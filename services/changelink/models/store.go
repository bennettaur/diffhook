package models

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const DEFAULT_FILE_STORE = ".changelink.yml"

type Store interface {
	Find ()
}

type LocalStore struct {
	filePath string
	Watchers []Watcher `json:"watchers"`
}

func GetLocalStore(filePath string) (*LocalStore, error) {
	if filePath == "" {
		filePath = DEFAULT_FILE_STORE
	}
	l := &LocalStore{filePath: filePath}
	data, err := ioutil.ReadFile(l.filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *LocalStore) Save() error {
	data, err := yaml.Marshal(l)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(l.filePath, data, 0644)
	return err
}
