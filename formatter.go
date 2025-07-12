package errtrace

import (
	"fmt"
	"path/filepath"
	"strings"
)

const red = "\033[31m"
const yellowBold = "\033[1;33m"
const blueBold = "\033[1;34m"
const greenBold = "\033[1;32m"
const reset = "\033[0m"

func defaultFormatErr(stacktrace []StackFrame, err error, vars []VarPoint) string {
	sb := &strings.Builder{}

	formatHeader(sb, err)
	formatStack(sb, stacktrace)
	formatVars(sb, vars)

	return sb.String()
}

func formatHeader(sb *strings.Builder, err error) {
	sb.WriteString(red)
	fmt.Fprintf(sb, "%v\n", err)
	sb.WriteString(reset)
}

func formatStack(sb *strings.Builder, stacktrace []StackFrame) {
	sb.WriteString(yellowBold)
	sb.WriteString("Stack trace:\n")
	sb.WriteString(reset)

	for _, frame := range stacktrace {
		fmt.Fprintf(sb, "  %s\n", frame.function)

		base := filepath.Base(frame.file)
		rest := frame.file[:len(frame.file)-len(base)]
		fmt.Fprintf(sb, "    %s%s%s:%d\n", rest, red, base, frame.line)
		sb.WriteString(reset)
	}
}

func formatVars(sb *strings.Builder, vars []VarPoint) {
	if len(vars) == 0 {
		return
	}

	sb.WriteString(yellowBold)
	sb.WriteString("Tracked variables:\n")
	sb.WriteString(reset)

	for _, varPoint := range vars {
		frame := varPoint.stacktrace[0]

		sb.WriteString(greenBold)
		fmt.Fprintf(sb, "  %s:%d\n", frame.file, frame.line)
		sb.WriteString(reset)

		for varName, variable := range varPoint.vars {
			sb.WriteString("  - ")
			sb.WriteString(blueBold)
			fmt.Fprintf(sb, "%s%s: %v\n", varName, reset, variable)
		}
	}
}
