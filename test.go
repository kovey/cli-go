package main

import (
	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/gui"
	"github.com/kovey/debug-go/debug"
)

func main() {
	//testCallBack()
	testServ()
}

type serv struct {
	*app.ServBase
}

func (s *serv) Init(a app.AppInterface) error {
	debug.Info("[%s] init", a.Name())
	return nil
}

func (s *serv) Run(a app.AppInterface) error {
	debug.Info("[%s] run", a.Name())
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
	cli.Show = func(table *gui.Table) {
		table.Add("custom")
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
