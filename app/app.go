package app

import (
	"time"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/debug-go/debug"
)

type CleanFlag byte

const (
	Clean_Defalut CleanFlag = 0
	Clean_Help    CleanFlag = 1
	Clean_Version CleanFlag = 1 << 1
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
	// @deprecated
	CleanCommandLine(isCleanHelp bool)
	CleanCommandLineWith(flags CleanFlag)
	RunChild(func(AppInterface))
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
