package main

import (
	"time"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
)

func main() {
	//testCallBack()
	testServ()
}

type serv struct {
	*app.ServBase
}

/**
func (s *serv) Usage() {
	fmt.Println(`
Usage:
	test <command> [arguments]
The commands are:
	test
Use "migrate help <command>" for more information about a command.
	`)
}
**/

func (s *serv) Flag(a app.AppInterface) error {
	a.FlagLong("to", "test", app.TYPE_STRING, "test")
	a.FlagNonValue("v", "show version")
	a.FlagNonValueLong("version", "show version")
	a.FlagLong("to-user", "user", app.TYPE_STRING, "user")
	return nil
}

func (s *serv) Init(a app.AppInterface) error {
	debug.Info("[%s] init", a.Name())
	return nil
}

func (s *serv) Run(a app.AppInterface) error {
	debug.Info("[%s] run", a.Name())
	if test, err := a.Get("to"); err == nil {
		debug.Info("test: %s", test.String())
	}
	if f, err := a.Arg(0, app.TYPE_STRING); err == nil {
		debug.Info("arg 0: %s", f.String())
	}
	if f, err := a.Arg(1, app.TYPE_STRING); err == nil {
		debug.Info("arg 1: %s", f.String())
	}
	if f, err := a.Get("to-user"); err == nil {
		debug.Info("to-user: %s", f.String())
	}

	time.Sleep(1 * time.Minute)
	//panic("run error")
	return nil
}

func (s *serv) Reload(a app.AppInterface) error {
	debug.Info("[%s] reload", a.Name())
	return nil
}

func (s *serv) Shutdown(a app.AppInterface) error {
	debug.Info("[%s] shutdown", a.Name())
	return nil
}

func testServ() {
	cli := app.NewApp("test")
	cli.SetDebugLevel(debug.Debug_Info)
	cli.SetServ(&serv{})
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}

func testCallBack() {
	cli := app.NewApp("test")
	cli.SetDebugLevel(debug.Debug_Info)
	cli.Action = func(a app.AppInterface) error {
		//panic("action is panic")
		debug.Info("app is running")
		debug.Warn("this is warning")
		debug.Erro("this is error")
		path, err := a.Get("config")
		if err != nil {
			return err
		}

		count, err := a.Get("count")
		if err != nil {
			return err
		}

		debug.Info("config[%s]", path.String())
		debug.Info("count[%d]", count.Int())
		return nil
	}

	cli.Reload = func(a app.AppInterface) error {
		debug.Info("app is reload")
		return nil
	}

	cli.Stop = func(a app.AppInterface) error {
		debug.Info("app is stop")
		return nil
	}
	cli.Panic = func(ai app.AppInterface) {
		debug.Info("app[%s] is panic", ai.Name())
	}

	cli.Flag("config", "", app.TYPE_STRING, "config path")
	cli.Flag("count", 100, app.TYPE_INT, "reload count")
	cli.Flag("count", 100, app.TYPE_INT, "reload count")
	cli.Flag("s", "", app.TYPE_STRING, "signal")
	err := cli.Run()
	if err != nil {
		debug.Erro(err.Error())
	}
}
