package gui

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/kovey/debug-go/color"
)

func TestTable(t *testing.T) {
	table := NewTable()
	table.Add(0, "ID")
	table.Add(0, "名称")
	table.Add(0, "版本")
	table.Add(0, "日期")
	for i := 1; i <= 2; i++ {
		table.Add(i, fmt.Sprintf("%d", i+100))
		table.AddColor(i, fmt.Sprintf("kovey_%d", i), color.Color_Green)
		table.Add(i, runtime.Version())
		table.Add(i, time.Now().Format(time.DateTime))
	}
	table.Show()
}
