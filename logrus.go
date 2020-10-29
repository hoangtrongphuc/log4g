package log4g

import (
	"github.com/sirupsen/logrus"
	"os"
)

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

func getFormatter(config Configuration) logrus.Formatter {
	if config.JSONFormat {
		return &logrus.JSONFormatter{
			TimestampFormat: config.TimestampFormat,
			FieldMap:        convertToLogrusFieldMap(config.FieldMap),
		}
	}
	return &logrus.TextFormatter{
		TimestampFormat: config.TimestampFormat,
		FieldMap:        convertToLogrusFieldMap(config.FieldMap),
	}
}

func newLogrusLogger(config Configuration) (LoggerLibrary, error) {
	logLevel := config.LogLevel

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}

	stdOutHandler := os.Stdout
	lLogger := &logrus.Logger{
		Out:       stdOutHandler,
		Formatter: getFormatter(config),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	return &logrusLogger{
		logger: lLogger,
	}, nil
}

func (l *logrusLogger) Msg(msg string) {
	l.logger.Info(msg)
}

func (l *logrusLogger) WithFields(fields Fields) LoggerLibrary {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) Msg(msg string) {
	l.entry.Info(msg)
}

func (l *logrusLogEntry) WithFields(fields Fields) LoggerLibrary {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}

func convertToLogrusFieldMap(fieldMap FieldMap) logrus.FieldMap {
	logrusFieldMap := logrus.FieldMap{}
	for key, val := range fieldMap {
		switch key {
		case FieldKeyLevel:
			logrusFieldMap[logrus.FieldKeyLevel] = val
		case FieldKeyTime:
			logrusFieldMap[logrus.FieldKeyTime] = val
		case FieldKeyMsg:
			logrusFieldMap[logrus.FieldKeyMsg] = val
		case FieldKeyFunc:
			logrusFieldMap[logrus.FieldKeyFunc] = val
		}
	}
	return logrusFieldMap
}
