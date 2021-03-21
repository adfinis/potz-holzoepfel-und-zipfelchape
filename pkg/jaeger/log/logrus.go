package log

/*
* logrus adapter for jaeger-client
 */

import log "github.com/sirupsen/logrus"

type logrusAdapter struct {
	logger *log.Logger
}

func NewLogrusAdapter(logger *log.Logger) *logrusAdapter {
	return &logrusAdapter{
		logger: logger,
	}
}

func (l *logrusAdapter) Error(msg string) {
	l.logger.Errorf(msg)
}

func (l *logrusAdapter) Infof(msg string, args ...interface{}) {
	l.logger.Infof(msg, args...)
}

func (l *logrusAdapter) Debugf(msg string, args ...interface{}) {
	l.logger.Debugf(msg, args...)
}
