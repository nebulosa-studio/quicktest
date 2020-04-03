package status

import (
	"fmt"

	"github.com/fatih/color"
)

// const of status
const (
	Success = iota
	Warning
	Error
)

// Print to stdout
func Print(target string, status int, message string) {
	var st string
	successString := color.New(color.BgGreen, color.FgBlack).Sprint("[success]")
	errorString := color.New(color.BgRed, color.FgBlack).Sprint("[error]")

	if status == Success {
		target = color.New(color.FgGreen).Sprint(target)
		st = successString
	}

	if status == Error {
		target = color.New(color.FgRed).Sprint(target)
		st = errorString
	}

	fmt.Println(
		target,
		st,
		message,
	)
}
