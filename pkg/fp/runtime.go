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
	return NewPlainRuntime().WithArithmeticExtension("sign", func(value ...Value) Value {
		v := value[len(value)-1].(int)
		switch {
		case v > 0:
			return +1
		case v < 0:
			return -1
		case v == 0:
			return 0
		}
		panic("runtime error")
	}).WithArithmeticExtension("tail", func(value ...Value) Value {
		return value[len(value)-1]
	}).WithArithmeticExtension("sub", func(value ...Value) Value {
		if len(value) != 2 {
			panic("runtime error")
		}
		return value[0].(int) - value[1].(int)
	}).WithArithmeticExtension("add", func(value ...Value) Value {
		v := 0
		for i := 0; i < len(value); i++ {
			v += value[i].(int)
		}
		return v
	})
}
