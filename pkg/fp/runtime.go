package fp

import "strconv"

// NewPlainRuntime - language specification
func NewPlainRuntime() *Runtime {
	return (&Runtime{
		Stack: []Frame{
			make(Frame),
		},
		parseToken: func(expr string) (interface{}, error) {
			return strconv.Atoi(expr)
		},
		extension: make(map[Name]func(r *Runtime, expr LambdaExpr) Value),
	}).
		WithExtension("let", letExtension).
		WithExtension("lambda", lambdaExtension).
		WithExtension("case", caseExtension)
}

// NewBasicRuntime : minimal set of extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewPlainRuntime().
		WithArithmeticExtension("tail", tailArithmeticExtension).
		WithArithmeticExtension("add", addArithmeticExtension).
		WithArithmeticExtension("sub", subArithmeticExtension).
		WithArithmeticExtension("sign", signArithmeticExtension)
}
