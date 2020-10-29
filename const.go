package log4g

import "time"

const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg            = "msg"
	FieldKeyLevel          = "level"
	FieldKeyTime           = "time"
	FieldKeyFunc           = "caller"
)

const (
	//Debug has verbose message
	DEBUG = "debug"
	//Info is default log4g level
	INFO = "info"
	//Warn is for logging messages about possible issues
	WARN = "warn"
	//Error is for logging errors
	ERROR = "error"
	//Fatal is for logging fatal messages. The system shutdown after logging the message.
	FATAL = "fatal"
)

// Logger key constants definition.
const (
	ErrKey        = "_err"
	InfoKey       = "_info"
	DebugKey      = "_debug"
	WarningKey    = "_warning"
	CustomDataKey = "_custom_data"
	XRequestID    = "x_request_id"
	AppName       = "app_name"
	Severity      = "severity"
)

// Level log4g constants definition.
type Level int8

const (
	// DebugLevel defines debug log4g level.
	DebugLevel Level = iota
	// InfoLevel defines info log4g level.
	InfoLevel
	// WarnLevel defines warn log4g level.
	WarnLevel
	// ErrorLevel defines error log4g level.
	ErrorLevel
	// FatalLevel defines fatal log4g level.
	FatalLevel
	// PanicLevel defines panic log4g level.
	PanicLevel
	// NoLevel defines an absent log4g level.
	NoLevel
	// Disabled disables the logger-interface.
	Disabled

	// TraceLevel defines trace log4g level.
	TraceLevel Level = -1
)

const (
	// InstanceZapLogger InstanceZapLogger
	InstanceZapLogger int = iota
	// InstanceLogrusLogger InstanceLogrusLogger
	InstanceLogrusLogger
	// InstanceGCPZapLogger InstanceGCPZapLogger
	InstanceGCPZapLogger
)
