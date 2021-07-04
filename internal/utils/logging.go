package utils

import (
	"os"

	"github.com/fatih/color"
)

var WarnStyle = color.New(color.FgHiYellow, color.Bold)
var ErrorStyle = color.New(color.FgHiRed, color.Bold)
var InfoStyle = color.New(color.FgHiGreen, color.Bold)
var PrintStyle = color.New(color.FgWhite)

var Warn = WarnStyle.PrintfFunc()
var Error = ErrorStyle.SprintfFunc()

// print errors with a format, interface like in printf
// Example:
// 
// ErrorP("%d",1)
func ErrorP(format string, a ...interface{}) {
	ErrorStyle.Printf(format, a...)
	os.Exit(1)
}

var Info = InfoStyle.PrintfFunc()
var Print = PrintStyle.PrintfFunc()
