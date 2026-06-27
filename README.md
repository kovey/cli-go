# cli-go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.23-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Test](https://github.com/kovey/cli-go/actions/workflows/test.yml/badge.svg)](https://github.com/kovey/cli-go/actions/workflows/test.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

A feature-rich CLI application framework for Go, with built-in support for hierarchical subcommands, short/long flags, daemon mode, signal handling, and hot-reload.

---

## Features

- **Hierarchical Subcommands** — nest commands arbitrarily deep with per-level flags
- **Short & Long Flags** — `-v` and `--verbose` style, with typed values (string, int, float, bool)
- **Daemon Mode** — background process management with PID files, auto-restart on crash
- **Signal Handling** — graceful shutdown (SIGTERM), hot-reload (SIGUSR1)
- **Built-in Commands** — `start`, `stop`, `restart`, `reload`, `kill` out of the box
- **`.env` File Loading** — parse environment files with comment support (#, ;, --, //)
- **Colorized Terminal Output** — table rendering with border-drawing characters and color support
- **Async Logging** — optional buffered async log writer

---

## Installation

```bash
go get -u github.com/kovey/cli-go
```

Requires **Go 1.23+**.

---

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/kovey/cli-go/app"
    "github.com/kovey/debug-go/debug"
)

type MyService struct {
    *app.ServBase
}

func (s *MyService) Flag(a app.AppInterface) error {
    a.FlagLong("name", "world", app.TYPE_STRING, "your name")
    a.FlagNonValue("v", "show version")
    return nil
}

func (s *MyService) Run(a app.AppInterface) error {
    name, _ := a.Get("name")
    fmt.Printf("Hello, %s!\n", name.String())
    return nil
}

func main() {
    cli := app.NewApp("hello")
    cli.SetDebugLevel(debug.Debug_Info)
    cli.SetServ(&MyService{})
    if err := cli.Run(); err != nil {
        debug.Erro("error: %s", err)
    }
}
```

```bash
$ go run main.go --name=Kovey
Hello, Kovey!
```

---

## API Reference

### Creating an Application

```go
cli := app.NewApp("myapp")       // create with name
cli.SetServ(&MyService{})        // attach service handler
cli.SetDebugLevel(debug.Debug_Info)
cli.Run()                        // parse args → Flag → Init → Run
```

### Service Interface

Implement `app.ServInterface` (embed `app.ServBase` for defaults):

| Method | Signature | Called When |
|--------|-----------|-------------|
| `Flag` | `func(AppInterface) error` | Before argument parsing — register flags here |
| `Init` | `func(AppInterface) error` | After parsing, before `Run` |
| `Run` | `func(AppInterface) error` | Main application logic |
| `Reload` | `func(AppInterface) error` | SIGUSR1 received (hot-reload) |
| `Panic` | `func(AppInterface)` | Panic recovered in Run |
| `AsyncLog` | `func(AppInterface)` | Start async log (if enabled) |
| `AsyncLogClose` | `func(AppInterface)` | Close async log |
| `PidFile` | `func(AppInterface) string` | Path to PID file |
| `Version` | `func() string` | Version string for `-v` |
| `Author` | `func() string` | Author name |
| `Usage` | `func()` | Custom usage/help output |

### Registering Flags

#### Value Flags (require a value argument)

| Method | CLI Form | Example |
|--------|----------|---------|
| `Flag(name, def, type, comment)` | `-name value` | `a.Flag("n", "default", app.TYPE_STRING, "...")` |
| `FlagLong(name, def, type, comment)` | `--name=value` | `a.FlagLong("name", "default", app.TYPE_STRING, "...")` |

#### Non-Value Flags (boolean presence, no value)

| Method | CLI Form | Example |
|--------|----------|---------|
| `FlagNonValue(name, comment)` | `-v` | `a.FlagNonValue("v", "verbose mode")` |
| `FlagNonValueLong(name, comment)` | `--verbose` | `a.FlagNonValueLong("verbose", "verbose mode")` |

#### Subcommands (positional arguments)

| Method | CLI Form | Example |
|--------|----------|---------|
| `FlagArg(name, comment)` | `command` | `a.FlagArg("create", "create resource")` |

#### Hierarchical Flags

Flags can be scoped under subcommands using variadic `parents`:

```go
// --path is only valid under the "create" subcommand
a.FlagLong("path", "/tmp", app.TYPE_STRING, "output path", "create")

