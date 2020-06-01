package log

import "github.com/boxofimagination/bxdk/go/log/logger"


// KV is a type for passing more fields information
type KV = logger.KV

// Debug prints debug level log like log.Print
func Debug(args ...interface{}) {
	debugLogger.Debug(args...)
}

// Debugln prints debug level log like log.Println
func Debugln(args ...interface{}) {
	debugLogger.Debugln(args...)
}

// Debugf prints debug level log like log.Printf
func Debugf(format string, v ...interface{}) {
	debugLogger.Debugf(format, v...)
}

// DebugWithFields prints debug level log with additional fields.
// useful when output is in json format
func DebugWithFields(msg string, fields KV) {
	debugLogger.DebugWithFields(msg, fields)
}

// Print info level log like log.Print
func Print(v ...interface{}) {
	infoLogger.Info(v...)
}

// Println info level log like log.Println
func Println(v ...interface{}) {
	infoLogger.Infoln(v...)
}

// Printf info level log like log.Printf
func Printf(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// Info prints info level log like log.Print
func Info(args ...interface{}) {
	infoLogger.Info(args...)
}

// Infoln prints info level log like log.Println
func Infoln(args ...interface{}) {
	infoLogger.Infoln(args...)
}

// Infof prints info level log like log.Printf
func Infof(format string, v ...interface{}) {
	infoLogger.Infof(format, v...)
}

// InfoWithFields prints info level log with additional fields.
// useful when output is in json format
func InfoWithFields(msg string, fields KV) {
	infoLogger.InfoWithFields(msg, fields)
}

// Warn prints warn level log like log.Print
func Warn(args ...interface{}) {
	warnLogger.Warn(args...)
}

// Warnln prints warn level log like log.Println
func Warnln(args ...interface{}) {
	warnLogger.Warnln(args...)
}

// Warnf prints warn level log like log.Printf
func Warnf(format string, v ...interface{}) {
	warnLogger.Warnf(format, v...)
}

// WarnWithFields prints warn level log with additional fields.
// useful when output is in json format
func WarnWithFields(msg string, fields KV) {
	warnLogger.WarnWithFields(msg, fields)
}

// Error prints error level log like log.Print
func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

// Errorln prints error level log like log.Println
func Errorln(args ...interface{}) {
	errorLogger.Errorln(args...)
}

// Errorf prints error level log like log.Printf
func Errorf(format string, v ...interface{}) {
	errorLogger.Errorf(format, v...)
}

// ErrorWithFields prints error level log with additional fields.
// useful when output is in json format
func ErrorWithFields(msg string, fields KV) {
	errorLogger.ErrorWithFields(msg, fields)
}

// Errors can handle error from bxdk/x/go/errors package
func Errors(err error) {
	errorLogger.Errors(err)
}

// Fatal prints fatal level log like log.Print
func Fatal(args ...interface{}) {
	fatalLogger.Fatal(args...)
}

// Fatalln prints fatal level log like log.Println
func Fatalln(args ...interface{}) {
	fatalLogger.Fatalln(args...)
}

// Fatalf prints fatal level log like log.Printf
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Fatalf(format, v...)
}

// FatalWithFields prints fatal level log with additional fields.
// useful when output is in json format
func FatalWithFields(msg string, fields KV) {
	fatalLogger.FatalWithFields(msg, fields)
}
