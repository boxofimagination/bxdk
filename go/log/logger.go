package log

import (
	"errors"

	"github.com/boxofimagination/bxdk/go/env"
	"github.com/boxofimagination/bxdk/go/log/logger"
	"github.com/boxofimagination/bxdk/go/log/logger/zerolog"
)

type (
	// Level of log
	Level = logger.Level

	// Logger interface
	Logger = logger.Logger
)

// Zerolog engine
const Zerolog Engine = logger.Zerolog

// Level option
const (
	DebugLevel = logger.DebugLevel
	InfoLevel  = logger.InfoLevel
	WarnLevel  = logger.WarnLevel
	ErrorLevel = logger.ErrorLevel
	FatalLevel = logger.FatalLevel
)

var (
	isDev          = env.IsDevelopment()
	debugLogger, _ = NewLogger(Zerolog, &logger.Config{Level: logger.DebugLevel, UseColor: isDev, UseJSON: !isDev})
	infoLogger, _  = NewLogger(Zerolog, &logger.Config{Level: logger.InfoLevel, UseColor: isDev, UseJSON: !isDev})
	warnLogger     = infoLogger
	errorLogger    = infoLogger
	fatalLogger    = infoLogger
	loggers        = [5]*Logger{
		&debugLogger,
		&infoLogger,
		&warnLogger,
		&errorLogger,
		&fatalLogger,
	}

	errInvalidLogger = errors.New("invalid logger")
	errInvalidLevel  = errors.New("invalid log level")
)

// NewLogger creates a new zerolog logger.
// Engine parameter is not used anymore.
func NewLogger(engine Engine, config *logger.Config) (Logger, error) {
	config.UseColor = env.IsDevelopment()
	config.UseJSON = !config.UseColor
	return zerolog.New(config)
}

// SetLogger for certain level
func SetLogger(level logger.Level, lgr logger.Logger) error {
	if level < logger.DebugLevel || level > logger.FatalLevel {
		return errInvalidLevel
	}
	if lgr == nil || !lgr.IsValid() {
		return errInvalidLogger
	}
	*loggers[level] = lgr
	return nil
}

// SetLevel adjusts log level threshold.
// Only log with level higher or equal with this level will be printed
func SetLevel(level Level) {
	if level < 0 {
		level = InfoLevel
	}

	debugLogger.SetLevel(level)
	infoLogger.SetLevel(level)
	warnLogger.SetLevel(level)
	errorLogger.SetLevel(level)
	fatalLogger.SetLevel(level)
}

// SetLevelString adjusts log level threshold using string
func SetLevelString(level string) {
	SetLevel(logger.StringToLevel(level))
}
