package debug

import (
	"fmt"
	"runtime"
	"time"

	"github.com/kovey/cli-go/color"
	"github.com/kovey/cli-go/util"
)

type DebugType string

type DebugValue int32

type DebugLevels map[DebugType]DebugValue

var level DebugValue = val_info

func (d DebugLevels) CanShow(t DebugType) bool {
	if l, ok := d[t]; ok {
		return l >= level
	}

	return false
}

const (
	val_info DebugValue = 1
	val_warn DebugValue = 2
	val_erro DebugValue = 3
	val_none DebugValue = 4
)

const (
	info DebugType = "info"
	erro DebugType = "erro"
	warn DebugType = "warn"
)

const (
	echoFormat = "[%s][%s] %s\n"
)

var maps = DebugLevels{
	info: val_info,
	warn: val_warn,
	erro: val_erro,
}

func echo(format string, t DebugType, args ...interface{}) {
	if !maps.CanShow(t) {
		return
	}

	switch t {
	case warn:
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Yellow(str))
	case erro:
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Red(str))
	default:
		fmt.Printf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
	}
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
