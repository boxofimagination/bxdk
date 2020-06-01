# BXDK Log

BXDK Log is based on [Zerolog](https://github.com/rs/zerolog)

## Log level

Logger supports log levels, available log levels are:

- Debug
- Info
- Warning
- Error
- Fatal

The log is disabled if `LogLevel` < `CurrentLogLevel`. For example, the `Debug` log is disabled when the current level is `Info`.

## Default Config

BXdk log default config:
- debug log level
- RFC3339 time format
- No caller
- No output file
- Colored console format if BOXENV=development
- JSON format if BOXENV!=development

## Customizing log

**IMPORTANT NOTE: the functions in this section are not thread-safe, call them when initializing the application.**

### 1. Change Log Level

You can set the log level using `SetLevel` or `SetLevelString`.

Example of `SetLevel`:

```go
import "github.com/boxofimagination/bxdk/go/log"

func main() {
    log.SetLevel(log.InfoLevel)
    log.Info("this is a log")
}
```

Example of `SetLevelString` (use lowercase)
```go
import "github.com/boxofimagination/bxdk/go/log"

func main() {
    log.SetLevelString("info")
    log.Info("this is a log")
}
```

### 2. Change more configuration

Use `Init` to initialize the logger with your configurations:

```go
import "github.com/boxofimagination/bxdk/go/log"

func main() {
    err := log.Init(&log.Option{
        Level:      "info",
        Output:     log.Output {
            Debug:  "info.log",
            Info:   "info.log",
            Warn:   "error.log",
            Error:  "error.log",
            Fatal:  "error.log",
        },
    })
    if err != nil {
        panic(err)
    }
    log.Info("this is a log")
}
```

`log.Option`:

| Field | Tag (JSON & YAML) | Type | Description |
|-|-|-|-|
| Level | level | string | specify log level, default: "info" |
| TimeFormat | time_format | string | specify time format, default: "2006-01-02T15:04:05Z07:00" |
| Caller | caller | bool | print caller position or not, default: false |
| CallerSkip | caller_skip | int | specify additional wrapper, default: 0 |
| Output | output | Output | specify file path for each level |

`log.Output`:

| Field | Tag (JSON & YAML) | Type | Description |
|-|-|-|-|
| Debug | debug | string | specify file path for debug level log |
| Info | info | string | specify file path for info level log |
| Warn | warn | string | specify file path for warn level log |
| Error | error | string | specify file path for error level log |
| Fatal | fatal | string | specify file path for fatal level log |


### 3. Change logger

If you want even more flexible configuration, you can set the logger per level by yourself using `SetLogger` \
Note: your custom logger will be replaced if you call `Init` after `SetLogger`

example: set a separate output file for `errorLevel` and `fatalLevel` log
```go
import "github.com/boxofimagination/bxdk/go/log"
import "github.com/boxofimagination/bxdk/go/log/logger"

func main() {
    errLogger, err := log.NewLogger(log.Zerolog, &logger.Config{
        Level:   log.DebugLevel,
        LogFile: "error.log",
    })
    if err != nil {
        panic(err)
    }
    err = log.SetLogger(log.ErrorLevel, errLogger)
    if err != nil {
        panic(err)
    }
    err = log.SetLogger(log.FatalLevel, errLogger)
    if err != nil {
        panic(err)
    }
    log.Error("this is a log")
}
```

## Available function

Each level has 3 function:
- unformatted -> like Println in standard log
- formatted -> like Printf in standard log
- with map[string]interface{} (or log.KV)

We also add several functions to ease the migration to bxdk log.
Function list and example:

```go
arg1 := "hello"
arg2 := "world"

log.Debug(arg1, arg2)
log.Debugln(arg1, arg2)
log.Debugf("message %v %v", arg1, arg2)
log.DebugWithFields("message", log.KV{"arg1":arg1, "arg2":arg2})

log.Info(arg1, arg2)
log.Infoln(arg1, arg2)
log.Infof("message %v %v", arg1, arg2)
log.InfoWithFields("message", log.KV{"arg1":arg1, "arg2":arg2})

// alias for Info
log.Print(arg1, arg2)
// alias for Infoln
log.Println(arg1, arg2)
// alias for Infof
log.Printf("message %v %v", arg1, arg2)

log.Warn(arg1, arg2)
log.Warnln(arg1, arg2)
log.Warnf("message %v %v", arg1, arg2)
log.WarnWithFields("message", log.KV{"arg1":arg1, "arg2":arg2})

log.Error(arg1, arg2)
log.Errorln(arg1, arg2)
log.Errorf("message %v %v", arg1, arg2)
log.ErrorWithFields("message", log.KV{"arg1":arg1, "arg2":arg2})

log.Fatal(arg1, arg2)
log.Fatalln(arg1, arg2)
log.Fatalf("message %v %v", arg1, arg2)
log.FatalWithFields("message", log.KV{"arg1":arg1, "arg2":arg2})
```

Output example:
```
2019-03-08 17:16:48+07:00 DBG log_test.go:35 > helloworld
2019-03-08 17:16:48+07:00 DBG log_test.go:35 > hello world
2019-03-08 17:16:48+07:00 DBG log_test.go:36 > message hello world
2019-03-08 17:16:48+07:00 DBG log_test.go:37 > message arg1=hello arg2=world
2019-03-08 17:16:48+07:00 INF log_test.go:39 > helloworld
2019-03-08 17:16:48+07:00 INF log_test.go:39 > hello world
2019-03-08 17:16:48+07:00 INF log_test.go:40 > message hello world
2019-03-08 17:16:48+07:00 INF log_test.go:41 > message arg1=hello arg2=world
2019-03-08 17:16:48+07:00 INF log_test.go:44 > helloworld
2019-03-08 17:16:48+07:00 INF log_test.go:45 > hello world
2019-03-08 17:16:48+07:00 INF log_test.go:47 > message hello world
2019-03-08 17:16:48+07:00 WRN log_test.go:49 > helloworld
2019-03-08 17:16:48+07:00 WRN log_test.go:49 > hello world
2019-03-08 17:16:48+07:00 WRN log_test.go:50 > message hello world
2019-03-08 17:16:48+07:00 WRN log_test.go:51 > message arg1=hello arg2=world
2019-03-08 17:16:48+07:00 ERR log_test.go:53 > helloworld
2019-03-08 17:16:48+07:00 ERR log_test.go:53 > hello world
2019-03-08 17:16:48+07:00 ERR log_test.go:54 > message hello world
2019-03-08 17:16:48+07:00 ERR log_test.go:55 > message arg1=hello arg2=world
2019-03-08 17:16:48+07:00 FTL log_test.go:57 > helloworld
2019-03-08 17:16:48+07:00 FTL log_test.go:57 > hello world
2019-03-08 17:16:48+07:00 FTL log_test.go:57 > message hello world
2019-03-08 17:16:48+07:00 FTL log_test.go:57 > message arg1=hello arg2=world
```

## Integration with BXDK Error package

BXDK error package has a features called `errors.Fields`. This fields can be used to add more context into the error, and then we can print the fields when needed. BXDK log will automatically print the fields if `error = bxdkerrors.Error` by using `log.Errors`. For example:

```go
import "github.com/boxofimagination/bxdk/go/log"
import "github.com/boxofimagination/bxdk/go/errors"

func main() {
    err := errors.E("this is an error", errors.Fields{"field1":"value1"})
    log.Errors(err)
}

// result is
// message=this is an error field1=value1
```
