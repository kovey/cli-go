package gui

import (
	"fmt"
)

type Row struct {
	Begin    string
	Data     string
	End      string
	Len      int
	prev     string
	suff     string
	isRepeat bool
}

func NewRow(data string) *Row {
	r := &Row{Begin: "+", End: "+", Data: data, isRepeat: false}
	r.init()
	return r
}

func NewRepeatRow(data string) *Row {
	r := NewRow(data)
	r.isRepeat = true
	return r
}

func (r *Row) init() {
	r.Len = len(r.Begin) + len(r.End) + len(r.Data)
}

func (r *Row) Adjust(maxLen int) {
	if maxLen <= r.Len {
		return
	}

	sub := maxLen - r.Len
	if r.isRepeat {
		r.prev += r.Data
	} else {
		r.prev += " "
	}

	for i := 1; i < sub; i++ {
		if r.isRepeat {
			r.suff += r.Data
			continue
		}

		r.suff += " "
	}
}

func (r *Row) Show() {
	fmt.Println(fmt.Sprintf("%s%s%s%s%s", r.Begin, r.prev, r.Data, r.suff, r.End))
}
