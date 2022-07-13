package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	GOLANG_BIRTHDAY = "2006-01-02 15:04:05"
)

func RunDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("get run path error")
	}

	return strings.Replace(dir, "\\", "/", -1)
}

func IsFile(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func CurrentDir() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}
