package logger

import (
	"fmt"
	"strings"
	"time"

	color "github.com/fatih/color"
)

// Log events with color
func Log(label string, message string, colorCode string) {
	msg := strings.ReplaceAll(message, "\n", "")

	// Get current date and time
	now := time.Now()

	// Format the date and time as desired
	dateTimeStr := now.Format("2006-01-02 15:04:05")
	output := fmt.Sprintf("[%s] [%s] %s", dateTimeStr, label, msg)

	switch colorCode {
	case "purple":
		color.Magenta(output)
	case "green":
		color.Green(output)
	case "cyan":
		color.Cyan(output)
	case "gray":
		color.White(output)
	case "orange":
		color.Yellow(output)
	case "red":
		color.Red(output)
	default:
		println(output)
	}
}
