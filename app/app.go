package app

import (
	"time"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
)

func init() {
	loadEnv(time.Now())
}

type AppInterface interface {
	Get(names ...string) (*Flag, error)
	Name() string
	SetDebugLevel(t debug.DebugType)
	Flag(name string, def any, t Type, comment string, parents ...string)
	FlagLong(name string, def any, t Type, comment string, parents ...string)
	FlagNonValueLong(name string, comment string, parents ...string)
	FlagNonValue(name string, comment string, parents ...string)
	FlagArg(name string, comment string, parents ...string)
	Arg(index int, t Type) (*Flag, error)
	UsageWhenErr()
}

type App struct {
	*Daemon
}

func loadEnv(now time.Time) {
	if !env.HasEnv() {
		return
	}

	if err := env.LoadDefault(now); err != nil {
		debug.Erro(err.Error())
	}
}

func NewApp(name string) *App {
	a := &App{Daemon: NewDaemon(name)}
	return a
}
