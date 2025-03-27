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
		Module: make(map[Name]Module),
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
		WithArithmeticExtension("tail", tailArithmeticExtension).
		WithArithmeticExtension("add", addArithmeticExtension).
		WithArithmeticExtension("sub", subArithmeticExtension).
		WithArithmeticExtension("sign", signArithmeticExtension)
}

func NewStdRuntime() *Runtime {
	return NewBasicRuntime().
		WithArithmeticExtension("print", printArithmeticExtension).
		WithArithmeticExtension("list", listArithmeticExtension).
		WithArithmeticExtension("append", appendArithmeticExtension).
		WithArithmeticExtension("slice", sliceArithmeticExtension).
		WithArithmeticExtension("peak", peakArithmeticExtension)
}

// NewDebugRuntime : NewBasicRuntime + debug extensions
func NewDebugRuntime() *Runtime {
	return NewStdRuntime().
		LoadModule("stack", stackModule).
		LoadModule("module", moduleModule)
}
