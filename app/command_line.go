package app

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kovey/debug-go/debug"
)

var Usage = func() {
	PrintDefaults()
}

var _commanLine = NewCommandLine()

func PrintDefaults() {
	_commanLine.PrintDefaults()
}

func GetHelp() *Help {
	return _commanLine.help
}

type CommandLine struct {
	flags    *Flags
	others   []*Flag
	args     []string
	isParsed bool
	hasFlag  bool
	help     *Help
}

func NewCommandLine() *CommandLine {
	return &CommandLine{flags: NewFlags(), help: NewHelp("")}
}

func (c *CommandLine) PrintDefaults() {
	c.help.Show()
}

func (c *CommandLine) _flag(name, comment string, def any, t Type, index int, hasValue, isShort, isArg bool, parents ...string) {
	var parent *Flag
	if len(parents) > 0 {
		parent = c.flags.Get(parents...)
		if parent == nil {
			debug.Warn("parent not found: %s", name)
			return
		}

		if parent.children != nil && parent.children.Has(name) {
			debug.Warn("flag[%s] is registed", name)
			return
		}

		parent.AddChild(&Flag{name: name, def: def, comment: comment, t: t, index: index, hasValue: hasValue, isShort: isShort, isArg: isArg})
		pCmd := c.help.Get(parents...)
		if isArg {
			pCmd.AddCommand(name, comment)
		} else {
			pCmd.AddArg(name, comment, isShort, false)
		}
	} else {
		if c.flags.Has(name) {
			debug.Warn("flag[%s] is registed", name)
			return
		}
		c.flags.Add(&Flag{name: name, def: def, comment: comment, t: t, index: index, hasValue: hasValue, isShort: isShort, isArg: isArg})
		if isArg {
			c.help.Commands.AddCommand(name, comment)
		} else {
			c.help.Args.Add(name, comment, isShort, false)
		}
	}
}

func (c *CommandLine) FlagArg(name, comment string, index int, parents ...string) {
	if err := c.checkLong(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	c._flag(name, comment, "", TYPE_STRING, index, false, false, true, parents...)
}

func (c *CommandLine) FlagLong(name string, def any, t Type, comment string, parents ...string) {
	if err := c.checkLong(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	c._flag(name, comment, def, t, 0, true, false, false, parents...)
}

func (c *CommandLine) Flag(name string, def any, t Type, comment string, parents ...string) {
	if err := c.checkShort(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	c._flag(name, comment, def, t, 0, true, true, false, parents...)
}

func (c *CommandLine) FlagNonValueLong(name string, comment string, parents ...string) {
	if err := c.checkLong(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	c._flag(name, comment, "", TYPE_STRING, 0, false, false, false, parents...)
}

func (c *CommandLine) FlagNonValue(name string, comment string, parents ...string) {
	if err := c.checkShort(name); err != nil {
		debug.Erro(err.Error())
		return
	}

	c._flag(name, comment, "", TYPE_STRING, 0, false, true, false, parents...)
}

func (c *CommandLine) Arg(index int) *Flag {
	if index >= len(c.others) {
		return nil
	}

	return c.others[index]
}

func (c *CommandLine) Get(names ...string) *Flag {
	return c.flags.Get(names...)
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

	commands := append(c.AllArgName(), name)
	flag := c.flags.Get(commands...)
	if flag == nil {
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

	flag.value = value
	flag.has = true
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

		commands := append(c.AllArgName(), arg)
		flag := c.flags.Get(commands...)
		if flag == nil {
			return false, fmt.Errorf("arg[%s] not defined: %+v", arg, c.flags)
		}

		if flag.hasValue {
			return false, fmt.Errorf("arg[%s] has value", arg)
		}

		flag.has = true
		c.args = c.args[1:]
		c.hasFlag = true
		return len(c.args) == 0, nil
	}

	info := strings.Split(arg, "=")
	if err := c.checkLong(info[0]); err != nil {
		return false, err
	}
	commands := append(c.AllArgName(), info[0])
	flag := c.flags.Get(commands...)
	if flag == nil {
		return false, fmt.Errorf("arg[%s] not defined", info[0])
	}
	if !flag.hasValue {
		return false, fmt.Errorf("arg[%s] has not value", c.args[0])
	}

	flag.has = true
	flag.value = info[1]
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

	commands := append(c.AllArgName(), c.args[0])
	flag := c.flags.Get(commands...)
	if flag == nil {
		return false, fmt.Errorf("arg[%s] not defined", c.args[0])
	}

	if c.args[0] == Ko_Command_Help {
		if c.args[0] == flag.name {
			flag.value = c.args[0]
			c.others = append(c.others, flag)
		} else {
			c.others = append(c.others, &Flag{name: c.args[0], value: c.args[0], def: c.args[0], t: TYPE_STRING, isArg: true})
		}
	} else {
		flag.value = c.args[0]
		c.others = append(c.others, flag)
	}
	c.args = c.args[1:]
	return len(c.args) == 0, nil
}

func (c *CommandLine) AllArgName() []string {
	var res []string
	for _, f := range c.others {
		res = append(res, f.String())
	}
	return res
}

func (c *CommandLine) hasHelp() bool {
	for _, other := range c.others {
		if other.name == Ko_Command_Help {
			return true
		}
	}

	return false
}

func (c *CommandLine) Help() {
	var argNames []string
	for _, name := range c.AllArgName() {
		if name == Ko_Command_Help {
			continue
		}

		argNames = append(argNames, name)
	}
	if len(argNames) == 0 {
		c.help.Show()
		return
	}

	c.help.Help(argNames...)
}
