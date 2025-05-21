package app

import (
	"fmt"
	"strings"
)

/**
	fmt.Printf(`%s command lines.

Usage:
	%s <command> [arguments]
The commands are:
	create  craete config .env file
	start   start app
	reload  reload app
	stop    stop app
	kill	kill app use -9 signal
Use "ksql help <command>" for more information about a command.
**/

type Arg struct {
	Name       string
	Comment    string
	IsShort    bool
	IsRequired bool
}

func (a *Arg) HelpName() string {
	if a.IsShort {
		return fmt.Sprintf("-%s", a.Name)
	}

	return fmt.Sprintf("--%s", a.Name)
}

func (a *Arg) NameLen() int {
	if a.IsShort {
		return len(a.Name) + 1
	}

	return len(a.Name) + 2
}

func (a *Arg) Format(firstSpace, middleSpace string) string {
	if a.IsShort {
		return fmt.Sprintf("%s-%s%s%s", firstSpace, a.Name, middleSpace, a.Comment)
	}

	return fmt.Sprintf("%s--%s%s%s", firstSpace, a.Name, middleSpace, a.Comment)
}

type Args struct {
	argNames []string
	args     map[string]*Arg
}

func NewArgs() *Args {
	return &Args{args: make(map[string]*Arg)}
}

func (a *Args) MaxLen() int {
	maxLen := 0
	for _, arg := range a.args {
		if arg.IsShort {
			if len(arg.Name)+1 > maxLen {
				maxLen = len(arg.Name) + 1
			}

			continue
		}

		if len(arg.Name)+2 > maxLen {
			maxLen = len(arg.Name) + 2
		}
	}

	return maxLen
}

func (a *Args) HelpTitle() string {
	var res []string
	for _, name := range a.argNames {
		if a.args[name].IsRequired {
			res = append(res, fmt.Sprintf("[%s]", a.args[name].HelpName()))
		}
	}
	return strings.Join(res, " ")
}

func (a *Args) HasOptions() bool {
	for _, name := range a.argNames {
		if !a.args[name].IsRequired {
			return true
		}
	}

	return false
}

func (a *Args) Add(name, comment string, isShort, isRequired bool) {
	if _, ok := a.args[name]; !ok {
		a.argNames = append(a.argNames, name)
		a.args[name] = &Arg{}
	}

	a.args[name].Comment = comment
	a.args[name].Name = name
	a.args[name].IsShort = isShort
	a.args[name].IsRequired = isRequired
}

func (a *Args) FormatSub(space string, isRequired bool) string {
	maxLen := a.MaxLen()
	var res []string
	for _, name := range a.argNames {
		if a.args[name].IsRequired != isRequired {
			continue
		}
		res = append(res, a.args[name].Format(space, getSpace(maxLen, a.args[name].NameLen())))
	}

	return strings.Join(res, "\r\n")
}

func (a *Args) Format(space string) string {
	maxLen := a.MaxLen()
	var res []string
	for _, name := range a.argNames {
		res = append(res, a.args[name].Format(space, getSpace(maxLen, a.args[name].NameLen())))
	}

	return strings.Join(res, "\r\n")
}

type Commands struct {
	commandNames []string
	commands     map[string]*Command
}

func NewCommands() *Commands {
	return &Commands{commands: make(map[string]*Command)}
}

func (c *Commands) MaxLen() int {
	maxLen := 0
	for _, command := range c.commands {
		if len(command.Name) > maxLen {
			maxLen = len(command.Name)
		}
	}

	return maxLen
}

func (c *Commands) AddCommand(name, comment string) *Command {
	if _, ok := c.commands[name]; !ok {
		c.commands[name] = &Command{commands: NewCommands(), args: NewArgs()}
		c.commandNames = append(c.commandNames, name)
	}

	c.commands[name].Name = name
	c.commands[name].Comment = comment
	return c.commands[name]
}

