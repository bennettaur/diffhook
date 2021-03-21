package models

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

const DefaultFileStore = ".changelink.yml"

var configuredStore = DefaultFileStore

func SetLocalStore(filePath string) {
	configuredStore = filePath
}

type Store interface {
	Find ()
}

type LocalStore struct {
	filePath string
	Watchers []Watcher `json:"watchers"`
}

func GetLocalStore(filePath string) (*LocalStore, error) {
	l := &LocalStore{filePath: configuredStore}
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
