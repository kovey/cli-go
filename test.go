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

		fmt.Printf("config[%s]\n", *(path.(*string)))
		fmt.Printf("count[%d]\n", *(count.(*int)))
		fmt.Printf("pid[%d]\n", a.Pid())
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
	cli.Flag("count", 0, app.TYPE_INT, "reload count")
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
