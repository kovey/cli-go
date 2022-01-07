package util

import (
	"fmt"
	"path"
	"runtime"
)

func RunDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("get path error")
	}

	fmt.Println(filename)
	return path.Dir(filename)
}
