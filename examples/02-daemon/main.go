// Example: Daemon Application
//
// This example demonstrates daemon mode usage of the cli-go framework:
//   - Running an app as a background daemon
//   - Signal handling (SIGUSR1 for reload, SIGTERM for graceful shutdown)
//   - PID file management
//   - Subcommands: start, stop, reload, restart, kill
//   - Async logging
//
// Build and run:
//
//	cd examples/02-daemon && go build -o mydaemon
//
// Start daemon:
//
//	./mydaemon start --daemon
//
// Check status:
//
//	cat mydaemon.run.pid
//
// Reload config:
//
//	./mydaemon reload
//
// Stop gracefully:
//
//	./mydaemon stop
//
// Restart:
//
//	./mydaemon restart --daemon
//
// Kill (force):
//
//	./mydaemon kill
package main

import (
	"fmt"
	"time"

	"github.com/kovey/cli-go/app"
	"github.com/kovey/debug-go/debug"
)

// DaemonService implements a service that runs as a daemon.
type DaemonService struct {
	*app.ServBase
	counter int
}

// Flag registers flags for the daemon service.
func (s *DaemonService) Flag(a app.AppInterface) error {
	// Custom flags for the service
	a.FlagLong("interval", 5, app.TYPE_INT, "work interval in seconds")
	a.FlagLong("max-retry", 3, app.TYPE_INT, "max retry count")

	return nil
}

// Init initializes the service.
func (s *DaemonService) Init(a app.AppInterface) error {
	debug.Info("[%s] daemon initializing...", a.Name())

	if interval, err := a.Get("interval"); err == nil {
		debug.Info("Work interval: %d seconds", interval.Int())
	}

	return nil
}

// Run is the main daemon loop.
func (s *DaemonService) Run(a app.AppInterface) error {
	debug.Info("[%s] daemon started, pid=%d", a.Name(), a.(*app.App).Pid())

	// Simulate a long-running daemon loop
	for {
		select {
		case <-a.Context().Done():
			debug.Info("[%s] shutting down gracefully...", a.Name())
			return nil
		default:
			s.counter++
			debug.Info("[%s] heartbeat #%d", a.Name(), s.counter)
			time.Sleep(time.Duration(s.counter) * time.Second)
		}
	}
}

// Reload handles the SIGUSR1 signal for hot-reload.
func (s *DaemonService) Reload(a app.AppInterface) error {
	debug.Info("[%s] reloading configuration...", a.Name())
	fmt.Printf("Configuration reloaded at %s\n", time.Now().Format(time.DateTime))
	return nil
}

// Panic handles panic recovery in the daemon.
func (s *DaemonService) Panic(a app.AppInterface) {
	debug.Erro("[%s] panic recovered, restarting...", a.Name())
}

func main() {
	cli := app.NewApp("mydaemon")
	cli.SetDebugLevel(debug.Debug_Info)
	debug.SetFileLine(debug.File_Line_Show)
	cli.SetServ(&DaemonService{})

	if err := cli.Run(); err != nil {
		debug.Erro("daemon error: %s", err)
	}
}
