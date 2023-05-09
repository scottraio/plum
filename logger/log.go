package logger

import (
	"fmt"
	"os"
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

	colorCodeText(output, colorCode)

	PersistLog(output)
}

func Print(message string, colorCode string) {
	colorCodeText(message, colorCode)
}

func LogAndFail(label string, message string, colorCode string) {
	Log(label, message, colorCode)
	os.Exit(1)
}

func colorCodeText(output string, colorCode string) {
	switch colorCode {
	case "purple":
		color.Magenta(output)
	case "green":
		color.Green(output)
	case "cyan":
		color.Cyan(output)
	case "hi-cyan":
		color.HiCyan(output)
	case "gray":
		color.White(output)
	case "yellow":
		color.Yellow(output)
	case "red":
		color.Red(output)
	default:
		println(output)
	}
}

func PersistLog(output string) error {
	// Open the log file for appending or create it if it does not exist
	f, err := os.OpenFile("./logs/plum.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer f.Close()

	// Write the log message to the file
	if _, err := f.WriteString(output + "\n"); err != nil {
		return fmt.Errorf("error writing to log file: %v", err)
	}

	return nil
}
