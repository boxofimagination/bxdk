package zerolog

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/boxofimagination/bxdk/go/log/logger"
	"github.com/boxofimagination/bxdk/go/errors"
)

// Logger struct
type Logger struct {
	logger zerolog.Logger
	config logger.Config
	valid  bool
}

var (
	_ logger.Logger = (*Logger)(nil) // Make sure Logger struct complies with logger.Logger interface
)

// DefaultLogger return default value of logger
func DefaultLogger() *Logger {
	l := Logger{
		config: logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		},
		valid: true,
	}

	lgr := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		NoColor:    !l.config.UseColor,
		TimeFormat: l.config.TimeFormat,
	})
	lgr = setLevel(lgr, l.config.Level)

	l.logger = lgr
	return &l
}

// New logger
func New(config *logger.Config, opts ...func(*logger.Config)) (*Logger, error) {
	if config == nil {
		config = &logger.Config{
			Level:      logger.InfoLevel,
			TimeFormat: logger.DefaultTimeFormat,
		}
	}

	if config.TimeFormat == "" {
		config.TimeFormat = logger.DefaultTimeFormat
	}
	for _, opt := range opts {
		opt(config)
	}

	lgr, err := newLogger(config)
	if err != nil {
		return nil, err
	}
	l := Logger{
		logger: lgr,
		config: *config,
		valid:  true,
	}
	return &l, nil
}

func newLogger(config *logger.Config) (zerolog.Logger, error) {
	var (
		lgr zerolog.Logger
	)

	zerolog.TimeFieldFormat = config.TimeFormat
	zerolog.CallerSkipFrameCount = 4 + config.CallerSkip

	var writers zerolog.LevelWriter
	if config.UseJSON {
		writers = zerolog.MultiLevelWriter(os.Stderr)
	} else {
		writers = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			NoColor:    !config.UseColor,
			TimeFormat: config.TimeFormat,
		})
	}

	file, err := config.OpenLogFile()
	if err != nil {
		return lgr, err
	} else if file != nil {
		writers = zerolog.MultiLevelWriter(writers, file)
	}

	lgr = zerolog.New(writers)
	lgr = setLevel(lgr, config.Level)
	if config.Caller {
		lgr = lgr.With().Caller().Logger()
	}

	return lgr, nil
}

func setLevel(lgr zerolog.Logger, level logger.Level) zerolog.Logger {
	switch level {
	case logger.DebugLevel:
		lgr = lgr.Level(zerolog.DebugLevel)
	case logger.InfoLevel:
		lgr = lgr.Level(zerolog.InfoLevel)
	case logger.WarnLevel:
		lgr = lgr.Level(zerolog.WarnLevel)
	case logger.ErrorLevel:
		lgr = lgr.Level(zerolog.ErrorLevel)
	case logger.FatalLevel:
		lgr = lgr.Level(zerolog.FatalLevel)
	default:
		lgr = lgr.Level(zerolog.InfoLevel)
	}
	return lgr
}

// SetLevel for setting log level
func (l *Logger) SetLevel(level logger.Level) {
	if level < logger.DebugLevel || level > logger.FatalLevel {
		level = logger.InfoLevel
	}
	if level != l.config.Level {
		l.logger = setLevel(l.logger, level)
		l.config.Level = level
	}
}

// IsValid check if Logger is created using constructor
func (l *Logger) IsValid() bool {
	return l.valid
}

// Debug function
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug().Timestamp().Msg(fmt.Sprint(args...))
}

// Debugln function
func (l *Logger) Debugln(args ...interface{}) {
	l.logger.Debug().Timestamp().Msg(fmt.Sprintln(args...))
}

// Debugf function
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debug().Timestamp().Msgf(format, v...)
}

// DebugWithFields function
func (l *Logger) DebugWithFields(msg string, KV logger.KV) {
	l.logger.Debug().Timestamp().Fields(KV).Msg(msg)
}

// Info function
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info().Timestamp().Msg(fmt.Sprint(args...))
}

// Infoln function
func (l *Logger) Infoln(args ...interface{}) {
	l.logger.Info().Timestamp().Msg(fmt.Sprintln(args...))
}

// Infof function
func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Info().Timestamp().Msgf(format, v...)
}

// InfoWithFields function
func (l *Logger) InfoWithFields(msg string, KV logger.KV) {
	l.logger.Info().Timestamp().Fields(KV).Msg(msg)
}

// Warn function
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn().Timestamp().Msg(fmt.Sprint(args...))
}

// Warnln function
func (l *Logger) Warnln(args ...interface{}) {
	l.logger.Warn().Timestamp().Msg(fmt.Sprintln(args...))
}

// Warnf function
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Warn().Timestamp().Msgf(format, v...)
}

// WarnWithFields function
func (l *Logger) WarnWithFields(msg string, KV logger.KV) {
	l.logger.Warn().Timestamp().Fields(KV).Msg(msg)
}

// Error function
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error().Timestamp().Msg(fmt.Sprint(args...))
}

// Errorln function
func (l *Logger) Errorln(args ...interface{}) {
	l.logger.Error().Timestamp().Msg(fmt.Sprintln(args...))
}

// Errorf function
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Error().Timestamp().Msgf(format, v...)
}

// ErrorWithFields function
func (l *Logger) ErrorWithFields(msg string, KV logger.KV) {
	l.logger.Error().Timestamp().Fields(KV).Msg(msg)
}

// Errors function to log errors package
func (l *Logger) Errors(err error) {
	switch e := err.(type) {
	case *errors.Error:
		fields := e.GetFields()
		if fields == nil {
			fields = make(errors.Fields)
		}
		fields["operations"] = e.OpTraces

		l.logger.Error().Timestamp().Fields(fields).Msg(e.Error())
	case error:
		l.logger.Error().Timestamp().Msg(err.Error())
	}
}

// Fatal function
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal().Timestamp().Msg(fmt.Sprint(args...))
}

// Fatalln function
func (l *Logger) Fatalln(args ...interface{}) {
	l.logger.Fatal().Timestamp().Msg(fmt.Sprintln(args...))
}

// Fatalf function
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Fatal().Timestamp().Msgf(format, v...)
}

// FatalWithFields function
func (l *Logger) FatalWithFields(msg string, KV logger.KV) {
	l.logger.Fatal().Timestamp().Fields(KV).Msg(msg)
}
