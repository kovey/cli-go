package app

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/cli-go/util"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/debug-go/run"
)

const (
	Ko_Cli_Daemon_Background = "KO_CLI_DAEMON_BACKGROUND"
	ko_command_start         = "start"
	ko_command_reload        = "reload"
	ko_command_stop          = "stop"
)

type Daemon struct {
	pid          int
	args         []string
	childPid     int
	cmd          *exec.Cmd
	wait         sync.WaitGroup
	sig          chan os.Signal
	openChild    chan bool
	isBackground bool
	pidFile      string
	serv         ServInterface
	name         string
	showUsage    bool
	check        *time.Ticker
}

func NewDaemon(name string) *Daemon {
	if len(name) == 0 {
		name = os.Getenv(env.APP_NAME)
		if len(name) == 0 {
			panic("app name is empty")
		}
	}

	d := &Daemon{name: name, wait: sync.WaitGroup{}, sig: make(chan os.Signal, 1), openChild: make(chan bool, 1), isBackground: false, check: time.NewTicker(1 * time.Second)}
	if ok, err := strconv.ParseBool(os.Getenv(Ko_Cli_Daemon_Background)); err == nil {
		d.isBackground = ok
	}
	for _, arg := range os.Args {
		if arg == "--daemon" && d.isBackground {
			continue
		}

		d.args = append(d.args, arg)
	}

	if dbl, err := env.Get(env.DEBUG_LEVEL); err == nil && len(dbl) > 0 {
		d.SetDebugLevel(debug.DebugType(dbl))
	}
	if showFile, err := env.GetInt(env.DEBUG_SHOW_FILE); err == nil {
		debug.SetFileLine(debug.FileLine(showFile))
	}

	_commanLine.FlagArg("start", fmt.Sprintf("start app[%s]", name), 0)
	_commanLine.FlagArg("reload", fmt.Sprintf("reload app[%s]", name), 0)
	_commanLine.FlagArg("stop", fmt.Sprintf("stop app[%s]", name), 0)
	_commanLine.FlagNonValueLong("daemon", fmt.Sprintf("run app[%s] with daemon mode", name))
	return d
}

func (d *Daemon) UsageWhenErr() {
	d.showUsage = true
}

func (d *Daemon) FlagNonValueLong(name string, comment string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagNonValueLong(name, comment)
}

func (d *Daemon) FlagNonValue(name string, comment string) {
	if name == "s" {
		debug.Warn("flag[%s] is used by sinal module", name)
		return
	}

	_commanLine.FlagNonValue(name, comment)
}

func (d *Daemon) FlagLong(name string, def any, t Type, comment string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagLong(name, def, t, comment)
}

func (d *Daemon) Flag(name string, def any, t Type, comment string) {
	if name == "s" {
		debug.Warn("flag[%s] is used by sinal module", name)
		return
	}

	_commanLine.Flag(name, def, t, comment)
}

func (d *Daemon) Get(name string) (*Flag, error) {
	f := _commanLine.Get(name)
	if f == nil {
		return nil, fmt.Errorf("[%s] is not exists", name)
	}

	return f, nil
}

func (d *Daemon) Arg(index int, t Type) (*Flag, error) {
	f := _commanLine.Arg(index)
	if f == nil {
		return nil, fmt.Errorf("[%d] is not exists", index)
	}

	f.t = t
	return f, nil
}

func (d *Daemon) Pid() int {
	return d.pid
}

func (d *Daemon) PidString() string {
	return strconv.Itoa(d.pid)
}

func (d *Daemon) getPid() int {
	d.pidFile = d.serv.PidFile(d)
	pid, err := os.ReadFile(d.pidFile)
	if err != nil {
		return -1
	}

	id, err := strconv.Atoi(string(pid))
	if err != nil {
		return -1
	}

	return id
}

func (d *Daemon) SetServ(serv ServInterface) {
	d.serv = serv
	Usage = serv.Usage
}

func (d *Daemon) Name() string {
	return d.name
}

func (d *Daemon) SetDebugLevel(t debug.DebugType) {
	debug.SetLevel(t)
}

func (d *Daemon) runChild() {
	defer d.wait.Done()
	d.childPid = -1
	if err := d.doRun(); err != nil {
		d.cmd = nil
		debug.Erro("run child error: %s", err)
		return
	}

	if err := d.cmd.Wait(); err != nil {
		debug.Erro("wait child error: %s", err)
		return
	}

	d.openChild <- true
}

