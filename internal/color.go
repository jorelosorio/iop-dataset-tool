package internal

import (
	"fmt"
	"math"
	"strings"
)

type Color string

var (
	ColorMagenta Color = "\x1b[35;1m"
	ColorBlue    Color = "\x1b[34;1m"
	ColorCyan    Color = "\x1b[36;1m"
	ColorYellow  Color = "\x1b[33;1m"
	Reset        Color = "\x1b[0m"
)

func Colorize(color Color, s string) string {
	return string(color) + s + string(Reset)
}

func Colorizef(color Color, format string, a ...any) string {
	return string(color) + fmt.Sprintf(format, a...) + string(Reset)
}

func PsychedelicGradient(text string) string {
	lines := strings.Split(text, "\n")
	var result strings.Builder

	colors := []Color{
		ColorMagenta,
		ColorYellow,
		ColorCyan,
		ColorBlue,
	}
	numColors := len(colors)
	numLines := len(lines)
	linesPerColor := int(math.Ceil(float64(numLines) / float64(numColors)))

	for i, line := range lines {
		colorIndex := i / linesPerColor
		color := colors[colorIndex]
		coloredLine := Colorize(color, line)
		result.WriteString(coloredLine)
		result.WriteString("\n")
	}

	return result.String()
}
