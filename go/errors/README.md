# TDK Error

## Error

Error is value in Go, and because error is a value, we need to check them. But don't only check them, handle them gracefully as Go proverbs said:

> Don't just check errors, handle them gracefully

and

> Log an error or return the error, never both - Dave Channey

### Why another error package

As `Dave Channey` said, we should log an error or just return the error, but never both. But how can we log a meaningful error in go and still can compare the error itself?

In order to do that, we need a modified implementation of error. Put more context into error and print the context when we need to log. That way we don't need to log and return an error at the same time, just to put more context into the log.

### Errors function

To create a meaningful error from this package, we need to use `errors.E(args...)` function. Why `errors.E()` instead of `errors.New()`like `errors` package from Go itself?

1. Following `upspin` convention to create the error
2. Let the standard be a standard(`errors.New`) and the new one should have a new convention.

## Example

### Simple error creation

```go
import "github.com/boxofimagination/bxdk/errors"

func main() {
    err := errors.E("this is error from bxdk")
    // do something with the error
}

```

### Error with fields

Error with fields is useful to give context to error. For example `userid` of user.

```go
import "github.com/boxofimagination/bxdk/errors"

func main() {
    err := errors.E("this is error from bxdk", errors.Fields{"user_id": 1234})
    // do something with the error
}
```

### Error with operations

Sometimes we need to know what kind of operations we do in error, we want to know where exactly error happens.

```go
import "github.com/boxofimagination/bxdk/errors"

func main() {
    err := SomeFunction()
    // do something with the error
}

func SomeFunction() error {
    const op errors.Op = "main/somefunction"
    return errors.E(op, "this is error from bxdk")
}
```

### Stack trace in error

Tracing stack when spawning errors is not desireable. Calling the whole stack and parse the stack will surely come with performance degradation. In high traffic applications, this is not preferrable and is not a recommended practice.

But to make development more easier stack trace will be enabled when `TKPENV=development`.

### Real life example

This is an example where we need to call a function from handler and we need to know the error context

```go
import (
    "net/http"
    "strings"

    "github.com/boxofimagination/bxdk/errors"
)

// Error variables for error matching example
var (
    // Given string parameter on errors.E() will directly converted to error message
    ErrParamIsEqual      = errors.E("param1 is equal to param2")
    ErrMoreThanConstanta = errors.E("param1 length is more than constanta")
)

func main() {
    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        param1 := r.FormValue("param1")
        param2 := r.FormValue("param2")

        err := BusinessLogic(param1, param2)
        // sample implementation of errors.Match() to handle error regarding to error types
        if errors.Match(err, ErrParamIsEqual) {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Not OK"))
            return
        }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    http.ListenAndServe(":9090", nil)
}

func BusinessLogic(param1, param2 string) error {
    const op errors.Op = "business/BusinessLogic"

    if strings.Compare(param1, param2) == 0 {
        return errors.E(ErrParamIsEqual, errors.Fields{
            "param1": param1,
            "param2": param2,
        }, op)
    }
    return ResourceLogic(param1)
}

const constVal string = "constanta"

func ResourceLogic(param1 string) error {
    const op errors.Op = "resource/ResourceLogic"

    if len(param1) > len(constVal) {
        return errors.E(ErrMoreThanConstanta, op, errors.Fields{"param1": param1})
    }
    return nil
}

```