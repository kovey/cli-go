package gui

import "fmt"

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

func (t *Table) Add(index int, text string) {
	if index > len(t.rows) {
		return
	}
	if index == len(t.rows) {
		t.addRow()
	}

	t.rows[index].Add(text)
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
