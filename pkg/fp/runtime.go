package fp

import (
	"errors"
	"math/big"
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
		extension: make(map[Name]func(r *Runtime, expr LambdaExpr) Object),
	}).
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

// NewIntegerRuntime - BasicRuntime with bigInt
func NewIntegerRuntime() *Runtime {
	return NewPlainRuntime().
		WithParseLiteral(func(lit Name) (Object, error) {
			if i, ok := (&big.Int{}).SetString(lit.String(), 10); ok {
				return i, nil
			}
			return nil, errors.New("integer parse error")
		}).
		WithExtension("reset", resetExtension).
		WithArithmeticExtension("tail", tailArithmeticExtension).
		WithArithmeticExtension("add", func(object ...Object) Object {
			s := (&big.Int{}).SetInt64(0)
			for _, v := range object {
				s.Add(s, v.(*big.Int))
			}
			return s
		}).
		WithArithmeticExtension("sub", func(object ...Object) Object {
			if len(object) != 2 {
				panicError("sub arithmetic error")
			}
			s := (&big.Int{}).SetInt64(object[0].(int64))
			s = s.Sub(s, object[1].(*big.Int))
			return s
		})
}
