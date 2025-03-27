package fp

import (
	"strconv"
)

// NewPlainRuntime - runtime + core control flow extensions
func NewPlainRuntime() *Runtime {
	return (&Runtime{
		debug: false,
		parseLiteral: func(lit Name) (Object, error) {
			return strconv.Atoi(lit.String())
		},
		Stack: []Frame{
			make(Frame),
		},
		Module: make(map[Name]Module),
	}).
		LoadModule("let", letExtension).
		LoadModule("lambda", lambdaExtension).
		LoadModule("case", caseExtension)
}

// NewBasicRuntime : NewPlainRuntime + minimal set of arithmetic extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewPlainRuntime().
		LoadModule("reset", resetExtension).
		WithArithmeticExtension("tail", tailArithmeticExtension).
		WithArithmeticExtension("add", addArithmeticExtension).
		WithArithmeticExtension("sub", subArithmeticExtension).
		WithArithmeticExtension("sign", signArithmeticExtension)
}
