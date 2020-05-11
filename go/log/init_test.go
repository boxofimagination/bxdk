package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/boxofimagination/bxdk/go/log/logger"
)

func TestSkip(t *testing.T) {
	opt := &Option{
		Level: logger.DebugLevelString,
		Output: Output{
			Debug: "debug.log",
			Info:  "info.log",
			Warn:  "warn.log",
			Error: "error.log",
		},
		Caller:     true,
		CallerSkip: 1,
	}

	defer func() {
		for i := DebugLevel; i < FatalLevel; i++ {
			os.Remove(opt.Output.at(i))
		}
	}()

	require.NoError(t, newInit(opt))
	_debug("hello")
	_info("hello")
	_warn("hello")
	_error("hello")

	for i := DebugLevel; i < FatalLevel; i++ {
		res, err := ioutil.ReadFile(opt.Output.at(i))
		require.NoError(t, err)
		require.Contains(t, string(res), "init_test.go")
	}
}
