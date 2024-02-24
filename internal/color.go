package internal

import "fmt"

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed    Color = "\u001b[31m"
	ColorGreen  Color = "\u001b[32m"
	ColorYellow Color = "\u001b[33m"
	ColorBlue   Color = "\u001b[34m"
	ColorReset  Color = "\u001b[0m"
)

func Colorize(color Color, s string) string {
	return string(color) + s + string(ColorReset)
}

func Colorizef(color Color, format string, a ...any) string {
	return string(color) + fmt.Sprintf(format, a...) + string(ColorReset)
}
