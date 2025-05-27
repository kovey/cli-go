# kovey cli of terminal by golang
### Description
#### This is a cli app library with golang
### Usage
    go get -u github.com/kovey/cli-go
### Examples
```golang
    package main

    import (
        "time"

        "github.com/kovey/cli-go/app"
        "github.com/kovey/debug-go/debug"
    )

    func main() {
        testServ()
    }

    type serv struct {
        *app.ServBase
    }

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
        cli := app.NewApp("test")
        cli.SetDebugLevel(debug.Debug_Info)
        cli.SetServ(&serv{})
        if err := cli.Run(); err != nil {
            debug.Erro(err.Error())
        }
    }


```
```bash
   go run test.go start # run app
   go run test.go start --daemon # run app with daemon mode
   go run test.go reload # reload app
   go run test.go stop # stop app
   go run test.go restart # restart app
   go run test.go restart --daemon # restart app and run app with daemon mode
   go run test.go kill # kill app
```
### Env
    .env file
#### Conent of .env file
```ini
; comment
-- comment
// comment
# comment

APP_NAME    = Test
APP_URL     = http://baidu.com?ab=test&cd=dev
APP_STATUS  = 1
APP_OPEN    = true
APP_PRICE   = 10.05
```

#### Get
```golang
    env.Get("APP_NAME")
    env.GetInt("APP_STATUS")
    env.GetFloat("APP_PRICE")
    env.GetBool("APP_OPEN")
```
