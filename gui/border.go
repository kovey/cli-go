package gui

import (
	"fmt"
	"regexp"
)

const (
	Border_Left_Up      = "┌"
	Border_Left_Bottom  = "└"
	Border_Right_Up     = "┐"
	Border_Right_Bottom = "┘"
	Border_Horizontal   = "─"
	Border_Vertical     = "│"
	Border_T            = "┬"
	Border_Un_T         = "┴"
	Border_Left_Center  = "├"
	Border_Right_Center = "┤"
	Border_Center       = "┼"

	Reg_Chinese      = "^[\u4E00-\u9FA5]+$"
	Reg_Chinese_Sign = "^[\u3002|\uff1f|\uff01|\uff0c|\u3001|\uff1b|\uff1a|\u201c|\u201d|\u2018|\u2019|\uff08|\uff09|\u300a|\u300b|\u3010|\u3011|\u007e]+$"
)

func IsChinese(text string) bool {
	if ok, _ := regexp.Match(Reg_Chinese, []byte(text)); ok {
		return ok
	}

	ok, _ := regexp.Match(Reg_Chinese_Sign, []byte(text))
	return ok
}

type Position byte

const (
	Position_Left_Up    Position = 1
	Position_Left       Position = 2
	Position_Right_Up   Position = 3
	Position_Right      Position = 4
	Position_Center     Position = 5
	Position_Up         Position = 6
	Position_Down       Position = 7
	Position_Left_Down  Position = 8
	Position_Right_Down Position = 9
)

type Border struct {
	Data string
	Len  int
	prev string
	suff string
}

func NewBorder(data string) *Border {
	b := &Border{Data: data}
	b.init()
	return b
}

func (b *Border) init() {
	texts := []rune(b.Data)
	other := 0
	for _, text := range texts {
		if IsChinese(string(text)) {
			other++
		}
	}
	b.Len = len(Border_Left_Up) + len(Border_Right_Up) + len(texts) + other
}
func (b *Border) Adjust(maxLen int) {
	if maxLen <= b.Len {
		return
	}

	sub := maxLen - b.Len
	b.prev += b.Data

	for i := 1; i < sub; i++ {
		b.suff += b.Data
	}
}

func (b *Border) Text() string {
	if b == nil {
		return ""
	}

	return fmt.Sprintf("%s%s%s", b.prev, b.Data, b.suff)
}
