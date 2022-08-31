package debug

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kovey/cli-go/util"
)

type DebugType string

const (
	info DebugType = "info"
	erro DebugType = "erro"
	warn DebugType = "warn"
)

const (
	echoFormat = "[%s][%s] %s\n"
)

func echo(format string, t DebugType, args ...interface{}) {
	fmt.Printf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	echo(format, info, args...)
}

func Erro(format string, args ...interface{}) {
	echo(format, erro, args...)
}

func Warn(format string, args ...interface{}) {
	echo(format, warn, args...)
}

func RunCo(f func()) {
	defer func() {
		Panic(recover())
	}()
	f()
}

func Panic(err interface{}) bool {
	if err == nil {
		return false
	}

	Erro("panic error[%s]", err)

	for i := 3; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		Erro("%s(%d)", file, line)
	}

	return true
}
