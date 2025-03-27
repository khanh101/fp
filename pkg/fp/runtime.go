package fp

import "strconv"

func defaultRuntimeOption() *runtimeOption {
	return &runtimeOption{
		debug: false,
		parseName: func(name Name) (interface{}, error) {
			return strconv.Atoi(string(name))
		},
	}
}

// NewPlainRuntime - runtime + core control flow extensions
func NewPlainRuntime() *Runtime {
	return (&Runtime{
		option: defaultRuntimeOption(),
		Stack: []Frame{
			make(Frame),
		},
		extension: make(map[Name]func(r *Runtime, expr LambdaExpr) Object),
	}).
		WithOption(ParseNameOption(func(name Name) (interface{}, error) {
			return strconv.Atoi(string(name))
		})).
		WithExtension("let", letExtension).
		WithExtension("lambda", lambdaExtension).
		WithExtension("case", caseExtension)
}

// NewBasicRuntime : NewPlainRuntime + minimal set of arithmetic extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewPlainRuntime().
		WithExtension("reset", resetExtension).
		WithArithmeticExtension("tail", tailArithmeticExtension).
		WithArithmeticExtension("add", addArithmeticExtension).
		WithArithmeticExtension("sub", subArithmeticExtension).
		WithArithmeticExtension("sign", signArithmeticExtension)
}
