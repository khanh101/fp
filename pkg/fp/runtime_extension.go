package fp

func (r *Runtime) WithArithmeticExtension(name Name, f func(...Value) Value) *Runtime {
	return r.WithExtension(name, func(r *Runtime, expr LambdaExpr) Value {
		var args []Value
		for i := 0; i < len(expr.Args); i++ {
			if i == len(expr.Args)-1 {
				args = append(args, r.Step(expr.Args[i], WithTailCallOptimization))
			} else {
				args = append(args, r.Step(expr.Args[i]))
			}
		}
		return f(args...)
	})
}

func letExtension(r *Runtime, expr LambdaExpr) Value {
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
}

func lambdaExtension(r *Runtime, expr LambdaExpr) Value {
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
}

func caseExtension(r *Runtime, expr LambdaExpr) Value {
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
}
