package app

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

type App struct {
	Action  func(*App) error
	Reload  func(*App) error
	Stop    func(*App) error
	flags   map[string]*Flag
	pid     int
	sigChan chan os.Signal
	isStop  bool
	wait    sync.WaitGroup
}

func NewApp() *App {
	return &App{flags: make(map[string]*Flag), sigChan: make(chan os.Signal, 1), isStop: false, wait: sync.WaitGroup{}}
}

func (a *App) Flag(name string, def interface{}, t Type, comment string) {
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

func (a *App) PidString() string {
	return strconv.Itoa(a.pid)
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
	signal.Notify(a.sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2)
	a.wait.Add(1)
	go a.listen()
	err := a.Action(a)
	if !a.isStop {
		a.sigChan <- os.Interrupt
	}
	a.wait.Wait()
	return err
}

func (a *App) listen() {
	defer a.wait.Done()
	defer close(a.sigChan)
loop:
	for {
		select {
		case sig := <-a.sigChan:
			if sig == syscall.SIGUSR2 {
				if a.Reload != nil {
					a.Reload(a)
				}

				continue loop
			}

			if a.Stop != nil {
				a.isStop = true
				a.Stop(a)
			}

			break loop
		}
	}
}
