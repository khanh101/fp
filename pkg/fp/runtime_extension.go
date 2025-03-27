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

func letExtension(r *Runtime, expr LambdaExpr) LambdaExpr {
	return LambdaExpr{}
}
