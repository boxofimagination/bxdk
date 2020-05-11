package logger

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	// KV is a type for logging with more information
	// this used by with function
	KV map[string]interface{}

	// Logger interface
	Logger interface {
		SetLevel(level Level)
		Debug(args ...interface{})
		Debugln(args ...interface{})
		Debugf(format string, args ...interface{})
		DebugWithFields(msg string, KV KV)
		Info(args ...interface{})
		Infoln(args ...interface{})
		Infof(format string, args ...interface{})
		InfoWithFields(msg string, KV KV)
		Warn(args ...interface{})
		Warnln(args ...interface{})
		Warnf(format string, args ...interface{})
		WarnWithFields(msg string, KV KV)
		Error(args ...interface{})
		Errorln(args ...interface{})
		Errorf(format string, args ...interface{})
		ErrorWithFields(msg string, KV KV)
		Errors(err error)
		Fatal(args ...interface{})
		Fatalln(args ...interface{})
		Fatalf(format string, args ...interface{})
		FatalWithFields(msg string, KV KV)
		IsValid() bool // IsValid check if Logger is created using constructor
	}

	// Level of log
	Level int

	// Engine of logger
	Engine string
)

// list of log level
const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

// DefaultTimeFormat of logger
const DefaultTimeFormat = time.RFC3339

// Logger engine option
const (
	Logrus  Engine = "logrus"
	Zerolog Engine = "zerolog"
)

// StringToLevel to set string to level
func StringToLevel(level string) Level {
	switch strings.ToLower(level) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		// TODO: make this more informative when happened
		return InfoLevel
	}
}

// LevelToString convert log level to readable string
func LevelToString(l Level) string {
	switch l {
	case DebugLevel:
		return DebugLevelString
	case InfoLevel:
		return InfoLevelString
	case WarnLevel:
		return WarnLevelString
	case ErrorLevel:
		return ErrorLevelString
	case FatalLevel:
		return FatalLevelString
	default:
		return InfoLevelString
	}
}

// Config of logger
type Config struct {
	Level      Level
	LogFile    string
	TimeFormat string
	CallerSkip int
	Caller     bool
	UseColor   bool
	UseJSON    bool
}

// OpenLogFile tries to open the log file (creates it if not exists) in write-only/append mode and return it
// Note: the func return nil for both *os.File and error if the file name is empty string
func (c *Config) OpenLogFile() (*os.File, error) {
	if c.LogFile == "" {
		return nil, nil
	}

	err := os.MkdirAll(filepath.Dir(c.LogFile), 0755)
	if err != nil && err != os.ErrExist {
		return nil, err
	}

	return os.OpenFile(c.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}
