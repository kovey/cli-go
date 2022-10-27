package main

import (
	"github.com/kovey/cli-go/app"
	"github.com/kovey/cli-go/debug"
)

func main() {
	cli := app.NewApp("test")
	cli.Action = func(a *app.App) error {
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

	cli.Reload = func(a *app.App) error {
		debug.Info("app is reload")
		return nil
	}

	cli.Stop = func(a *app.App) error {
		debug.Info("app is stop")
		return nil
	}

	cli.Flag("config", "", app.TYPE_STRING, "config path")
	cli.Flag("count", 100, app.TYPE_INT, "reload count")
	cli.Flag("count", 100, app.TYPE_INT, "reload count")
	cli.Flag("s", "", app.TYPE_STRING, "signal")
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
