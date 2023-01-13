package gui

import (
	"fmt"
)

type Row struct {
	Data     string
	Len      int
	prev     string
	suff     string
	isRepeat bool
	IsHeader bool
	IsFoot   bool
}

func NewRow(data string) *Row {
	r := &Row{Data: data, isRepeat: false}
	r.init()
	return r
}

func NewRepeatRow(data string) *Row {
	r := NewRow(data)
	r.isRepeat = true
	return r
}

func NewHeader() *Row {
	r := NewRepeatRow(Border_Horizontal)
	r.IsHeader = true
	return r
}

func NewFoot() *Row {
	r := NewRepeatRow(Border_Horizontal)
	r.IsFoot = true
	return r
}

func (r *Row) init() {
	texts := []rune(r.Data)
	other := 0
	for _, text := range texts {
		if IsChinese(string(text)) {
			other++
		}
	}
	r.Len = len(Border_Left_Up) + len(Border_Right_Up) + len(texts) + other
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
	if r.IsHeader {
		fmt.Println(fmt.Sprintf("%s%s%s%s%s", Border_Left_Up, r.prev, r.Data, r.suff, Border_Right_Up))
		return
	}

	if r.IsFoot {
		fmt.Println(fmt.Sprintf("%s%s%s%s%s", Border_Left_Bottom, r.prev, r.Data, r.suff, Border_Right_Bottom))
		return
	}

	fmt.Println(fmt.Sprintf("%s%s%s%s%s", Border_Vertical, r.prev, r.Data, r.suff, Border_Vertical))
}
