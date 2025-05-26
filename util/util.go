package util

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	Unit_Second int64 = 1
	Unit_Minute int64 = 60 * Unit_Second
	Unit_Hour   int64 = 60 * Unit_Minute
	Unit_Day    int64 = 24 * Unit_Hour
)

func ExecFilePath() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return path
}

func RunDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("get run path error")
	}

	return dir
}

func IsFile(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}

	return os.IsExist(err)
}

func CurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path
}

func IsRunWithGoRunCmd() bool {
	return strings.Contains(RunDir(), os.TempDir())
}
