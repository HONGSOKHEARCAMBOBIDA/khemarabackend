package helper

import "fmt"

func FormatDuration(totalMinutes float64) (hours float64, display string) {
	h := int(totalMinutes) / 60
	m := int(totalMinutes) % 60
	hours = totalMinutes / 60
	display = fmt.Sprintf("%dម៉ោង %dនាទី", h, m)
	// fmt.Sprintf() → creates a formatted string
	// %d → placeholder for an integer number
	return
}
