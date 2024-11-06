package app

import (
	"fmt"
	"os"

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
}

type ServBase struct {
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
	fmt.Printf("Usage of %s \n", os.Args[0])
	PrintDefaults()
}
