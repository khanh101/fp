package fp

import (
	"fmt"
)

type Extension = func(...Object) (Object, error)

func (r *Runtime) LoadExtension(name Name, f Extension) *Runtime {
	return r.LoadModule(name, func(r *Runtime, expr LambdaExpr) (Object, error) {
		args, err := r.stepMany(expr.Args...)
		if err != nil {
			return nil, err
		}
		return f(args...)
	})
}

func letModule(r *Runtime, expr LambdaExpr) (Object, error) {
	if len(expr.Args) < 2 {
		return nil, fmt.Errorf("not enough arguments for let")
	}
	name := expr.Args[0].(Name)
	outputs, err := r.stepMany(expr.Args[1:]...)
	if err != nil {
		return nil, err
	}
	r.Stack[len(r.Stack)-1][name] = outputs[len(outputs)-1]
	return outputs[len(outputs)-1], nil
}

func delModule(r *Runtime, expr LambdaExpr) (Object, error) {
	if len(expr.Args) < 1 {
		return nil, fmt.Errorf("not enough arguments for del")
	}
	name := expr.Args[0].(Name)
	_, err := r.stepMany(expr.Args[1:]...)
	if err != nil {
		return nil, err
	}
	delete(r.Stack[len(r.Stack)-1], name)
	return nil, nil
}

func lambdaModule(r *Runtime, expr LambdaExpr) (Object, error) {
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

func caseModule(r *Runtime, expr LambdaExpr) (Object, error) {
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
	return r.Step(expr.Args[i+1])
}

func resetModule(r *Runtime, expr LambdaExpr) (Object, error) {
	r.Stack = []Frame{
		make(Frame),
	}
	return nil, nil
}

func tailExtension(value ...Object) (Object, error) {
	return value[len(value)-1], nil
}

func addExtension(value ...Object) (Object, error) {
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

func subExtension(value ...Object) (Object, error) {
	if len(value) != 2 {
		return nil, fmt.Errorf("subtract requires 2 arguments")
	}
	a, ok := value[0].(int)
	if !ok {
		return nil, fmt.Errorf("subtract non-integer value")
	}
	b, ok := value[1].(int)
	if !ok {
		return nil, fmt.Errorf("subtract non-integer value")
	}
	return a - b, nil
}

func signExtension(value ...Object) (Object, error) {
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

func listExtension(value ...Object) (Object, error) {
	var l []Object
	for _, v := range value {
		l = append(l, v)
	}
	return l, nil
}

func appendExtension(value ...Object) (Object, error) {
	l, ok := value[0].([]Object)
	if !ok {
		return nil, fmt.Errorf("first argument must be list")
	}
	return append(l, value[1:]...), nil
}

func sliceExtension(value ...Object) (Object, error) {
	if len(value) != 3 {
		return nil, fmt.Errorf("slice requires 3 arguments")
	}
	l, ok := value[0].([]Object)
	if !ok {
		return nil, fmt.Errorf("first argument must be list")
	}
	if len(l) < 1 {
		return nil, fmt.Errorf("empty list")
	}
	i, ok := value[1].(int)
	if !ok {
		return nil, fmt.Errorf("second argument must be integer")
	}
	j, ok := value[2].(int)
	if !ok {
		return nil, fmt.Errorf("third argument must be integer")
	}
	return l[i:j], nil
}

func peakExtension(value ...Object) (Object, error) {
	if len(value) != 2 {
		return nil, fmt.Errorf("peak requires 2 arguments")
	}
	l, ok := value[0].([]Object)
	if !ok {
		return nil, fmt.Errorf("first argument must be list")
	}
	if len(l) < 1 {
		return nil, fmt.Errorf("empty list")
	}
	i, ok := value[1].(int)
	if !ok {
		return nil, fmt.Errorf("second argument must be integer")
	}
	return l[i], nil
}

// TODO - implement map filter reduce

func stackModule(r *Runtime, expr LambdaExpr) (Object, error) {
	_, err := r.stepMany(expr.Args...)
	if err != nil {
		return nil, err
	}
	return r.Stack, nil
}

func moduleModule(r *Runtime, expr LambdaExpr) (Object, error) {
	_, err := r.stepMany(expr.Args...)
	if err != nil {
		return nil, err
	}
	return r.Module, nil
}

func printExtension(value ...Object) (Object, error) {
	for _, v := range value {
		fmt.Printf("%v ", v)
	}
	fmt.Println()
	return len(value), nil
}
