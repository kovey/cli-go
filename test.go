package main

import (
	"cli/app"
	"fmt"
)

func main() {
	cli := app.NewApp("test")
	cli.Action = func(a *app.App) error {
		fmt.Println("app is running")
		path, err := a.Get("config")
		if err != nil {
			return err
		}

		count, err := a.Get("count")
		if err != nil {
			return err
		}

		fmt.Printf("config[%s]\n", path.String())
		fmt.Printf("count[%d]\n", count.Int())
		return nil
	}

	cli.Reload = func(a *app.App) error {
		fmt.Println("app is reload")
		return nil
	}

	cli.Stop = func(a *app.App) error {
		fmt.Println("app is stop")
		return nil
	}

	cli.Flag("config", "", app.TYPE_STRING, "config path")
	cli.Flag("count", 100, app.TYPE_INT, "reload count")
	cli.Flag("s", "", app.TYPE_STRING, "signal")
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