func (c *Commands) AddArg(command, name, comment string, isShort, isRequired bool) {
	if com, ok := c.commands[command]; ok {
		com.AddArg(name, command, isShort, isRequired)
	}
}

func (c *Commands) AddSubCommand(parent, name, comment string) *Command {
	if com, ok := c.commands[parent]; ok {
		return com.AddCommand(name, comment)
	}

	return nil
}

func (c *Commands) Format(maxLen int) string {
	var res []string
	for _, name := range c.commandNames {
		res = append(res, c.commands[name].Format("    ", getSpace(maxLen, len(name))))
	}

	return strings.Join(res, "\r\n")
}

type Command struct {
	Name     string
	Comment  string
	commands *Commands
	args     *Args
}

func (c *Command) HelpSub(appName, command string) {
	if sub, ok := c.commands.commands[command]; ok {
		sub.Help(fmt.Sprintf("%s %s", appName, c.Name), appName)
	}
}

func (c *Command) Help(commandPrefix, appName string) {
	if !c.args.HasOptions() {
		fmt.Printf(`"%s" command of %s help details information.

Usage:
    %s %s %s
%s
		`, c.Name, appName, commandPrefix, c.Name, c.args.HelpTitle(), c.args.Format("        "))
		return
	}

	fmt.Printf(`"%s" command of %s help details information.

Usage:
    %s %s %s [options]
%s
options
%s
		`, c.Name, appName, commandPrefix, c.Name, c.args.HelpTitle(), c.args.FormatSub("        ", true), c.args.FormatSub("    ", false))
}

func (c *Command) AddArg(name, comment string, isShort, isRequired bool) {
	c.args.Add(name, comment, isShort, isRequired)
}

func (c *Command) AddCommand(name, comment string) *Command {
	return c.commands.AddCommand(name, comment)
}

func (c *Command) Format(firstSpace, middleSpace string) string {
	return fmt.Sprintf("%s%s%s%s", firstSpace, c.Name, middleSpace, c.Comment)
}

func (c *Command) FormatArg(firstSpace string) string {
	return c.args.Format(firstSpace)
}

type Help struct {
	Title    string
	AppName  string
	Commands *Commands
	Args     *Args
}

func NewHelp(appName string) *Help {
	return &Help{AppName: appName, Commands: NewCommands(), Args: NewArgs()}
}

func getSpace(maxLen, curLen int) string {
	res := "  "
	for i := 0; i < maxLen-curLen; i++ {
		res += " "
	}

	return res
}

func (h *Help) Show() {
	if len(h.Commands.commands) == 0 {
		fmt.Printf(`%s.

Usage:
    %s [options]
The options are:
%s
`, h.Title, h.AppName, h.Args.Format("    "))
		return
	}

	maxLen := h.Commands.MaxLen()
	if len(h.Args.argNames) == 0 {
		fmt.Printf(`%s.

Usage:
    %s <command> [arguments]
The commands are:
%s
Use "%s help <command>" for more information about a command.
`, h.Title, h.AppName, h.Commands.Format(maxLen), h.AppName)
		return
	}

	fmt.Printf(`%s.

Usage:
    %s <command> [arguments] [options]
The commands are:
%s
The options are:
%s
Use "%s help <command>" for more information about a command.
`, h.Title, h.AppName, h.Commands.Format(maxLen), h.Args.Format("    "), h.AppName)
}

func (h *Help) Help(commands ...string) {
	var last *Command
	var prefix []string
	commMap := h.Commands
	for _, command := range commands {
		sub, ok := commMap.commands[command]
		if !ok {
			break
		}

		last = sub
		prefix = append(prefix, command)
		commMap = sub.commands
	}

	if last == nil {
		return
	}

	cs := prefix[:len(prefix)-1]
	if len(cs) > 0 {
		last.Help(fmt.Sprintf("%s %s", h.AppName, strings.Join(prefix[:len(prefix)-1], " ")), h.AppName)
		return
	}

	last.Help(h.AppName, h.AppName)
}
