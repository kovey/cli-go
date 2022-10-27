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

func (d DebugLevels) Get(t DebugType) DebugValue {
	if val, ok := d[t]; ok {
		return val
	}

	return val_none
}

const (
	val_info DebugValue = 1
	val_warn DebugValue = 2
	val_erro DebugValue = 3
	val_none DebugValue = 4
)

const (
	Debug_Info DebugType = "info"
	Debug_Erro DebugType = "erro"
	Debug_Warn DebugType = "warn"
)

const (
	echoFormat = "[%s][%s] %s\n"
)

var maps = DebugLevels{
	Debug_Info: val_info,
	Debug_Warn: val_warn,
	Debug_Erro: val_erro,
}

func SetLevel(t DebugType) {
	level = maps.Get(t)
}

func echo(format string, t DebugType, args ...interface{}) {
	if !maps.CanShow(t) {
		return
	}

	switch t {
	case Debug_Warn:
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Yellow(str))
	case Debug_Erro:
		str := fmt.Sprintf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
		fmt.Print(color.Red(str))
	default:
		fmt.Printf(echoFormat, time.Now().Format(util.GOLANG_BIRTHDAY), t, fmt.Sprintf(format, args...))
	}
}

func Info(format string, args ...interface{}) {
	echo(format, Debug_Info, args...)
}

func Erro(format string, args ...interface{}) {
	echo(format, Debug_Erro, args...)
}

func Warn(format string, args ...interface{}) {
	echo(format, Debug_Warn, args...)
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
