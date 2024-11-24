package color

import (
	"fmt"
	"os"
	"strings"
	"regexp"
)

type Color struct {
	FG int
	BG int
}

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Colors = []*Color{
		{FG: 31, BG: 41}, // Red
		{FG: 32, BG: 42}, // Green
		{FG: 33, BG: 43}, // Yellow
		{FG: 34, BG: 44}, // Blue
		{FG: 35, BG: 45}, // Magenta
		{FG: 36, BG: 46}, // Cyan
		{FG: 37, BG: 47}, // White
	}
)

func init() {
	if os.Getenv("NO_COLOR") != "" {
		Reset = ""
		Bold = ""
		for i := range Colors {
			Colors[i] = &Color{FG: 0, BG: 0}
		}
	}
}

func (c *Color) Foreground(s string) string {
	if c.FG == 0 {
		return s
	}
	return fmt.Sprintf("\033[%dm%s%s", c.FG, s, Reset)
}

func (c *Color) Background(s string) string {
	if c.BG == 0 {
		return s
	}
	return fmt.Sprintf("\033[%dm%s%s", c.BG, s, Reset)
}

func Rainbow(s string) string {
	var result strings.Builder
	for i, r := range s {
		color := Colors[i%len(Colors)]
		result.WriteString(color.Foreground(string(r)))
	}
	return result.String()
}

func ColoredLabel(label, value string) string {
	colorIndex := 4 // Default to blue
	color := Colors[colorIndex]
	return fmt.Sprintf("%s%s:%s %s", Bold, color.Foreground(label), Reset, value)
}

func DisplayColorBlocks() string {
	var result strings.Builder
	for _, color := range Colors {
		result.WriteString(color.Background("  ") + Reset + " ")
	}
	return strings.TrimRight(result.String(), " ")
}

func StripANSI(str string) string {
    ansi := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)
    return ansi.ReplaceAllString(str, "")
}
