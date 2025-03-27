package fp

import "strconv"

// NewPlainRuntime - runtime + core control flow extensions
func NewPlainRuntime() *Runtime {
	return (&Runtime{
		Stack: []Frame{
			make(Frame),
		},
		parseToken: func(expr string) (interface{}, error) {
			return strconv.Atoi(expr)
		},
		extension: make(map[Name]func(r *Runtime, expr LambdaExpr) Object),
	}).
		WithExtension("let", letExtension).
		WithExtension("lambda", lambdaExtension).
		WithExtension("case", caseExtension)
}

// NewBasicRuntime : NewPlainRuntime + minimal set of arithmetic extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewPlainRuntime().
		WithArithmeticExtension("tail", tailArithmeticExtension).
		WithArithmeticExtension("add", addArithmeticExtension).
		WithArithmeticExtension("sub", subArithmeticExtension).
		WithArithmeticExtension("sign", signArithmeticExtension)
}
