package app

import (
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
	FlagLong(name string, def any, t Type, comment string)
	FlagNonValueLong(name string, comment string)
	FlagNonValue(name string, comment string)
	Arg(index int, t Type) (*Flag, error)
	UsageWhenErr()
}

type App struct {
	Action     func(AppInterface) error
	Reload     func(AppInterface) error
	Stop       func(AppInterface) error
	Panic      func(AppInterface)
	PidFile    func(AppInterface) string
	ticker     *time.Ticker
	pid        int
	sigChan    chan os.Signal
	isStop     bool
	wait       sync.WaitGroup
	pidFile    string
	name       string
	isShowInfo bool
	serv       ServInterface
	showUsage  bool
}

func NewApp(name string) *App {
	if len(name) == 0 {
		panic("app name is empty")
	}

	a := &App{
		sigChan: make(chan os.Signal, 1), isStop: false, isShowInfo: false,
		wait: sync.WaitGroup{}, pidFile: util.RunDir() + "/" + name + ".pid", name: name, ticker: time.NewTicker(1 * time.Minute),
	}
	_commanLine.Flag("s", "no", TYPE_STRING, "signal: reload|info|stop")
	return a
}

func (a *App) UsageWhenErr() {
	a.showUsage = true
}

func (a *App) SetServ(serv ServInterface) {
	a.serv = serv
	Usage = serv.Usage
}

func (a *App) FlagNonValueLong(name string, comment string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagNonValueLong(name, comment)
}

func (a *App) FlagNonValue(name string, comment string) {
	if name == "s" {
		debug.Warn("flag[%s] is used by sinal module", name)
		return
	}

	_commanLine.FlagNonValue(name, comment)
}

func (a *App) FlagLong(name string, def any, t Type, comment string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagLong(name, def, t, comment)
}

func (a *App) Flag(name string, def any, t Type, comment string) {
	if name == "s" {
		debug.Warn("flag[%s] is used by sinal module", name)
		return
	}

	_commanLine.Flag(name, def, t, comment)
}

func (a *App) Get(name string) (*Flag, error) {
	f := _commanLine.Get(name)
	if f == nil {
		return nil, fmt.Errorf("[%s] is not exists", name)
	}

	return f, nil
}

func (a *App) Arg(index int, t Type) (*Flag, error) {
	f := _commanLine.Arg(index)
	if f == nil {
		return nil, fmt.Errorf("[%d] is not exists", index)
	}

	f.t = t
	return f, nil
}

func (a *App) Pid() int {
	return a.pid
}

func (a *App) PidString() string {
	return strconv.Itoa(a.pid)
}

func (a *App) parse() {
	_commanLine.Parse(os.Args[1:])
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
	err = a._run()
	if !a.isStop {
		a.sigChan <- os.Interrupt
	}
	a.wait.Wait()
	if err != nil && a.showUsage {
		Usage()
	}

	return err
}

func (a *App) _run() error {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		if a.serv != nil {
			a.serv.Panic(a)
		} else {
			if a.Panic != nil {
				a.Panic(a)
			}
		}

		run.Panic(err)
	}()

	if a.serv != nil {
		return a.serv.Run(a)
	}

	return a.Action(a)
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
	table.Add(0, "Cli Info Name")
	table.Add(0, "Value")
	table.Add(1, "app")
	table.Add(1, a.name)
	table.Add(2, "golang version")
	table.Add(2, runtime.Version())
	table.Add(3, "start time")
	table.Add(3, StartTimestamp())
	table.Add(4, "run time")
	table.Add(4, GetFormatRunTime())
	table.Add(5, "total alloc(bytes)")
	table.AddAny(5, m.TotalAlloc)
	table.Add(6, "alloc(bytes)")
	table.AddAny(6, m.Alloc)
	table.Add(7, "active objects")
	table.AddAny(7, m.Mallocs)
	table.Add(8, "free objects")
	table.AddAny(8, m.Frees)
	table.Add(9, "heap alloc(bytes)")
	table.AddAny(9, m.HeapAlloc)
	table.Add(10, "heap idle(bytes)")
	table.AddAny(10, m.HeapIdle)
	table.Add(11, "heap released(bytes)")
	table.AddAny(11, m.HeapReleased)
	table.Add(12, "heap sys(bytes)")
	table.AddAny(12, m.HeapSys)
	table.Add(13, "heap in use(bytes)")
	table.AddAny(13, m.HeapInuse)
	table.Add(14, "heap objects")
	table.AddAny(14, m.HeapObjects)
	table.Add(15, "stack in use(bytes)")
	table.AddAny(15, m.StackInuse)
	table.Add(16, "stack sys(bytes)")
	table.AddAny(16, m.StackSys)
	table.Add(17, "gc cpu fraction(ms)")
	table.AddAny(17, m.GCCPUFraction)
	table.Add(18, "gc sys(bytes)")
	table.AddAny(18, m.GCSys)

	table.Show()
}