// --to-user is only valid under "create build"
a.FlagLong("to-user", "admin", app.TYPE_STRING, "target user", "create", "build")
```

```bash
$ ./app create build --to-user=kovey
```

### Flag Types

```go
app.TYPE_STRING   // string value
app.TYPE_INT      // integer value
app.TYPE_BOOL     // boolean value (true/false/1/0)
app.TYPE_FLOAT    // float64 value
```

### Accessing Flag Values

```go
f, err := a.Get("name")          // root-level flag
f, err := a.Get("create", "path") // hierarchical flag

f.String()   // get as string
f.Int()      // get as int
f.Bool()     // get as bool
f.Float()    // get as float64
f.IsInput()  // true if flag was provided on command line
```

### Accessing Subcommands

```go
arg, err := a.Arg(0, app.TYPE_STRING)  // first positional argument
arg, err := a.Arg(1, app.TYPE_STRING)  // second positional argument
```

---

## Daemon Mode

cli-go provides a full daemon lifecycle out of the box:

```bash
# Start as daemon (forks to background)
./myapp start --daemon

# Hot-reload configuration (sends SIGUSR1)
./myapp reload

# Graceful stop (sends SIGTERM)
./myapp stop

# Force kill (sends SIGKILL)
./myapp kill

# Stop and restart as daemon
./myapp restart --daemon

# Show version
./myapp -v
```

The daemon process:
- Writes a PID file (`<appname>.run.pid`)
- Listens for signals and proxies them to the child process
- Auto-respawns the child on unexpected exit
- Reloads `.env` changes every second

---

## Environment Files (`.env`)

cli-go automatically loads a `.env` file if present and modified since last load.

### Supported Formats

```ini
# comment
; comment
-- comment
// comment

APP_NAME   = myapp
APP_URL    = http://example.com?a=1&b=2
APP_STATUS = 1
APP_OPEN   = true
APP_PRICE  = 10.05
```

### API

```go
import "github.com/kovey/cli-go/env"

env.Get("APP_NAME")           // → "myapp", error
env.GetInt("APP_STATUS")      // → 1, error
env.GetFloat("APP_PRICE")     // → 10.05, error
env.GetBool("APP_OPEN")       // → true, error

env.Load("./config.env")      // load a specific file
env.LoadDefault(time.Now())   // load .env with timestamp tracking
env.HasEnv()                  // true if .env exists and is newer than last load
```

### Environment Variables for Configuration

| Variable | Purpose | Default |
|----------|---------|---------|
| `APP_NAME` | Application name (fallback) | — |
| `DEBUG_LEVEL` | Log level | — |
| `DEBUG_SHOW_FILE` | Show file:line in logs | — |
| `DEBUG_ASYNC_OPEN` | Enable async logging (`"On"`) | off |
| `LOG_DIR` | Async log output directory | — |
| `LOG_MAX` | Max log buffer size | 10240 |
| `PID_FILE` | Custom PID file path | `<cwd>/<appname>.run.pid` |

---

## Terminal Output

### Table Rendering

```go
import "github.com/kovey/cli-go/gui"

table := gui.NewTable()
table.Add(0, "ID")
table.Add(0, "Name")
table.Add(1, "1")
table.Add(1, "Kovey")
table.Show()
```

```
┌──┬─────┐
│ID│Name │
├──┼─────┤
│1 │Kovey│
└──┴─────┘
```

### Colorized Status

```go
gui.PrintlnOk("build completed")       // green [ok]
gui.PrintlnFailure("build failed")     // red [failure]
gui.PrintlnNormal("left", "right")     // left-aligned with spacing
```

---

## Project Structure

```
cli-go/
├── app/            # Core framework: App, Daemon, Flag, CommandLine, Help
├── env/            # .env file parser and accessors
├── gui/            # Terminal table rendering and color output
├── util/           # File system and runtime utilities
├── examples/
│   ├── 01-basic/   # Basic CLI with flags and subcommands
│   └── 02-daemon/  # Daemon mode with start/stop/reload
└── .github/
    └── workflows/  # CI pipeline (go vet + go test)
```

## Running the Examples

```bash
# Basic CLI
cd examples/01-basic
go run main.go --name=myapp --count=5 --verbose
go run main.go create --path=/tmp/config --type=json
go run main.go create help

# Daemon mode
cd examples/02-daemon
go build -o mydaemon
./mydaemon start --daemon
./mydaemon reload
./mydaemon stop
```

---

## License

Apache 2.0 — see [LICENSE](LICENSE)
