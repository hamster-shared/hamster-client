package logger

import (
	"go.uber.org/zap"
)

var zapLog *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	zapLog, _ = zap.NewDevelopment()
	sugar = zapLog.Sugar()
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}
