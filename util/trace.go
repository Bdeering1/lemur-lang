package util

import (
	"fmt"
	"strings"
)

const traceIdentPlaceholder string = "  "

var traceLevel int = 0

func Trace(msg string) string {
	indent()
	PrintTrace("BEGIN " + msg)
	return msg
}
func Untrace(msg string) {
	PrintTrace("END " + msg)
	dedent()
}

func PrintTrace(fs string) {
	fmt.Printf("%s%s\n", indentStr(), fs)
}

func indentStr() string { return strings.Repeat(traceIdentPlaceholder, traceLevel - 1) }
func indent() { traceLevel = traceLevel + 1 }
func dedent() { traceLevel = traceLevel - 1 }
