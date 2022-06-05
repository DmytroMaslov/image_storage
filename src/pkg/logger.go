package pkg

import "go.uber.org/zap"

type logger struct {
	logger *zap.SugaredLogger
}

type Logger interface {
	InitLogger()
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

func NewLogger() Logger {
	return &logger{}
}
func (l *logger) InitLogger() {
	loggerInstance, _ := zap.NewDevelopment()
	l.logger = loggerInstance.Sugar()
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
	l.logger.Sync()
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
	l.logger.Sync()

}
func (l *logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
	l.logger.Sync()
}
