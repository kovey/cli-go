package app

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/kovey/cli-go/debug"
	"github.com/kovey/cli-go/util"
)

type App struct {
	Action   func(*App) error
	Reload   func(*App) error
	Stop     func(*App) error
	Maintain func(*App) error
	PidFile  func(*App) string
	flags    map[string]*Flag
	pid      int
	sigChan  chan os.Signal
	isStop   bool
	wait     sync.WaitGroup
	pidFile  string
	name     string
}

func NewApp(name string) *App {
	if len(name) == 0 {
		panic("app name is empty")
	}

	a := &App{
		flags: make(map[string]*Flag), sigChan: make(chan os.Signal, 1), isStop: false,
		wait: sync.WaitGroup{}, pidFile: util.RunDir() + "/" + name + ".pid", name: name,
	}
	a.flag("s", "no", TYPE_STRING, "signal: reload|maintain|stop")
	return a
}

func (a *App) flag(name string, def interface{}, t Type, comment string) {
	_, ok := a.flags[name]
	if ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}

	a.flags[name] = &Flag{name: name, def: def, t: t, comment: comment}
}

func (a *App) Flag(name string, def interface{}, t Type, comment string) {
	if name == "s" {
		debug.Warn("flag[%s] is used by sinal module", name)
		return
	}

	a.flag(name, def, t, comment)
}

func (a *App) Get(name string) (*Flag, error) {
	f, ok := a.flags[name]
	if !ok {
		return nil, fmt.Errorf("[%s] is not exists", name)
	}

	return f, nil
}

func (a *App) Pid() int {
	return a.pid
}

func (a *App) PidString() string {
	return strconv.Itoa(a.pid)
}

func (a *App) parse() {
	for _, flag := range a.flags {
		flag.parse()
	}

	flag.Parse()
}

func (a *App) getPid() int {
	if a.PidFile != nil {
		a.pidFile = a.PidFile(a)
	}

	pid, err := ioutil.ReadFile(a.pidFile)
	if err != nil {
		return -1
	}

	id, err := strconv.Atoi(string(pid))
	if err != nil {
		return -1
	}

	return id
}

func (a *App) signal() bool {
	s, err := a.Get("s")
	if err != nil {
		return false
	}

	sig := s.String()
	if sig == "no" {
		return false
	}

	switch sig {
	case "reload":
		pid := a.getPid()
		if pid < 2 {
			debug.Erro("[%s] is not running", a.name)
			return true
		}

		syscall.Kill(pid, syscall.SIGUSR2)
		debug.Info("%s[%d] reload", a.name, pid)
		return true
	case "maintain":
		pid := a.getPid()
		if pid < 2 {
			debug.Erro("[%s] is not running", a.name)
			return true
		}
		syscall.Kill(pid, syscall.SIGUSR1)
		debug.Info("%s[%d] maintain", a.name, pid)
		return true
	case "stop":
		pid := a.getPid()
		if pid < 2 {
			debug.Erro("[%s] is not running", a.name)
			return true
		}

		syscall.Kill(pid, syscall.SIGTERM)
		debug.Info("%s[%d] stop", a.name, pid)
		return true
	default:
		debug.Warn("unknown signal")
		return true
	}
}

func (a *App) Run() error {
	a.parse()
	if a.signal() {
		return nil
	}

	if a.Action == nil {
		return fmt.Errorf("Action is nil")
	}

	if a.PidFile != nil {
		a.pidFile = a.PidFile(a)
	}

	if util.IsFile(a.pidFile) {
		return fmt.Errorf("app[%s] is running", a.name)
	}

	a.pid = os.Getpid()

	err := ioutil.WriteFile(a.pidFile, []byte(a.PidString()), 0644)
	if err != nil {
		return err
	}

	signal.Notify(a.sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
	a.wait.Add(1)
	go a.listen()
	debug.Info("app[%s] run, pid[%s]", a.name, a.PidString())

	err = a.Action(a)
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
			switch sig {
			case syscall.SIGUSR2:
				if a.Reload != nil {
					a.Reload(a)
				}
				break
			case syscall.SIGUSR1:
				if a.Maintain != nil {
					a.Maintain(a)
				}
				break
			default:
				if a.Stop != nil {
					a.isStop = true
					os.Remove(a.pidFile)
					a.Stop(a)
				}
				break loop
			}
		}
	}
}
