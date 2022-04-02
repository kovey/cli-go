package app

import (
	"cli/util"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
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
		panic("name is empty")
	}

	a := &App{
		flags: make(map[string]*Flag), sigChan: make(chan os.Signal, 1), isStop: false,
		wait: sync.WaitGroup{}, pidFile: util.RunDir() + "/" + name + ".pid", name: name,
	}
	a.Flag("s", "no", TYPE_STRING, "signal: reload|maintain|stop")
	return a
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

	sig, ok := s.(*string)
	if !ok {
		return false
	}

	if *sig == "no" {
		return false
	}

	switch *sig {
	case "reload":
		pid := a.getPid()
		if pid < 2 {
			fmt.Printf("%s is not running\n", a.name)
			return true
		}

		syscall.Kill(pid, syscall.SIGUSR2)
		fmt.Printf("%s[%d] reload\n", a.name, pid)
		return true
	case "maintain":
		pid := a.getPid()
		if pid < 2 {
			fmt.Printf("%s is not running\n", a.name)
			return true
		}
		syscall.Kill(pid, syscall.SIGUSR1)
		fmt.Printf("%s[%d] maintain\n", a.name, pid)
		return true
	case "stop":
		pid := a.getPid()
		if pid < 2 {
			fmt.Printf("%s is not running\n", a.name)
			return true
		}

		syscall.Kill(pid, syscall.SIGTERM)
		fmt.Printf("%s[%d] stop\n", a.name, pid)
		return true
	default:
		fmt.Println("unknown signal")
		return true
	}
}

func (a *App) Run() error {
	a.Parse()
	if a.signal() {
		return nil
	}

	if a.Action == nil {
		return fmt.Errorf("Action is nil")
	}

	a.pid = os.Getpid()

	if a.PidFile != nil {
		a.pidFile = a.PidFile(a)
	}

	err := ioutil.WriteFile(a.pidFile, []byte(a.PidString()), 0644)
	if err != nil {
		return err
	}

	signal.Notify(a.sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
	a.wait.Add(1)
	go a.listen()
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
