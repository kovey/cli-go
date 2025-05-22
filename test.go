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
	a.FlagArg("create", "create config")
	a.FlagArg("build", "build config", "create")
	a.FlagArg("make", "make config", "create", "build")
	a.FlagLong("to", nil, app.TYPE_STRING, "test", "create", "build")
	a.FlagLong("from", nil, app.TYPE_STRING, "from path", "create", "build")
	a.FlagNonValue("v", "show version", "create")
	a.FlagNonValueLong("version", "show version", "create", "build")
	a.FlagLong("to-user", "user", app.TYPE_STRING, "user", "create", "build", "make")
	a.FlagLong("path", nil, app.TYPE_STRING, "config path")
	a.FlagLong("name", nil, app.TYPE_STRING, "name")
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

	time.Sleep(30 * time.Second)
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
	/**
	file, err := os.OpenFile("./test.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	debug.SetWriter(file)
	**/
	cli := app.NewApp("test")
	cli.SetDebugLevel(debug.Debug_Info)
	cli.SetServ(&serv{})
	if err := cli.Run(); err != nil {
		debug.Erro(err.Error())
	}
}
