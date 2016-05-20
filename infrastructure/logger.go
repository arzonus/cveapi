package infrastructure

import (
	log "github.com/Sirupsen/logrus"
)

type Logger struct{}

func (logger Logger) Log(args ...interface{}) {
	log.Info(args)
}

func (logger Logger) Error(args ...interface{}) {
	log.Error(args)
}

func (logger Logger) Fatal(args ...interface{}) {
	log.Fatal(args)
}

func (logger Logger) Debug(args ...interface{}) {
	log.Debug(args)
}

func (logger Logger) Info(args ...interface{}) {
	log.Info(args)
}
