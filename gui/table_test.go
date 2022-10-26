package gui

import (
	"runtime"
	"testing"
)

func TestTable(t *testing.T) {
	table := NewTable()
	table.Add("kovey")
	table.Add("golang version: " + runtime.Version())
	table.Add("test")
	table.Add("hello world")
	table.Add("yes i do")
	table.Show()
}
