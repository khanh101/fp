package fp

import (
	"strconv"
)

// NewPlainRuntime - runtime + core control flow extensions
func NewPlainRuntime() *Runtime {
	return (&Runtime{
		parseLiteral: func(lit Name) (Object, error) {
			return strconv.Atoi(lit.String())
		},
		Stack: []Frame{
			make(Frame),
		},
	}).
		LoadModule("let", letModule).
		LoadModule("del", delModule).
		LoadModule("lambda", lambdaModule).
		LoadModule("case", caseModule)
}

// NewBasicRuntime : NewPlainRuntime + minimal set of arithmetic extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewPlainRuntime().
		LoadModule("reset", resetModule).
		LoadExtension("tail", tailExtension).
		LoadExtension("add", addExtension).
		LoadExtension("sub", subExtension).
		LoadExtension("sign", signExtension)
}

func NewStdRuntime() *Runtime {
	return NewBasicRuntime().
		LoadExtension("print", printExtension).
		LoadExtension("list", listExtension).
		LoadExtension("append", appendExtension).
		LoadExtension("slice", sliceExtension).
		LoadExtension("peak", peakExtension).
		LoadModule("map", mapModule).
		LoadExtension("type", typeExtension)
}

// NewDebugRuntime : NewBasicRuntime + debug extensions
func NewDebugRuntime() *Runtime {
	return NewStdRuntime().
		LoadModule("stack", stackModule)
}
