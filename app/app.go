package app

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/kovey/cli-go/gui"
	"github.com/kovey/cli-go/util"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/debug-go/run"
)

type AppInterface interface {
	Get(name string) (*Flag, error)
	Name() string
	SetDebugLevel(t debug.DebugType)
	Flag(name string, def any, t Type, comment string)
}

type App struct {
	Action     func(AppInterface) error
	Reload     func(AppInterface) error
	Stop       func(AppInterface) error
	Show       func(*gui.Table)
	PidFile    func(AppInterface) string
	flags      map[string]*Flag
	ticker     *time.Ticker
	pid        int
	sigChan    chan os.Signal
	isStop     bool
	wait       sync.WaitGroup
	pidFile    string
	name       string
	isShowInfo bool
	serv       ServInterface
}

func NewApp(name string) *App {
	if len(name) == 0 {
		panic("app name is empty")
	}

	a := &App{
		flags: make(map[string]*Flag), sigChan: make(chan os.Signal, 1), isStop: false, isShowInfo: false,
		wait: sync.WaitGroup{}, pidFile: util.RunDir() + "/" + name + ".pid", name: name, ticker: time.NewTicker(1 * time.Minute),
	}
	a.flag("s", "no", TYPE_STRING, "signal: reload|info|stop")
	return a
}

func (a *App) SetServ(serv ServInterface) {
	a.serv = serv
}

func (a *App) flag(name string, def any, t Type, comment string) {
	_, ok := a.flags[name]
	if ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}

	a.flags[name] = &Flag{name: name, def: def, t: t, comment: comment}
}

func (a *App) Flag(name string, def any, t Type, comment string) {
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

	pid, err := os.ReadFile(a.pidFile)
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

		if err := syscall.Kill(pid, syscall.SIGUSR2); err != nil {
			debug.Erro("%s[%d] reload failure, error: %s", a.name, pid, err)
			return true
		}

		debug.Info("%s[%d] reload", a.name, pid)
		return true
	case "info":
		pid := a.getPid()
		if pid < 2 {
			debug.Erro("[%s] is not running", a.name)
			return true
		}
		if err := syscall.Kill(pid, syscall.SIGUSR1); err != nil {
			debug.Erro("%s[%d] show or hide info failure, error: %s", a.name, pid, err)
		}
		debug.Info("%s[%d] show or hide info", a.name, pid)
		return true
	case "stop":
		pid := a.getPid()
		if pid < 2 {
			debug.Erro("[%s] is not running", a.name)
			return true
		}

		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			debug.Erro("%s[%d] stop failure, error: %s", a.name, pid, err)
			return true
		}
		debug.Info("%s[%d] stop", a.name, pid)
		return true
	default:
		debug.Warn("unknown signal")
		return true
	}
}

func (a *App) Run() error {
	if a.serv != nil {
		if err := a.serv.Flag(a); err != nil {
			return err
		}
	}

	a.parse()
	if a.signal() {
		return nil
	}

	if a.PidFile != nil {
		a.pidFile = a.PidFile(a)
	}

	if util.IsFile(a.pidFile) {
		return fmt.Errorf("app[%s] is running", a.name)
	}

	if a.serv == nil {
		if a.Action == nil {
			return fmt.Errorf("Action is nil")
		}
	} else {
		if err := a.serv.Init(a); err != nil {
			return err
		}
	}

	defer func() {
		if a.pidFile != "" {
			os.Remove(a.pidFile)
		}
		run.Panic(recover())
	}()

	a.pid = os.Getpid()

	err := os.WriteFile(a.pidFile, []byte(a.PidString()), 0644)
	if err != nil {
		return err
	}

	signal.Notify(a.sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
	a.wait.Add(1)
	go a.listen()
	debug.Info("app[%s] run, pid[%s]", a.name, a.PidString())

	startTime = time.Now()
	if a.serv != nil {
		err = a.serv.Run(a)
	} else {
		err = a.Action(a)
	}
	if !a.isStop {
		a.sigChan <- os.Interrupt
	}
	a.wait.Wait()
	return err
}

func (a *App) listen() {
	defer a.wait.Done()
	defer a.ticker.Stop()

	for {
		select {
		case _, ok := <-a.ticker.C:
			if !ok {
				a.isStop = true
				return
			}
			if a.isShowInfo {
				a.show()
			}
		case sig := <-a.sigChan:
			switch sig {
			case syscall.SIGUSR2:
				if a.serv != nil {
					if err := a.serv.Reload(a); err != nil {
						debug.Erro("serv[%s] reload failure, error: %s", a.name, err)
					}
				} else {
					if a.Reload != nil {
						if err := a.Reload(a); err != nil {
							debug.Erro("serv[%s] reload failure, error: %s", a.name, err)
						}
					}
				}
			case syscall.SIGUSR1:
				a.isShowInfo = !a.isShowInfo
			default:
				a.isStop = true
				if a.serv != nil {
					if err := a.serv.Shutdown(a); err != nil {
						debug.Erro("serv[%s] shutdown failure, error: %s", a.name, err)
						return
					}
				}
				if a.Stop != nil {
					if err := a.Stop(a); err != nil {
						debug.Erro("serv[%s] shutdown failure, error: %s", a.name, err)
					}
				}
				return
			}
		}
	}
}

func (a *App) Name() string {
	return a.name
}

func (a *App) SetDebugLevel(t debug.DebugType) {
	debug.SetLevel(t)
}

func (a *App) show() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	table := gui.NewTable()
	table.Add(fmt.Sprintf("app[%s]", a.name))
	table.Add(fmt.Sprintf("golang version[%s]", runtime.Version()))
	table.Add(fmt.Sprintf("start time[%s]", StartTimestamp()))
	table.Add(fmt.Sprintf("run time[%s]", GetFormatRunTime()))
	table.Add(fmt.Sprintf("total alloc[%d](bytes)", m.TotalAlloc))
	table.Add(fmt.Sprintf("alloc[%d](bytes)", m.Alloc))
	table.Add(fmt.Sprintf("active objects[%d]", m.Mallocs))
	table.Add(fmt.Sprintf("free objects[%d]", m.Frees))
	table.Add(fmt.Sprintf("heap alloc[%d](bytes)", m.HeapAlloc))
	table.Add(fmt.Sprintf("heap idle[%d](bytes)", m.HeapIdle))
	table.Add(fmt.Sprintf("heap released[%d](bytes)", m.HeapReleased))
	table.Add(fmt.Sprintf("heap sys[%d](bytes)", m.HeapSys))
	table.Add(fmt.Sprintf("heap in use[%d](bytes)", m.HeapInuse))
	table.Add(fmt.Sprintf("heap objects[%d]", m.HeapObjects))
	table.Add(fmt.Sprintf("stack in use[%d](bytes)", m.StackInuse))
	table.Add(fmt.Sprintf("stack sys[%d](bytes)", m.StackSys))
	table.Add(fmt.Sprintf("gc cpu fraction[%f](ms)", m.GCCPUFraction))
	table.Add(fmt.Sprintf("gc sys[%d](bytes)", m.GCSys))
	if a.Show != nil {
		table.AddRepeat(gui.Border_Horizontal)
		a.Show(table)
	}
	table.Show()
}
