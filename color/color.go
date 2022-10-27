package color

import "fmt"

func format(c Color, content string) string {
	return fmt.Sprintf("%s%s%s", color(c), content, color(Color_Reset))
}

func Green(content string) string {
	return format(Color_Green, content)
}

func Yellow(content string) string {
	return format(Color_Yellow, content)
}

func Red(content string) string {
	return format(Color_Red, content)
}

func White(content string) string {
	return format(Color_White, content)
}

func Blue(content string) string {
	return format(Color_Blue, content)
}

func Magenta(content string) string {
	return format(Color_Magenta, content)
}

func Cyan(content string) string {
	return format(Color_Cyan, content)
}
