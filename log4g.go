package log4g

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var LoggerKey = "LOGGER_INTERFACE"

// A variable so that log4g functions can be directly accessed to log library
var reflog LoggerLibrary

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// FieldMap allows customization of the key names for default fields.
type FieldMap map[string]string

// LogArr is array of string.
type LogArr []string

//LoggerLibrary is our contract for the logger-library
type LoggerLibrary interface {
	Msg(msg string)
	WithFields(keyValues Fields) LoggerLibrary
}

//Logger is our contract for the logger-interface
type Logger interface {
	Msg(msg string)

	Debug(format string, args ...interface{})

	Info(format string, args ...interface{})

	Error(format string, args ...interface{})

	WithFields(keyValues Fields) Logger

	GetSeverity() string
}

type Map interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

// Configuration stores the config for the logger-interface
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	JSONFormat      bool
	LogLevel        string
	FieldMap        FieldMap
	TimestampFormat string
}

func Get(c Map) Logger {
	data, _ := c.Get(LoggerKey)
	return data.(Logger)
}

//NewContextLogger returns an instance of context with configuration for logger
func NewContextLogger(ctx context.Context, config Configuration, loggerInstance int) (Logger, error) {
	switch loggerInstance {
	case InstanceLogrusLogger:
		log, err := newLogrusLogger(config)
		if err != nil {
			return nil, err
		}
		reflog = log
	default:
		return nil, errors.New("invalid logger-interface instance")
	}
	logger := &ContextLogger{
		level:  parseLevel(config.LogLevel),
		ctx:    ctx,
		reflog: reflog,
	}
	mapCtx, ok := interface{}(ctx).(Map)
	if ok && mapCtx != nil {
		mapCtx.Set(LoggerKey, logger)
	}
	return logger, nil
}

type ContextLogger struct {
	level  Level
	ctx    context.Context
	reflog LoggerLibrary
}

func (_this *ContextLogger) Error(format string, args ...interface{}) {
	_this.Append(ErrKey, format, args...)
}

func (_this *ContextLogger) Debug(format string, args ...interface{}) {
	_this.Append(DebugKey, format, args...)
}

func (_this *ContextLogger) Info(format string, args ...interface{}) {
	_this.Append(InfoKey, format, args...)
}

// WithFields WithFields
func (_this *ContextLogger) WithFields(keyValues Fields) Logger {
	_this.reflog = _this.reflog.WithFields(keyValues)
	return _this
}

func (_this *ContextLogger) Set(key string, value interface{}) {
	if _this.ctx != nil {
		_this.ctx = context.WithValue(_this.ctx, key, value)
	}
}

func (_this *ContextLogger) Get(key string) interface{} {
	if _this.ctx != nil {
		if v := _this.ctx.Value(key); v != nil {
			return v
		}
	}
	return nil
}

func (_this *ContextLogger) GetLogData(key string) LogArr {
	val := _this.Get(key)
	if val != nil {
		return val.(LogArr)
	}
	return nil
}

func (_this *ContextLogger) Initial(key string) {
	if _this.ctx != nil {
		if v := _this.ctx.Value(key); v != nil {
			return
		}
		_this.ctx = context.WithValue(_this.ctx, key, LogArr{})
	}
}

func (_this *ContextLogger) Append(key string, format string, args ...interface{}) {
	_this.Initial(key)
	if value := _this.Get(key); value != nil {
		logArr := value.(LogArr)
		logArr = append(logArr, fmt.Sprintf(format, args...))
		_this.Set(key, logArr)
	}
}

func (_this *ContextLogger) Msg(msg string) {
	for _, key := range []string{ErrKey, InfoKey, DebugKey} {
		_this.reflog = _this.prepareLogData(key, _this.level)
	}
	_this.reflog.Msg(msg)
}

func parseLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "error":
		return ErrorLevel
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	}
	return InfoLevel
}

func (_this *ContextLogger) prepareLogData(key string, currentLogLevel Level) LoggerLibrary {
	logData := _this.GetLogData(key)
	if (key == ErrKey && ErrorLevel >= currentLogLevel) || (key == InfoKey && InfoLevel >= currentLogLevel) || (key == DebugKey && DebugLevel >= currentLogLevel) {
		for i, data := range logData {
			_this.reflog = _this.reflog.WithFields(Fields{
				fmt.Sprintf("%v_%v", key, i): data,
			})
		}
	}
	return _this.reflog
}

func (_this *ContextLogger) GetSeverity() string {
	var currentLogLevel = _this.level
	if logData := _this.GetLogData(DebugKey); len(logData) != 0 && DebugLevel >= _this.level {
		currentLogLevel = DebugLevel
	}

	if logData := _this.GetLogData(InfoKey); len(logData) != 0 && InfoLevel >= _this.level {
		currentLogLevel = InfoLevel
	}

	if logData := _this.GetLogData(WarningKey); len(logData) != 0 && WarnLevel >= _this.level {
		currentLogLevel = WarnLevel
	}

	if logData := _this.GetLogData(ErrKey); len(logData) != 0 && ErrorLevel >= _this.level {
		currentLogLevel = ErrorLevel
	}

	switch currentLogLevel {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	}
	return "DEFAULT"
}
