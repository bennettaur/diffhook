package models

type Collection interface {
	New(args ...interface{}) (interface{}, error)
	Save() error
	Validate () error
}