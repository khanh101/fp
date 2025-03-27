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

// NewStdRuntime : NewPlainRuntime + standard functions
func NewStdRuntime() *Runtime {
	return NewBasicRuntime().
		LoadExtension("mul", mulExtension).
		LoadExtension("div", divExtension).
		LoadExtension("mod", modExtension).
		LoadExtension("print", printExtension).
		LoadExtension("list", listExtension).
		LoadExtension("append", appendExtension).
		LoadExtension("slice", sliceExtension).
		LoadExtension("peak", peakExtension).
		LoadModule("map", mapModule).
		LoadExtension("type", typeExtension).
		LoadModule("stack", stackModule).
		LoadExtension("unicode", unicodeExtension)
}
