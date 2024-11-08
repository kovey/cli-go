package gui

import (
	"fmt"

	"github.com/kovey/debug-go/color"
)

type Table struct {
	rows     []*Row
	isAdjust bool
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) addRow() {
	row := NewRow(len(t.rows))
	t.rows = append(t.rows, row)
}

func (t *Table) AddRow(row *Row) {
	row.index = len(t.rows)
	t.rows = append(t.rows, row)
}

func (t *Table) AddColor(index int, text string, color color.Color) {
	if index > len(t.rows) {
		return
	}
	if index == len(t.rows) {
		t.addRow()
	}

	t.rows[index].AddColor(text, color)
}

func (t *Table) Add(index int, text string) {
	t.AddColor(index, text, color.Color_None)
}

func (t *Table) AddInt(index, data int) {
	t.Add(index, fmt.Sprintf("%d", data))
}

func (t *Table) AddAny(index int, data any) {
	t.Add(index, fmt.Sprintf("%v", data))
}

func (t *Table) adjust() {
	if t.isAdjust {
		return
	}
	t.isAdjust = true
	rowCount := len(t.rows)
	if rowCount == 0 {
		return
	}

	for _, r := range t.rows {
		r.Final(rowCount)
	}

	columnCount := len(t.rows[0].columns)
	for i := 0; i < columnCount; i++ {
		maxLen := 0
		for _, r := range t.rows {
			if r.ColumnLen(i) > maxLen {
				maxLen = r.ColumnLen(i)
			}
		}

		for _, r := range t.rows {
			r.Adjust(i, maxLen)
		}
	}
}

func (t *Table) Show() {
	t.adjust()
	for _, r := range t.rows {
		fmt.Printf("%s\n%s\n%s", r.UpBorder(), r.Text(), r.DownBorder())
	}

	fmt.Println("")
}
