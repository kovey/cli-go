package env

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	trim_set = " \r\n\t"
)

var Err_Key_Not_Found = errors.New("key not found")

func Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if err := parse(line); err != nil {
				return err
			}
			break
		}

		if err != nil {
			return err
		}

		if err := parse(line); err != nil {
			return err
		}
	}

	return nil
}

func parse(line string) error {
	line = strings.Trim(line, trim_set)
	if len(line) == 0 {
		return nil
	}

	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "--") || strings.HasPrefix(line, "//") {
		return nil
	}

	info := strings.Split(line, "=")
	itemCount := len(info)
	if itemCount < 2 {
		return fmt.Errorf("%s format error", line)
	}

	key := strings.Trim(info[0], trim_set)
	value := strings.Trim(info[1], trim_set)
	if itemCount > 2 {
		info[1] = strings.TrimLeft(info[1], trim_set)
		info[itemCount-1] = strings.TrimRight(info[itemCount-1], trim_set)
		value = strings.Join(info[1:], "=")
	}

	return os.Setenv(key, value)
}

func GetFloat(key string) (float64, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, Err_Key_Not_Found
	}

	return strconv.ParseFloat(value, 64)
}

func GetBool(key string) (bool, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false, Err_Key_Not_Found
	}

	return strconv.ParseBool(value)
}

func GetInt(key string) (int, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return 0, Err_Key_Not_Found
	}

	return strconv.Atoi(value)
}

func Get(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", Err_Key_Not_Found
	}

	return value, nil
}
