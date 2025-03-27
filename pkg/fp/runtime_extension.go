package fp

import (
	"errors"
	"fmt"
)

type ArithmeticExtension = func(...Object) (Object, error)

func (r *Runtime) stepWithTailCallOptimization(exprList ...Expr) ([]Object, error) {
	return r.stepWithTailOption(TCOStepOption(true), exprList...)
}

func (r *Runtime) WithArithmeticExtension(name Name, f ArithmeticExtension) *Runtime {
	return r.LoadModule(name, func(r *Runtime, expr LambdaExpr) (Object, error) {
		args, err := r.stepWithTailCallOptimization(expr.Args...)
		if err != nil {
			return nil, err
		}
		return f(args...)
	})
}

func letExtension(r *Runtime, expr LambdaExpr) (Object, error) {
	name := expr.Args[0].(Name)
	outputs, err := r.stepWithTailCallOptimization(expr.Args[1:]...)
	if err != nil {
		return nil, err
	}
	if len(outputs) == 0 {
		return nil, errors.New("let of nothing")
	}
	r.Stack[len(r.Stack)-1][name] = outputs[len(outputs)-1]
	return 0, nil
}

func lambdaExtension(r *Runtime, expr LambdaExpr) (Object, error) {
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
	return v, nil
}

func caseExtension(r *Runtime, expr LambdaExpr) (Object, error) {
	cond, err := r.Step(expr.Args[0])
	if err != nil {
		return nil, err
	}
	i, err := func() (int, error) {
		for i := 1; i < len(expr.Args); i += 2 {
			if arg, ok := expr.Args[i].(Name); ok && arg == "_" {
				return i, nil
			}
			comp, err := r.Step(expr.Args[i])
			if err != nil {
				return 0, err
			}
			if comp == cond {
				return i, nil
			}
		}
		return 0, fmt.Errorf("runtime error: no case matched %s", expr)
	}()
	if err != nil {
		return nil, err
	}
	return r.Step(expr.Args[i+1], TCOStepOption(true))
}

func resetExtension(r *Runtime, expr LambdaExpr) (Object, error) {
	r.Stack = []Frame{
		make(Frame),
	}
	return nil, nil
}

func tailArithmeticExtension(value ...Object) (Object, error) {
	return value[len(value)-1], nil
}

func addArithmeticExtension(value ...Object) (Object, error) {
	sum := 0
	for i := 0; i < len(value); i++ {
		v, ok := value[i].(int)
		if !ok {
			return nil, fmt.Errorf("adding non-integer value")
		}
		sum += v
	}
	return sum, nil
}

func subArithmeticExtension(value ...Object) (Object, error) {
	if len(value) != 2 {
		return nil, fmt.Errorf("subtract requires 2 arguments")
	}
	a, ok := value[0].(int)
	if !ok {
		return nil, fmt.Errorf("subtract non-integer value")
	}
	b, ok := value[0].(int)
	if !ok {
		return nil, fmt.Errorf("subtract non-integer value")
	}
	return a - b, nil
}

func signArithmeticExtension(value ...Object) (Object, error) {
	v, ok := value[len(value)-1].(int)
	if !ok {
		return nil, fmt.Errorf("sign non-integer value")
	}
	switch {
	case v > 0:
		return +1, nil
	case v < 0:
		return -1, nil
	default:
		return 0, nil
	}
}
