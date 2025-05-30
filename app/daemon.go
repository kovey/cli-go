package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/cli-go/gui"
	"github.com/kovey/cli-go/util"
	"github.com/kovey/debug-go/async"
	"github.com/kovey/debug-go/debug"
	"github.com/kovey/debug-go/run"
)

const (
	Ko_Cli_Daemon_Background = "KO_CLI_DAEMON_BACKGROUND"
	Ko_Command_Start         = "start"
	Ko_Command_Reload        = "reload"
	Ko_Command_Stop          = "stop"
	Ko_Command_Kill          = "kill"
	Ko_Command_Restart       = "restart"
	Ko_Command_Daemon        = "daemon"
	ko_command_daemon_arg    = "--daemon"
	Ko_Command_Help          = "help"
	ko_command_help_arg      = "h"
)

var Err_App_Init = errors.New("app init failure")
var Err_App_Run = errors.New("app run failure")
var Err_Not_Restart = errors.New("app not restart")
var Err_App_Process_Exit_Unexpected = errors.New("app run process exit unexpected")

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
	workdir      string
	internalSig  chan bool
	childRunErr  error
}

func NewDaemon(name string) *Daemon {
	if len(name) == 0 {
		name = os.Getenv(env.APP_NAME)
		if len(name) == 0 {
			panic("app name is empty")
		}
	}

	d := &Daemon{name: name, wait: sync.WaitGroup{}, sig: make(chan os.Signal, 1), openChild: make(chan bool, 1), isBackground: false, check: time.NewTicker(1 * time.Second)}
	d.internalSig = make(chan bool, 1)
	if ok, err := strconv.ParseBool(os.Getenv(Ko_Cli_Daemon_Background)); err == nil {
		d.isBackground = ok
		gui.Background()
	}
	for _, arg := range os.Args {
		if arg == ko_command_daemon_arg && d.isBackground {
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

	_commanLine.help.AppName = name
	_commanLine.help.Title = fmt.Sprintf("command line of %s", name)
	_commanLine.FlagArg(Ko_Command_Start, fmt.Sprintf("start app[%s]", name), 0)
	_commanLine.FlagArg(Ko_Command_Reload, fmt.Sprintf("reload app[%s]", name), 0)
	_commanLine.FlagArg(Ko_Command_Stop, fmt.Sprintf("stop app[%s]", name), 0)
	_commanLine.FlagArg(Ko_Command_Kill, fmt.Sprintf("kill app[%s] with -9", name), 0)
	_commanLine.FlagArg(Ko_Command_Restart, fmt.Sprintf("restart app[%s]", name), 0)
	_commanLine.FlagNonValueLong(Ko_Command_Daemon, fmt.Sprintf("start app[%s] with daemon mode", name), Ko_Command_Start)
	_commanLine.FlagNonValueLong(Ko_Command_Daemon, fmt.Sprintf("restart app[%s] and runned with daemon mode", name), Ko_Command_Restart)
	d._help()
	return d
}

func (d *Daemon) _help() {
	_commanLine.FlagNonValue(ko_command_help_arg, fmt.Sprintf("show app[%s] usage details", _commanLine.help.AppName))
	_commanLine.FlagArg(Ko_Command_Help, fmt.Sprintf("show app[%s] command usage details", _commanLine.help.AppName), 0)
	_commanLine.help.Args.Add(Ko_Command_Help, fmt.Sprintf("show app[%s] usage details", _commanLine.help.AppName), false, "")
}

func (d *Daemon) CleanCommandLine(isCleanHelp bool) {
	_commanLine.CleanDefaults()
	if isCleanHelp {
		return
	}

	d._help()
}

func (d *Daemon) UsageWhenErr() {
	d.showUsage = true
}

func (d *Daemon) FlagArg(name string, comment string, parents ...string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagArg(name, comment, 0, parents...)
}

func (d *Daemon) FlagNonValueLong(name string, comment string, parents ...string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagNonValueLong(name, comment, parents...)
}

func (d *Daemon) FlagNonValue(name string, comment string, parents ...string) {
	_commanLine.FlagNonValue(name, comment, parents...)
}

func (d *Daemon) FlagLong(name string, def any, t Type, comment string, parents ...string) {
	if len(name) < 2 {
		debug.Warn("flag[%s] is too short", name)
		return
	}

	_commanLine.FlagLong(name, def, t, comment, parents...)
}

func (d *Daemon) Flag(name string, def any, t Type, comment string, parents ...string) {
	_commanLine.Flag(name, def, t, comment, parents...)
}

func (d *Daemon) Get(names ...string) (*Flag, error) {
	f := _commanLine.Get(names...)
	if f == nil {
		return nil, fmt.Errorf("[%s] is not exists", strings.Join(names, "->"))
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
	return fmt.Sprintf("%d-%d", d.pid, d.childPid)
}

func (d *Daemon) getPid() int {
	d.pidFile = d.serv.PidFile(d)
	pid, err := os.ReadFile(d.pidFile)
	if err != nil {
		return -1
	}

	pidInfo := strings.Split(string(pid), "-")
	id, err := strconv.Atoi(pidInfo[0])
	if err != nil {
		return -1
	}

	return id
}

func (d *Daemon) getPidAndChildPid() []int {
	d.pidFile = d.serv.PidFile(d)
	pid, err := os.ReadFile(d.pidFile)
	if err != nil {
		return nil
	}

	pidInfo := strings.Split(string(pid), "-")
	ids := make([]int, len(pidInfo))
	for index, idInfo := range pidInfo {
		id, err := strconv.Atoi(idInfo)
		if err != nil {
			return nil
		}
		ids[index] = id
	}

	return ids
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

	if err := os.WriteFile(d.pidFile, []byte(d.PidString()), 0644); err != nil {
		debug.Erro("write pid file error: %s", err)
	}

	if err := d.cmd.Wait(); err != nil {
		debug.Erro("wait child error: %s", err)
		d.internalSig <- true
		if d.childRunErr == nil {
			d.childRunErr = Err_App_Process_Exit_Unexpected
		}
		return
	}

	d.openChild <- true
}

func (d *Daemon) listen() {
	if !d.isBackground {
		return
	}

	signal.Notify(d.sig, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR2, syscall.SIGUSR1)
	defer d.check.Stop()

	for {
		select {
		case <-d.internalSig:
			return
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

			if d.childPid == 0 && d.serv != nil {
				switch s {
				case os.Interrupt, syscall.SIGTERM:
					if err := d.serv.Shutdown(d); err != nil {
						debug.Erro("app[%s] stop failure, error: %s", d.name, err)
					}
					return
				case syscall.SIGUSR1:
					if err := d.serv.Reload(d); err != nil {
						debug.Erro("app[%s] reload failure, error: %s", d.name, err)
						gui.PrintlnFailure("app[%s] reloaded", d.name)
					} else {
						gui.PrintlnOk("app[%s] reloaded", d.name)
					}
				}
			}
		case <-d.openChild:
			d.wait.Add(1)
			go d.runChild()
		}
	}
}

func (d *Daemon) _runAppEnd() {
	if !d.isBackground {
		return
	}

	d.internalSig <- true
}

func (d *Daemon) _runApp() {
	defer d.wait.Done()
	defer d._runAppEnd()
	defer async.Stop()
	defer func() {
		err := recover()
		if err == nil {
			debug.Info("app[%s] total run time: [%s]", d.name, GetFormatRunTime())
			return
		}

		debug.Info("app[%s] total run time: [%s]", d.name, GetFormatRunTime())
		d.serv.Panic(d)
		run.Panic(err)
	}()

	startTime = time.Now()
	d.SetDebugLevel(debug.DebugType(os.Getenv(env.DEBUG_LEVEL)))
	if line, err := env.GetInt(os.Getenv(env.DEBUG_SHOW_FILE)); err == nil {
		debug.SetFileLine(debug.FileLine(line))
	}
	d.serv.AsyncLog(d)

	if err := d.serv.Init(d); err != nil {
		debug.Erro("run app[%s] init error: %s", d.name, err)
		d.childRunErr = Err_App_Init
		return
	}

	if err := d.serv.Run(d); err != nil {
		if d.showUsage {
			d.serv.Usage()
		}

		d.childRunErr = err
		debug.Erro("run app[%s] error: %s", d.name, err)
	}
}

func (d *Daemon) runApp() error {
	d.wait.Add(1)
	go d._runApp()

	d.listen()
	d.wait.Wait()
	return d.childRunErr
}

func (d *Daemon) _runDaemon() error {
	if util.IsRunWithGoRunCmd() {
		debug.Erro("daemon is unsupport when app run with go run command")
		gui.PrintlnFailure("app[%s] started", d.name)
		os.Exit(0)
		return nil
	}

	d.pidFile = d.serv.PidFile(d)
	if util.IsFile(d.pidFile) {
		debug.Erro("app[%s] is running", d.name)
		gui.PrintlnFailure("app[%s] started", d.name)
		os.Exit(0)
		return nil
	}

	if err := d.doRun(); err != nil {
		debug.Erro("run background process error: %s", err)
		gui.PrintlnFailure("app[%s] started", d.name)
		os.Exit(0)
		return nil
	}

	go func() {
		if err := d.cmd.Wait(); err != nil {
			gui.PrintlnFailure("app[%s] started", d.name)
			os.Exit(0)
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	<-ticker.C
	fmt.Print(".")
	<-ticker.C
	fmt.Println(".")
	ticker.Stop()
	gui.PrintlnOk("pid[%d] of app[%s] started", d.childPid, d.name)
	os.Exit(0)
	return nil
}

func (d *Daemon) _run(commands ...string) error {
	if f := _commanLine.Get(append(commands, Ko_Command_Daemon)...); f == nil || !f.has {
		err := d.runApp()
		if d.isBackground {
			if err == Err_App_Init || err == Err_Not_Restart {
				os.Exit(1)
			}
		}
		return err
	}

	if !d.isBackground {
		return d._runDaemon()
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

	d.wait.Add(1)
	go d.runChild()
	d.listen()
	d.wait.Wait()
	return d.childRunErr
}

func (d *Daemon) _reload() error {
	pid := d.getPid()
	if pid < 1 {
		debug.Erro("app[%s] not running", d.name)
		gui.PrintlnFailure("app[%s] reloaded", d.name)
		return nil
	}

	if err := syscall.Kill(pid, syscall.SIGUSR1); err != nil {
		debug.Erro("app[%s] reloaded failure: %s", d.name, err)
		gui.PrintlnFailure("app[%s] reloaded", d.name)
		return nil
	}

	gui.PrintlnOk("app[%s] reloaded", d.name)
	return nil
}

func (d *Daemon) _stop() error {
	pid := d.getPid()
	if pid < 1 {
		debug.Erro("app[%s] not running", d.name)
		gui.PrintlnFailure("app[%s] stopped", d.name)
		return nil
	}

	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		debug.Erro("app[%s] stopped failure: %s", d.name, err)
		gui.PrintlnFailure("pid[%d] of app[%s] stopped", pid, d.name)
		return nil
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	count := 0
	for {
		<-ticker.C
		count++
		if count%10 == 0 {
			fmt.Print(".")
		}
		if err := syscall.Kill(pid, 0); err != nil {
			fmt.Println(".")
			break
		}
	}

	gui.PrintlnOk("pid[%d] of app[%s] stopped", pid, d.name)
	return nil
}

func (d *Daemon) _restart() error {
	if err := d._stop(); err != nil {
		debug.Erro(err.Error())
	}

	d.args[1] = Ko_Command_Start
	return d._run(Ko_Command_Restart)
}

func (d *Daemon) _kill() error {
	pids := d.getPidAndChildPid()
	if len(pids) < 1 {
		debug.Erro("app[%s] not running", d.name)
		gui.PrintlnFailure("app[%s] killed", d.name)
		return nil
	}

	for _, pid := range pids {
		if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
			debug.Erro("pid[%d] of app[%s] killed failure", pid, d.name)
		}
	}

	gui.PrintlnFailure("app[%s] killed", d.name)
	return nil
}

func (d *Daemon) _runCommand(command string) error {
	switch command {
	case Ko_Command_Start:
		return d._run(command)
	case Ko_Command_Reload:
		return d._reload()
	case Ko_Command_Stop:
		return d._stop()
	case Ko_Command_Kill:
		return d._kill()
	case Ko_Command_Restart:
		return d._restart()
	case Ko_Command_Help:
		_commanLine.Help()
		return nil
	default:
		return d._run(_commanLine.AllArgName()...)
	}
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
	method := Ko_Command_Start
	if _commanLine.hasHelp() {
		method = Ko_Command_Help
	} else {
		if f, err := d.Arg(0, TYPE_STRING); err == nil {
			method = f.String()
		}
	}

	err := d._runCommand(method)
	if err == Err_App_Process_Exit_Unexpected {
		os.Exit(1)
	}

	return err
}

func (d *Daemon) doRun() error {
	env := append(os.Environ(), fmt.Sprintf("%s=%t", Ko_Cli_Daemon_Background, true))
	d.cmd = &exec.Cmd{Path: util.ExecFilePath(), Args: d.args, SysProcAttr: &syscall.SysProcAttr{Setsid: true}, Env: env, Dir: util.CurrentDir(), Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
	err := d.cmd.Start()
	if err == nil {
		d.childPid = d.cmd.Process.Pid
	}
	return err
}

func (d *Daemon) _runChild(call func(a AppInterface)) {
	defer d.wait.Done()
	call(d)
}

func (d *Daemon) RunChild(call func(a AppInterface)) {
	d.wait.Add(1)
	go d._runChild(call)
}
