package log

import (
	"github.com/fatih/color"
)

var WarnStyle = color.New(color.FgHiYellow, color.Bold)
var ErrorStyle = color.New(color.FgHiRed, color.Bold)
var InfoStyle = color.New(color.FgHiYellow, color.Bold)

var Warn = WarnStyle.PrintlnFunc()
var Error = ErrorStyle.SprintfFunc()
var Info = InfoStyle.PrintlnFunc()
