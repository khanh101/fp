package fp

import (
	"strconv"
)

// NewCoreRuntime - runtime + core control flow extensions
func NewCoreRuntime() *Runtime {
	return (&Runtime{
		parseLiteral: func(lit String) (Object, error) {
			if lit == "_" {
				return Wildcard{}, nil
			}
			if lit == "*" {
				return Unwrap{}, nil
			}
			i, err := strconv.Atoi(lit.String())
			return Int(i), err
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

// NewBasicRuntime : NewCoreRuntime + minimal set of arithmetic extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewCoreRuntime().
		LoadExtension("tail", tailExtension).
		LoadExtension("add", addExtension).
		LoadExtension("sub", subExtension).
		LoadExtension("sign", signExtension)
}

// NewStdRuntime : NewCoreRuntime + standard functions
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
		LoadExtension("unicode", unicodeExtension).
		LoadModule("kaboom", kaboomModule).
		LoadExtension("doom", doomExtension)
}
