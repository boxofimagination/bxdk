package log

import (
	"github.com/boxofimagination/bxdk/go/env"

	"github.com/boxofimagination/bxdk/go/log/logger"
	"github.com/boxofimagination/bxdk/go/log/logger/zerolog"
)

// Output file of log
type Output struct {
	Debug string `json:"debug" yaml:"debug"`
	Info  string `json:"info" yaml:"info"`
	Warn  string `json:"warn" yaml:"warn"`
	Error string `json:"error" yaml:"error"`
	Fatal string `json:"fatal" yaml:"fatal"`
}

func (o *Output) at(level Level) string {
	switch level {
	case DebugLevel:
		return o.Debug
	case InfoLevel:
		return o.Info
	case WarnLevel:
		return o.Warn
	case ErrorLevel:
		return o.Error
	case FatalLevel:
		return o.Fatal
	default:
		return ""
	}
}

// Option of logger
type Option struct {
	Level      string `json:"level" yaml:"level"`
	TimeFormat string `json:"time_format" yaml:"time_format"`
	Output     Output `json:"output" yaml:"output"`

	// Caller, option to print caller line numbers.
	// make sure you understand the overhead when use this
	Caller bool `json:"caller" yaml:"caller"`
	// CallerSkip skips additional wrapper
	CallerSkip int `json:"caller_skip" yaml:"caller_skip"`

	// TODO: Metadata / default fields in each log
}

// newInit configures log to use the given option
// Send nil to opt parameter to use default config
func newInit(opt *Option) error {
	var (
		lgr   [5]Logger
		err   error
		isDev = env.IsDevelopment()
	)

	cfg := &logger.Config{
		Level:      InfoLevel,
		TimeFormat: logger.DefaultTimeFormat,
		UseColor:   isDev,
		UseJSON:    !isDev,
	}
	if opt != nil {
		cfg.Level = logger.StringToLevel(opt.Level)
		cfg.TimeFormat = opt.TimeFormat
		cfg.Caller = opt.Caller
		cfg.CallerSkip = opt.CallerSkip
	}

	for i := DebugLevel; i <= FatalLevel; i++ {
		lgr[i], err = zerolog.New(cfg, withOutputFile(opt.Output.at(i)))
		if err != nil {
			return err
		}
	}

	setLoggers(lgr[0], lgr[1], lgr[2], lgr[3], lgr[4])

	return nil
}

func withOutputFile(filepath string) func(*logger.Config) {
	return func(cfg *logger.Config) {
		cfg.LogFile = filepath
	}
}

func setLoggers(debug, info, warn, err, fatal Logger) {
	if debug != nil {
		debugLogger = debug
	}
	if info != nil {
		infoLogger = info
	}
	if warn != nil {
		warnLogger = warn
	}
	if err != nil {
		errorLogger = err
	}
	if fatal != nil {
		fatalLogger = fatal
	}
}
