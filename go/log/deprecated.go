package log

import (
	"github.com/boxofimagination/bxdk/go/log/logger"
)

// Engine of logger.
// Deprecated, there is only 1 engine now
type Engine = logger.Engine

// Logrus engine.
// Deprecated, logrus engine is dropped, use zerolog engine instead
const Logrus Engine = logger.Logrus

// Debugw prints debug level log with additional fields.
// Deprecated: use DebugWithFields
func Debugw(msg string, keyValues KV) {
	debugLogger.DebugWithFields(msg, keyValues)
}

// Infow prints info level log with additional fields.
// Deprecated: use InfoWithFields
func Infow(msg string, keyValues KV) {
	infoLogger.InfoWithFields(msg, keyValues)
}

// Warnw prints warn level log with additional fields.
// Deprecated: use WarnWithFields
func Warnw(msg string, keyValues KV) {
	warnLogger.WarnWithFields(msg, keyValues)
}

// Errorw prints error level log with additional fields.
// Deprecated: use ErrorWithFields
func Errorw(msg string, keyValues KV) {
	errorLogger.ErrorWithFields(msg, keyValues)
}

// Fatalw prints fatal level log with additional fields.
// Deprecated: use FatalWithFields
func Fatalw(msg string, keyValues KV) {
	fatalLogger.FatalWithFields(msg, keyValues)
}

// Config of log
// Deprecated
type Config struct {
	Level      string
	TimeFormat string
	// Caller, option to print caller line numbers.
	// make sure you understand the overhead when use this
	Caller bool

	// LogFile for log to file.
	// this is not needed by default,
	// application is expected to run in containerized environment
	LogFile string
	// DebugFile for debug to file.
	// this is not needed by default,
	// application is expected to run in containerized environment
	DebugFile string

	// Deprecated, this field will have no effect
	Engine Engine
	// UseColor, option to colorize log in console.
	// Deprecated true if and only if TKPENV=development
	UseColor bool
	// UseJSON, option to print in json format.
	// Deprecated, true if and only if TKPENV!=development
	UseJSON bool
}

// SetConfig creates new logger based on given config
// Deprecated: use Init
func SetConfig(config *Config) error {
	var (
		newDebugLogger    logger.Logger
		newLogger         logger.Logger
		err               error
		debugLoggerConfig = logger.Config{Level: logger.DebugLevel}
		loggerConfig      = logger.Config{Level: logger.InfoLevel}
		engine            = Zerolog
	)

	if config != nil {
		engine = config.Engine

		loggerConfig = logger.Config{
			Level:      logger.StringToLevel(config.Level),
			LogFile:    config.LogFile,
			TimeFormat: config.TimeFormat,
			Caller:     config.Caller,
		}

		// copy
		debugLoggerConfig = loggerConfig

		// custom output file
		debugLoggerConfig.LogFile = config.DebugFile
	}

	loggerConfig.UseColor = isDev
	debugLoggerConfig.UseColor = isDev
	loggerConfig.UseJSON = !isDev
	debugLoggerConfig.UseJSON = !isDev

	newLogger, err = NewLogger(engine, &loggerConfig)
	if err != nil {
		return err
	}
	// extra check because it is very difficult to debug if the log itself causes the panic
	if newLogger != nil {
		infoLogger = newLogger
		warnLogger = newLogger
		errorLogger = newLogger
		fatalLogger = newLogger
	}

	newDebugLogger, err = NewLogger(engine, &debugLoggerConfig)
	if err != nil {
		return err
	}
	if newDebugLogger != nil {
		debugLogger = newDebugLogger
	}

	return nil
}
