package utils

import (
	"os"

	"github.com/fatih/color"
)

var (
	// WarnStyle is a bold yellow hue
	WarnStyle = color.New(color.FgHiYellow, color.Bold)
	// ErrorStyle is a bold red hue
	ErrorStyle = color.New(color.FgHiRed, color.Bold)
	// InfoStyle is a bold green hue
	InfoStyle = color.New(color.FgHiGreen, color.Bold)
	// PrintStyle is a white hue
	PrintStyle = color.New(color.FgWhite)
)

var (
	// Warn is a printf function with WarnStyle
	Warn = WarnStyle.PrintfFunc()
	// Error is a sprintf function with ErrorStyle
	Error = ErrorStyle.SprintfFunc()
	// Info is a printf function with InfoStyle
	Info = InfoStyle.PrintfFunc()
	// Print is a printf function with PrintStyle
	Print = PrintStyle.PrintfFunc()
)

// ErrorP prints errors with a format and interface like in printf and exits the program.
func ErrorP(format string, a ...interface{}) {
	ErrorStyle.Printf(format, a...)
	os.Exit(1)
}
