package errtrace

import (
	"fmt"
	"path/filepath"
	"strings"
)

const red = "\033[31m"
const green = "\033[32m"
const blue = "\033[34m"
const yellow = "\033[33m"
const bold = "\033[1m"
const reset = "\033[0m"

func defaultFormatErr(stacktrace []StackFrame, err error, vars []VarPoint) string {
	sb := &strings.Builder{}

	// Error
	sb.WriteString(red)
	sb.WriteString(fmt.Sprintf("%v\n", err))

	// Stack trace
	sb.WriteString(yellow)
	sb.WriteString(bold)
	sb.WriteString("Stack trace:\n")
	sb.WriteString(reset)

	for _, frame := range stacktrace {
		sb.WriteString("  ")
		sb.WriteString(frame.function)
		sb.WriteString("\n")
		base := filepath.Base(frame.file)
		rest := frame.file[:len(frame.file)-len(base)]
		sb.WriteString(fmt.Sprintf("    %s%s%s:%d\n", rest, red, base, frame.line))
		sb.WriteString(reset)
	}

	// Variables
	if len(vars) != 0 {
		sb.WriteString(yellow)
		sb.WriteString(bold)
		sb.WriteString("Tracked variables:\n")
		sb.WriteString(reset)

		for _, varPoint := range vars {
			frame := varPoint.stacktrace[0]
			sb.WriteString(green)
			sb.WriteString(bold)
			sb.WriteString(fmt.Sprintf("  %s:%d\n", frame.file, frame.line))
			sb.WriteString(reset)

			for varName, variable := range varPoint.vars {
				sb.WriteString("  - ")
				sb.WriteString(blue)
				sb.WriteString(bold)
				sb.WriteString(fmt.Sprintf("%s%s: %v\n", varName, reset, variable))
			}
		}
	}

	return sb.String()
}
