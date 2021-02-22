package logging

import (
	"fmt"

	"go.uber.org/zap"
)

// ZapLogger is a zap sugared logger wrapper.
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger returns a new zap logger.
func NewZapLogger() (*ZapLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("could not initialise a zap logger")
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	return &ZapLogger{logger: sugar}, nil
}

func (l *ZapLogger) Errorf(msg string, args ...interface{}) {
	l.Error(fmt.Sprintf(msg, args...))
}

func (l *ZapLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *ZapLogger) Infof(msg string, args ...interface{}) {
	l.logger.Infof(fmt.Sprintf(msg, args...))
}
