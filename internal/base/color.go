package base

import (
	"fmt"
	"strconv"

	"github.com/gookit/color"
)

func Bold(text string) string {
	return color.Bold.Sprint(text)
}

func ColorizeLine(order string, text string) string {
	colorNo := colorNo(order)
	if colorNo == 0 {
		return text
	}
	s := color.S256(colorNo)
	return s.Sprint(text)
}

func Colorize(order string, text string, length int) string {
	colorNo := colorNo(order)
	if colorNo == 0 {
		return fmt.Sprintf("%"+strconv.Itoa(length)+"s ", text)
	}
	s := color.S256(colorNo)
	return s.Sprintf("%"+strconv.Itoa(length)+"s ", text)
}

func colorNo(order string) uint8 {
	switch order {
	case "HBT gestalten":
		return 93 // lila
	case "Krankheit":
		return 226 // gelb
	case "Fortbildung, Coach.":
		return 21 // blau
	case "Elternzeit":
		return 51 // blaugrün
	case "Urlaub":
		return 46 // light green
	default:
		return 0 // means nothing
	}
}
