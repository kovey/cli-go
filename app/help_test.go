package app

import (
	"testing"
)

func TestHelp(t *testing.T) {
	h := NewHelp("test")
	command := h.Commands.AddCommand("create", "create module")
	command.AddArg("path", "config path", false, true)
	command.AddArg("type", "config type: json|.env|yaml, default is .env", false, false)
	command.AddArg("c", "short test", true, false)
	h.Args.Add("daemon", "run with daemon mode", false, false)
	h.Args.Add("demo", "run with demo mode", false, false)
	sub := command.AddCommand("config", "create config")
	sub.AddArg("o", "open src", true, false)
	open := h.Commands.AddCommand("open", "open module")
	open.AddArg("path", "config path", false, true)
	h.Show()
	h.Help("create")
	h.Help("create", "config")
	h.Help("open")
	h.Commands.commandNames = nil
	h.Commands.commands = make(map[string]*Command)
	h.Show()
}
