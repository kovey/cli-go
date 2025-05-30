package app

import (
	"fmt"
	"os"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/cli-go/util"
	"github.com/kovey/debug-go/async"
	"github.com/kovey/debug-go/debug"
)

type ServInterface interface {
	Flag(AppInterface) error
	Init(AppInterface) error
	Run(AppInterface) error
	AsyncLog(AppInterface)
	Shutdown(AppInterface) error
	Reload(AppInterface) error
	Panic(AppInterface)
	Usage()
	PidFile(AppInterface) string
	Version() string
}

type ServBase struct {
}

func (s *ServBase) Version() string {
	return "0.0.0"
}

func (s *ServBase) AsyncLog(a AppInterface) {
	if os.Getenv(env.DEBUG_ASYNC_OPEN) != "On" {
		return
	}

	logDir := os.Getenv(env.LOG_DIR)
	if logDir == "" {
		return
	}

	maxLen, _ := env.GetInt(env.LOG_MAX)
	if maxLen == 0 {
		maxLen = 10240
	}

	if err := async.Start(logDir, maxLen); err != nil {
		debug.Erro("async log start error: %s", err)
		return
	}

	a.RunChild(func(ai AppInterface) {
		async.Listen()
	})
}

func (s *ServBase) PidFile(a AppInterface) string {
	path := os.Getenv(env.PID_FILE)
	if path != "" {
		return path
	}

	return fmt.Sprintf("%s/%s.run.pid", util.CurrentDir(), a.Name())
}

func (s *ServBase) Flag(AppInterface) error {
	debug.Info("run flag")
	return nil
}

func (s *ServBase) Reload(AppInterface) error {
	debug.Info("run reload")
	return nil
}

func (s *ServBase) Panic(a AppInterface) {
	debug.Info("app[%s] is panic", a.Name())
}

func (s *ServBase) Usage() {
	PrintDefaults()
}
