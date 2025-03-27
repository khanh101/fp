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
	}).WithExtension("let", func(r *Runtime, expr LambdaExpr) Value {
		name := expr.Args[0].(Name)
		var v Value
		for i := 1; i < len(expr.Args); i++ {
			if i == len(expr.Args)-1 {
				v = r.Step(expr.Args[i], WithTailCallOptimization)
			} else {
				v = r.Step(expr.Args[i])
			}
		}
		r.Stack[len(r.Stack)-1][name] = v
		return v
	}).WithExtension("lambda", func(r *Runtime, expr LambdaExpr) Value {
		v := Lambda{
			Params: nil,
			Impl:   nil,
			Frame:  nil,
		}
		for i := 0; i < len(expr.Args)-1; i++ {
			paramName := expr.Args[i].(Name)
			v.Params = append(v.Params, paramName)
		}
		v.Impl = expr.Args[len(expr.Args)-1]
		v.Frame = make(Frame).Update(r.Stack[len(r.Stack)-1])
		return v
	}).WithExtension("case", func(r *Runtime, expr LambdaExpr) Value {
		cond := r.Step(expr.Args[0])
		i := func() int {
			for i := 1; i < len(expr.Args); i += 2 {
				if arg, ok := expr.Args[i].(Name); ok && arg == "_" {
					return i
				}
				if r.Step(expr.Args[i]) == cond {
					return i
				}
			}
			panic("runtime error")
		}()
		return r.Step(expr.Args[i+1], WithTailCallOptimization)
	})
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
