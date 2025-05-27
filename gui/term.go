package gui

import (
	"fmt"
	"os"

	"github.com/kovey/debug-go/color"
	"golang.org/x/term"
)

func Middle(left, right string) string {
	w, _, _ := term.GetSize(int(os.Stdin.Fd()))
	middle := " "
	leftLen := len(left)
	rightLen := len(right)
	middleLen := len(middle)
	if leftLen+rightLen+middleLen > w {
		return middle
	}

	for i := 0; i < w-(leftLen+middleLen+rightLen); i++ {
		middle += " "
	}

	return middle
}

func Println(left, middle, right string) {
	fmt.Printf("%s%s%s\n", left, middle, right)
}

func PrintlnNormal(left, right string) {
	Println(left, Middle(left, right), right)
}

func PrintlnOk(format string, args ...any) {
	content := fmt.Sprintf(format, args...)
	middle := Middle(content, "[ok]")
	Println(color.Green(content), middle, fmt.Sprintf("[%s]", color.Green("ok")))
}

func PrintlnFailure(format string, args ...any) {
	content := fmt.Sprintf(format, args...)
	middle := Middle(content, "[failure]")
	Println(color.Red(content), middle, fmt.Sprintf("[%s]", color.Red("failure")))
}
