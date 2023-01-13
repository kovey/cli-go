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
	table.Add("yes i do a,这是好事")
	table.Add("这是很好的事情啦！")
	table.Show()
}
