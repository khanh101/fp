package fp

type ArithmeticExtension = func(...Object) Object

func (r *Runtime) stepWithTailCallOptimization(exprList ...Expr) []Object {
	var outputs []Object
	for i := 0; i < len(exprList)-1; i++ {
		outputs = append(outputs, r.Step(exprList[i]))
	}
	if len(exprList) > 0 {
		outputs = append(outputs, r.Step(exprList[len(exprList)-1], WithTailCallOptimization))
	}
	return outputs
}

func (r *Runtime) WithArithmeticExtension(name Name, f ArithmeticExtension) *Runtime {
	return r.WithExtension(name, func(r *Runtime, expr LambdaExpr) Object {
		return f(r.stepWithTailCallOptimization(expr.Args...)...)
	})
}

func letExtension(r *Runtime, expr LambdaExpr) Object {
	name := expr.Args[0].(Name)
	outputs := r.stepWithTailCallOptimization(expr.Args[1:]...)
	r.Stack[len(r.Stack)-1][name] = outputs[len(outputs)-1]
	return 0
}

func lambdaExtension(r *Runtime, expr LambdaExpr) Object {
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

func caseExtension(r *Runtime, expr LambdaExpr) Object {
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

func tailArithmeticExtension(value ...Object) Object {
	return value[len(value)-1]
}

func addArithmeticExtension(value ...Object) Object {
	v := 0
	for i := 0; i < len(value); i++ {
		v += value[i].(int)
	}
	return v
}

func subArithmeticExtension(value ...Object) Object {
	if len(value) != 2 {
		panic("runtime error")
	}
	return value[0].(int) - value[1].(int)
}

func signArithmeticExtension(value ...Object) Object {
	v := value[len(value)-1].(int)
	switch {
	case v > 0:
		return +1
	case v < 0:
		return -1
	default:
		return 0
	}
}
