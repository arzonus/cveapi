package infrastructure

import (
	log "github.com/Sirupsen/logrus"
)

type Logger struct{}

func (logger Logger) Log(args ...interface{}) {
	log.Info("CVA: ", args)
}

func (logger Logger) Error(args ...interface{}) {
	log.Error("CVA: ", args)
}

func (logger Logger) Fatal(args ...interface{}) {
	log.Fatal("CVA: ", args)
}

func (logger Logger) Debug(args ...interface{}) {
	log.Debug("CVA: ", args)
}

func (logger Logger) Info(args ...interface{}) {
	log.Info("CVA: ", args)
}
