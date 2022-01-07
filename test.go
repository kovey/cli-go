package main

import (
	"cli/app"
	"fmt"
)

func main() {
	cli := app.NewApp()
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
	cli.Add("config", "", app.TYPE_STRING, "config path")
	cli.Add("count", 0, app.TYPE_INT, "reload count")
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
