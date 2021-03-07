package models

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Store interface {
	Find ()
}

type LocalStore struct {
	Watchers []Watcher
}

func GetLocalStore() (*LocalStore, error) {
	l := &LocalStore{}
	data, err := ioutil.ReadFile("./data/watchers.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, l)
	if err != nil {
		return nil, err
	}
	return l, nil
}
