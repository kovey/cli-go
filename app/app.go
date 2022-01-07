package app

import (
	"flag"
	"fmt"
	"os"
)

type App struct {
	Action func(*App) error
	flags  map[string]*Flag
	pid    int
}

func NewApp() *App {
	return &App{flags: make(map[string]*Flag)}
}

func (a *App) Add(name string, def interface{}, t Type, comment string) {
	a.flags[name] = &Flag{Name: name, Default: def, Type: t, Comment: comment}
}

func (a *App) Get(name string) (interface{}, error) {
	f, ok := a.flags[name]
	if !ok {
		return nil, fmt.Errorf("[%s] is not exists", name)
	}

	return f.Value, nil
}

func (a *App) Pid() int {
	return a.pid
}

func (a *App) Parse() {
	for _, flag := range a.flags {
		flag.Parse()
	}

	flag.Parse()
}

func (a *App) Run() error {
	if a.Action == nil {
		return fmt.Errorf("Action is nil")
	}

	a.pid = os.Getpid()
	a.Parse()
	return a.Action(a)
}
