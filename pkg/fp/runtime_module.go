package fp

import (
	"fmt"
)

var letModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
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
	},
	Man: "module: (let x 3) - assign value 3 to local variable x",
}

var delModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
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
	},
	Man: "module: (del x) - delete variable x",
}

var lambdaModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
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
	},
	Man: "module: (lambda x y (add x y) - declare a function",
}

var caseModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
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
	},
	Man: "module: (case x 1 2 4 5) - case, if x=1 then return 3, if x=4 the return 5",
}

var resetModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
		r.Stack = []Frame{
			make(Frame),
		}
		return nil, nil
	},
	Man: "module: (reset) - reset stack",
}

var tailExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		return values[len(values)-1], nil
	},
	Man: "extension: (tail (print 1) (print 2) 3) - exec a sequence of expressions and return the last one",
}

var addExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		sum := 0
		for i := 0; i < len(values); i++ {
			v, ok := values[i].(int)
			if !ok {
				return nil, fmt.Errorf("adding non-integer values")
			}
			sum += v
		}
		return sum, nil
	},
	Man: "extension: (add 1 (add 2 3) 3) - exec a sequence of expressions and return the sum",
}

var subExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		if len(values) != 2 {
			return nil, fmt.Errorf("subtract requires 2 arguments")
		}
		a, ok := values[0].(int)
		if !ok {
			return nil, fmt.Errorf("subtract non-integer value")
		}
		b, ok := values[1].(int)
		if !ok {
			return nil, fmt.Errorf("subtract non-integer value")
		}
		return a - b, nil
	},
	Man: "extension: (sub 2 (add 1 1)) - exec two expressions and return difference",
}

var signExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		v, ok := values[len(values)-1].(int)
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
	},
	Man: "extension: (sign 3) - exec an expression and return the sign",
}

var listExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		var l List
		for _, v := range values {
			l = append(l, v)
		}
		return l, nil
	},
	Man: "extension: (list 1 2 (lambda x (add x 1))) - make a list",
}

var appendExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		l, ok := values[0].(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		return append(l, values[1:]...), nil
	},
	Man: "extension: (append l 2 (add 1 1)) - append elements into list l and return a new list",
}

var sliceExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		if len(values) != 3 {
			return nil, fmt.Errorf("slice requires 3 arguments")
		}
		l, ok := values[0].(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		if len(l) < 1 {
			return nil, fmt.Errorf("empty list")
		}
		i, ok := values[1].(int)
		if !ok {
			return nil, fmt.Errorf("second argument must be integer")
		}
		j, ok := values[2].(int)
		if !ok {
			return nil, fmt.Errorf("third argument must be integer")
		}
		return l[i:j], nil
	},
	Man: "extension: (slice l 2 3) - make a slice of a list l[2, 3]",
}

var peakExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		if len(values) >= 2 {
			return nil, fmt.Errorf("peak requires at least 2 arguments")
		}
		l, ok := values[0].(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		if len(l) < 1 {
			return nil, fmt.Errorf("empty list")
		}
		var outputs List
		for j := 1; j < len(values); j++ {
			i, ok := values[j].(int)
			if !ok {
				return nil, fmt.Errorf("second argument must be integer")
			}
			outputs = append(outputs, l[i])
		}
		if len(outputs) == 1 {
			return outputs[0], nil
		}
		return outputs, nil
	},
	Man: "extension: (peak l 3 2) - get elem from list (can get multiple elements)",
}

var mapModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
		if len(expr.Args) != 2 {
			return nil, fmt.Errorf("map requires 2 arguments")
		}
		l1, err := r.Step(expr.Args[0])
		if err != nil {
			return nil, err
		}
		l, ok := l1.(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		f1, err := r.Step(expr.Args[1])
		if err != nil {
			return nil, err
		}
		var outputs List
		switch f := f1.(type) {
		case Lambda:
			if len(f.Params) != 1 {
				return nil, fmt.Errorf("map function requires 1 argument")
			}
			for _, v := range l {
				// 2. add argument to local Frame
				localFrame := make(Frame).Update(f.Frame)
				localFrame[f.Params[0]] = v
				// 3. push Frame to Stack
				r.Stack = append(r.Stack, localFrame)
				// 4. exec function
				o, err := r.Step(f.Impl)
				// 5. pop Frame from Stack
				r.Stack = r.Stack[:len(r.Stack)-1]
				// 6. append o
				if err != nil {
					return nil, err
				}
				outputs = append(outputs, o)
			}
		case Module:
			for _, v := range l {
				// 2. add argument to local Frame
				localFrame := make(Frame)
				localFrame["x"] = v // dummy variable
				// 3. make dummy expr and exec
				o, err := f.Exec(r, LambdaExpr{
					Name: "",
					Args: []Expr{Name("x")}, // dummy variable
				})
				// 5. pop Frame from Stack
				r.Stack = r.Stack[:len(r.Stack)-1]
				// 6. append o
				if err != nil {
					return nil, err
				}
				outputs = append(outputs, o)
			}
		default:
			return nil, fmt.Errorf("runtime error: map module requires a function")
		}
		return outputs, nil
	},
	Man: "extension: (map l (lambda y (add 1 y))) - map",
}

// TODO - implement map filter reduce

var typeExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		var types List
		for _, v := range values {
			types = append(types, getType(v))
		}
		if len(types) == 1 {
			return types[0], nil
		}
		return types, nil
	},
	Man: "extension: (type x 1 (lambda y (add 1 y))) - get types of objects (can get multiple ones)",
}

var stackModule = Module{
	Exec: func(r *Runtime, expr LambdaExpr) (Object, error) {
		return r.Stack, nil
	},
	Man: "module: (stack) - get stack",
}

var printExtension = Extension{
	Exec: func(values ...Object) (Object, error) {
		for _, v := range values {
			fmt.Printf("%v ", v)
		}
		fmt.Println()
		return len(values), nil
	},
	Man: "extension: (print 1 x (lambda 3)) - print values",
}
