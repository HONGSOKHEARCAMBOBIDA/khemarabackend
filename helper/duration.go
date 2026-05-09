package helper

import "fmt"

func FormatDuration(totalMinutes float64) (hours float64, display string) {
	h := int(totalMinutes) / 60
	m := int(totalMinutes) % 60
	hours = totalMinutes / 60
	display = fmt.Sprintf("%dh %dm", h, m)
	return
}
