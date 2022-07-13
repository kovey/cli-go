# kovey cli of terminal by golang
### Description
#### This is a cli app library
#### Default flag "s" is signal, contains reload, stop and maintain
#### app will call Stop callback when recevie stop signal
#### app will call Reload callback when recevie reload signal
#### app will call Maintain callback when recevie maintain signal
### Usage
    go get -u github.com/kovey/cli-go
### Examples
```golang
    package main

    import (
        "github.com/kovey/cli-go/app"
        "github.com/kovey/cli-go/debug"
    )

    func main() {
        cli := app.NewApp("sample example")
        cli.Action = func(a *app.App) error {
            debug.Info("app is running")
            path, err := a.Get("c")
            if err != nil {
                return err
            }

            i, err := a.Get("i")
            if err != nil {
                return err
            }

            b, err := a.Get("b")
            if err != nil {
                return err
            }

            debug.Info("c is [%s]", path.String())
            debug.Info("i is [%s]", i.Int())
            debug.Info("b is [%t]", b.Bool())
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

        cli.Flag("c", "", app.TYPE_STRING, "app config path, type string")
        cli.Flag("i", 0, app.TYPE_INT, "type int config")
        cli.Flag("b", false, app.TYPE_BOOL, "type bool config")

        err := cli.Run()
        if err != nil {
            panic(err)
        }
    }

```
```bash
   go run main.go -c test -i 100 -b false # run app
   go run main.go -s reload # reload app
   go run main.go -s stop # stop app
   go run main.go -s maintain # maintain app
```
