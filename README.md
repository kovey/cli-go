# kovey cli of terminal by golang
### Description
#### This is a cli app library
### Usage
    go get -u github.com/kovey/cli-go
### Examples
```golang
    package main

    import (
        "github.com/kovey/cli-go/app"
        "fmt"
    )

    func main() {
        cli := app.NewApp("sample example")
        cli.Action = func(a *app.App) error {
            fmt.Println("app is running")
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

            fmt.Printf("c is [%s]\n", path.String())
            fmt.Printf("i is [%s]\n", i.Int())
            fmt.Printf("b is [%t]\n", b.Bool())
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
   go run main.go -c test -i 100 -b false
```
