package app

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kovey/debug-go/debug"
)

var Usage = func() {
	fmt.Printf("Usage of %s \n", os.Args[0])
	PrintDefaults()
}

var _commanLine = NewCommandLine()

func PrintDefaults() {
	_commanLine.PrintDefaults()
}

type CommandLine struct {
	flags    map[string]*Flag
	others   []*Flag
	args     []string
	isParsed bool
	keys     []string
	hasFlag  bool
}

func NewCommandLine() *CommandLine {
	return &CommandLine{flags: make(map[string]*Flag)}
}

func (c *CommandLine) PrintDefaults() {
	maxLen := 0
	for _, flag := range c.flags {
		fLen := len(flag.name) + 1
		if !flag.isShort {
			fLen++
		}

		if fLen > maxLen {
			maxLen = fLen
		}
	}

	for _, key := range c.keys {
		c.flags[key].print(maxLen)
	}
}

func (c *CommandLine) FlagArg(name, comment string, index int) {
	if err := c.checkLong(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	if _, ok := c.flags[name]; ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}
	c.keys = append(c.keys, name)
	c.flags[name] = &Flag{name: name, def: "", comment: comment, t: TYPE_STRING, hasValue: false, isShort: false, isArg: true}
}

func (c *CommandLine) FlagLong(name string, def any, t Type, comment string) {
	if err := c.checkLong(name); err != nil {
		debug.Erro(err.Error())
		return
	}
	if _, ok := c.flags[name]; ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}
	c.keys = append(c.keys, name)
	c.flags[name] = &Flag{name: name, def: def, comment: comment, t: t, hasValue: true, isShort: false}
}

func (c *CommandLine) Flag(name string, def any, t Type, comment string) {
	if err := c.checkShort(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	if _, ok := c.flags[name]; ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}
	c.keys = append(c.keys, name)
	c.flags[name] = &Flag{name: name, def: def, comment: comment, t: t, hasValue: true, isShort: true}
}

func (c *CommandLine) FlagNonValueLong(name string, comment string) {
	if err := c.checkLong(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	if _, ok := c.flags[name]; ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}
	c.keys = append(c.keys, name)
	c.flags[name] = &Flag{name: name, hasValue: false, isShort: false, comment: comment}
}

func (c *CommandLine) FlagNonValue(name string, comment string) {
	if err := c.checkShort(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	if _, ok := c.flags[name]; ok {
		debug.Warn("flag[%s] is registed", name)
		return
	}

	c.keys = append(c.keys, name)
	c.flags[name] = &Flag{name: name, hasValue: false, isShort: true, comment: comment}
}

func (c *CommandLine) Arg(index int) *Flag {
	if index >= len(c.others) {
		return nil
	}

	return c.others[index]
}

func (c *CommandLine) Get(name string) *Flag {
	return c.flags[name]
}

func (c *CommandLine) Args() []*Flag {
	return c.others
}

func (c *CommandLine) Parse(args []string) {
	if c.isParsed {
		return
	}

	c.args = args
	c.isParsed = true
	for {
		end, err := c.parseOne()
		if err != nil {
			debug.Erro(err.Error())
			Usage()
			os.Exit(1)
			return
		}

		if end {
			break
		}
	}

	return
}

func (c *CommandLine) checkShort(name string) error {
	reg, err := regexp.Compile("^[a-zA-Z]+$")
	if err != nil {
		return err
	}

	if reg.Match([]byte(name)) {
		return nil
	}

	return fmt.Errorf("expected short name[%s] is [a-zA-Z]", name)
}

func (c *CommandLine) checkLong(name string) error {
	reg, err := regexp.Compile("^[a-zA-Z-]+[a-zA-Z]$")
	if err != nil {
		return err
	}

	if reg.Match([]byte(name)) {
		return nil
	}

	return fmt.Errorf("expected long name[%s] is [a-zA-Z-]", name)
}

func (c *CommandLine) parseShort() (bool, error) {
	name := strings.TrimLeft(c.args[0], "-")
	if err := c.checkShort(name); err != nil {
		return false, err
	}

	if name == "h" || name == "help" {
		Usage()
		os.Exit(0)
	}

	flag, ok := c.flags[name]
	if !ok {
		return false, fmt.Errorf("arg[%s] not defined", name)
	}

	var value = ""
	if len(c.args) >= 2 && c.args[1][0] != '-' {
		if !flag.hasValue {
			return false, fmt.Errorf("arg[%s] has not value", c.args[0])
		}

		value = c.args[1]
		c.args = c.args[2:]
	} else {
		if flag.hasValue {
			return false, fmt.Errorf("arg[%s] has value", c.args[0])
		}

		c.args = c.args[1:]
	}

	c.flags[name].value = value
	c.flags[name].has = true
	c.hasFlag = true
	return len(c.args) == 0, nil
}

func (c *CommandLine) parseLong() (bool, error) {
	arg := strings.TrimLeft(c.args[0], "--")
	if arg == "h" || arg == "help" {
		Usage()
		os.Exit(0)
	}
	if !strings.Contains(arg, "=") {
		if err := c.checkLong(arg); err != nil {
			return false, err
		}

		flag, ok := c.flags[arg]
		if !ok {
			return false, fmt.Errorf("arg[%s] not defined: %+v", arg, c.flags)
		}

		if flag.hasValue {
			return false, fmt.Errorf("arg[%s] has value", arg)
		}

		c.flags[arg].has = true
		c.args = c.args[1:]
		c.hasFlag = true
		return len(c.args) == 0, nil
	}

	info := strings.Split(arg, "=")
	if err := c.checkLong(info[0]); err != nil {
		return false, err
	}
	flag, ok := c.flags[info[0]]
	if !ok {
		return false, fmt.Errorf("arg[%s] not defined", info[0])
	}
	if !flag.hasValue {
		return false, fmt.Errorf("arg[%s] has not value", c.args[0])
	}

	c.flags[info[0]].has = true
	c.flags[info[0]].value = info[1]
	c.args = c.args[1:]
	c.hasFlag = true
	return len(c.args) == 0, nil
}

func (c *CommandLine) parseOne() (bool, error) {
	if len(c.args) == 0 {
		return true, nil
	}

	if c.hasFlag {
		if len(c.args[0]) < 2 || c.args[0][0] != '-' {
			return false, fmt.Errorf("arg[%s] format error", c.args[0])
		}

		if len(c.args[0]) >= 3 && c.args[0][2] == '-' {
			return false, fmt.Errorf("arg[%s] format error", c.args[0])
		}
	}

	if c.args[0][0] == '-' && c.args[0][1] != '-' {
		return c.parseShort()
	}

	if c.args[0][0] == '-' && c.args[0][1] == '-' {
		return c.parseLong()
	}

	f, ok := c.flags[c.args[0]]
	if !ok || f.index != len(c.others) {
		return false, fmt.Errorf("arg[%s] not defined", c.args[0])
	}

	f.value = c.args[0]
	c.others = append(c.others, f)
	c.args = c.args[1:]
	return len(c.args) == 0, nil
}
