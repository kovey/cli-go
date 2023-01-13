package gui

type Table struct {
	header *Row
	rows   []*Row
	foot   *Row
}

func NewTable() *Table {
	return &Table{rows: make([]*Row, 0), header: NewHeader(), foot: NewFoot()}
}

func (t *Table) Add(data string) {
	t.rows = append(t.rows, NewRow(data))
	t.AddRepeat(Border_Horizontal)
}

func (t *Table) AddRepeat(repeat string) {
	t.rows = append(t.rows, NewRepeatRow(repeat))
}

func (t *Table) rowMaxLen() int {
	var maxLen = 0
	for _, row := range t.rows {
		if row.Len > maxLen {
			maxLen = row.Len
		}
	}

	return maxLen + 2
}

func (t *Table) Show() {
	maxLen := t.rowMaxLen()
	t.header.Adjust(maxLen)
	t.header.Show()
	length := len(t.rows) - 1
	for i := 0; i < length; i++ {
		t.rows[i].Adjust(maxLen)
		t.rows[i].Show()
	}
	t.foot.Adjust(maxLen)
	t.foot.Show()
}