func (d *Daemon) listen() {
	signal.Notify(d.sig, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
	defer d.check.Stop()

	for {
		select {
		case now := <-d.check.C:
			loadEnv(now)
		case s := <-d.sig:
			if d.childPid > 0 {
				if err := syscall.Kill(d.childPid, s.(syscall.Signal)); err != nil {
					debug.Erro("sinal to app[%s] child[%d] failure, error: %s", d.name, d.childPid, err)
				}
				switch s {
				case os.Interrupt, syscall.SIGTERM:
					return
				}
			}

			if d.serv != nil {
				switch s {
				case os.Interrupt, syscall.SIGTERM:
					if err := d.serv.Shutdown(d); err != nil {
						debug.Erro("app[%s] stop failure, error: %s", d.name, err)
					}
					return
				case syscall.SIGUSR1:
					if err := d.serv.Reload(d); err != nil {
						debug.Erro("app[%s] reload failure, error: %s", d.name, err)
					}
				}
			}
		case <-d.openChild:
			d.wait.Add(1)
			go d.runChild()
		}
	}
}

func (d *Daemon) runApp() error {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		d.serv.Panic(d)
		run.Panic(err)
	}()

	if err := d.serv.Init(d); err != nil {
		return fmt.Errorf("run app[%s] init error: %s", d.name, err)
	}

	if err := d.serv.Run(d); err != nil {
		if d.showUsage {
			d.serv.Usage()
		}

		return fmt.Errorf("run app[%s] error: %s", d.name, err)
	}

	return nil
}

func (d *Daemon) _run() error {
	if f := _commanLine.Get("daemon"); f == nil || !f.has {
		return d.runApp()
	}

	if !d.isBackground {
		if err := d.doRun(); err != nil {
			return fmt.Errorf("run background process error: %s", err)
		}

		os.Exit(0)
	}

	d.pidFile = d.serv.PidFile(d)
	if util.IsFile(d.pidFile) {
		return fmt.Errorf("app[%s] is running", d.name)
	}

	defer func() {
		if d.pidFile != "" {
			os.Remove(d.pidFile)
		}

		if err := recover(); err != nil {
			if d.childPid > 0 {
				if err := syscall.Kill(d.childPid, syscall.SIGTERM); err != nil {
					debug.Erro("stop child[%d] failure, error: %s", d.childPid, err)
				}
			}

			run.Panic(err)
			d.serv.Panic(d)
		}
	}()

	if err := os.WriteFile(d.pidFile, []byte(d.PidString()), 0644); err != nil {
		debug.Erro("write pid file error: %s", err)
	}

	debug.Info("app[%s] run, pid[%s]", d.name, d.PidString())
	d.wait.Add(1)
	go d.runChild()
	d.listen()
	d.wait.Wait()
	return nil
}

func (d *Daemon) _reload() error {
	pid := d.getPid()
	if pid < 1 {
		return fmt.Errorf("app[%s] not running", d.name)
	}

	return syscall.Kill(pid, syscall.SIGUSR1)
}

func (d *Daemon) _stop() error {
	pid := d.getPid()
	if pid < 1 {
		return fmt.Errorf("app[%s] not running", d.name)
	}

	return syscall.Kill(pid, syscall.SIGTERM)
}

func (d *Daemon) Run() error {
	if d.serv == nil {
		return fmt.Errorf("serv not init")
	}

	d.pid = os.Getpid()
	if err := d.serv.Flag(d); err != nil {
		return fmt.Errorf("run app flag error: %s", err)
	}

	_commanLine.Parse(os.Args[1:])
	f, err := d.Arg(0, TYPE_STRING)
	if err != nil {
		f, _ = d.Get("start")
	}

	switch f.String() {
	case ko_command_start:
		return d._run()
	case ko_command_reload:
		return d._reload()
	case ko_command_stop:
		return d._stop()
	}

	return nil
}

func (d *Daemon) doRun() error {
	env := append(os.Environ(), fmt.Sprintf("%s=%t", Ko_Cli_Daemon_Background, true))
	d.cmd = &exec.Cmd{Path: d.args[0], Args: d.args, SysProcAttr: &syscall.SysProcAttr{Setsid: true}, Env: env}
	err := d.cmd.Start()
	if err == nil {
		d.childPid = d.cmd.Process.Pid
	}
	return err
}
