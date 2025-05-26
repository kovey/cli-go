package app

import (
	"fmt"
	"os"

	"github.com/kovey/cli-go/env"
	"github.com/kovey/cli-go/util"
	"github.com/kovey/debug-go/debug"
)

type ServInterface interface {
	Flag(AppInterface) error
	Init(AppInterface) error
	Run(AppInterface) error
	Shutdown(AppInterface) error
	Reload(AppInterface) error
	Panic(AppInterface)
	Usage()
	PidFile(AppInterface) string
}

type ServBase struct {
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
