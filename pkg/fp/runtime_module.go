package fp

import (
	"fmt"
)

var letModule = Module{
	Name: "let",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		if len(expr.Args) < 2 {
			return nil, fmt.Errorf("not enough arguments for let")
		}
		name := String(expr.Args[0].(Name))
		outputs, err := r.stepMany(interruptCh, expr.Args[1:]...)
		if err != nil {
			return nil, err
		}
		r.Stack[len(r.Stack)-1][name] = outputs[len(outputs)-1]
		return outputs[len(outputs)-1], nil
	},
	Man: "module: (let x 3) - assign value 3 to local variable x",
}

var delModule = Module{
	Name: "del",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		if len(expr.Args) < 1 {
			return nil, fmt.Errorf("not enough arguments for del")
		}
		name := String(expr.Args[0].(Name))
		_, err := r.stepMany(interruptCh, expr.Args[1:]...)
		if err != nil {
			return nil, err
		}
		delete(r.Stack[len(r.Stack)-1], name)
		return nil, nil
	},
	Man: "module: (del x) - delete variable x",
}

var lambdaModule = Module{
	Name: "lambda",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		v := Lambda{
			Params: nil,
			Impl:   nil,
			Frame:  nil,
		}
		for i := 0; i < len(expr.Args)-1; i++ {
			paramName := String(expr.Args[i].(Name))
			v.Params = append(v.Params, paramName)
		}
		v.Impl = expr.Args[len(expr.Args)-1]
		v.Frame = make(Frame).Update(r.Stack[len(r.Stack)-1])
		return v, nil
	},
	Man: "module: (lambda x y (add x y) - declare a function",
}

var caseModule = Module{
	Name: "case",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		cond, err := r.Step(expr.Args[0], interruptCh)
		if err != nil {
			return nil, err
		}
		i, err := func() (int, error) {
			for i := 1; i < len(expr.Args); i += 2 {
				comp, err := r.Step(expr.Args[i], interruptCh)
				if err != nil {
					return 0, err
				}
				if _, ok := comp.(Wildcard); ok || comp == cond {
					return i, nil
				}
			}
			return 0, fmt.Errorf("runtime error: no case matched %s", expr)
		}()
		if err != nil {
			return nil, err
		}
		return r.Step(expr.Args[i+1], interruptCh)
	},
	Man: "module: (case x 1 2 4 5) - case, if x=1 then return 3, if x=4 the return 5",
}

var kaboomModule = Module{
	Name: "kaboom",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		r.Stack = r.Stack[0:1]
		return nil, nil
	},
	Man: "module: (kaboom) - remove everything except global frame",
}

var doomExtension = Extension{
	Name: "doom",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		return String(fmt.Sprintf("i told you - we don't have Doom yet")), nil
	},
	Man: "module: (doom) - extra modules required https://youtu.be/dQw4w9WgXcQ",
}

var tailExtension = Extension{
	Name: "tail",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		return values[len(values)-1], nil
	},
	Man: "module: (tail (print 1) (print 2) 3) - exec a sequence of expressions and return the last one",
}

var addExtension = Extension{
	Name: "add",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		var sum Int = 0
		for i := 0; i < len(values); i++ {
			v, ok := values[i].(Int)
			if !ok {
				return nil, fmt.Errorf("adding non-integer values")
			}
			sum += v
		}
		return sum, nil
	},
	Man: "module: (add 1 (add 2 3) 3) - exec a sequence of expressions and return the sum",
}

var mulExtension = Extension{
	Name: "mul",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		var sum Int = 1
		for i := 0; i < len(values); i++ {
			v, ok := values[i].(Int)
			if !ok {
				return nil, fmt.Errorf("multiplying non-integer values")
			}
			sum *= v
		}
		return sum, nil
	},
	Man: "module: (mul 1 (add 2 3) 3) - exec a sequence of expressions and return the product",
}

var subExtension = Extension{
	Name: "sub",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		if len(values) != 2 {
			return nil, fmt.Errorf("subtract requires 2 arguments")
		}
		a, ok := values[0].(Int)
		if !ok {
			return nil, fmt.Errorf("subtract non-integer value")
		}
		b, ok := values[1].(Int)
		if !ok {
			return nil, fmt.Errorf("subtract non-integer value")
		}
		return a - b, nil
	},
	Man: "module: (sub 2 (add 1 1)) - exec two expressions and return difference",
}

var divExtension = Extension{
	Name: "div",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		if len(values) != 2 {
			return nil, fmt.Errorf("dividing requires 2 arguments")
		}
		a, ok := values[0].(Int)
		if !ok {
			return nil, fmt.Errorf("dividing non-integer value")
		}
		b, ok := values[1].(Int)
		if !ok {
			return nil, fmt.Errorf("dividing non-integer value")
		}
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return a / b, nil
	},
	Man: "module: (div 2 (add 1 1)) - exec two expressions and return ratio",
}

var modExtension = Extension{
	Name: "mod",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		if len(values) != 2 {
			return nil, fmt.Errorf("dividing requires 2 arguments")
		}
		a, ok := values[0].(Int)
		if !ok {
			return nil, fmt.Errorf("dividing non-integer value")
		}
		b, ok := values[1].(Int)
		if !ok {
			return nil, fmt.Errorf("dividing non-integer value")
		}
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return a % b, nil
	},
	Man: "module: (mod 2 (add 1 1)) - exec two expressions and return modulo",
}

var signExtension = Extension{
	Name: "sign",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		v, ok := values[len(values)-1].(Int)
		if !ok {
			return nil, fmt.Errorf("sign non-integer value")
		}
		switch {
		case v > 0:
			return Int(+1), nil
		case v < 0:
			return Int(-1), nil
		default:
			return Int(0), nil
		}
	},
	Man: "module: (sign 3) - exec an expression and return the sign",
}

var listExtension = Extension{
	Name: "list",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		var l List
		for _, v := range values {
			l = append(l, v)
		}
		return l, nil
	},
	Man: "module: (list 1 2 (lambda x (add x 1))) - make a list",
}

var appendExtension = Extension{
	Name: "append",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		l, ok := values[0].(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		return append(l, values[1:]...), nil
	},
	Man: "module: (append l 2 (add 1 1)) - append elements into list l and return a new list",
}

var sliceExtension = Extension{
	Name: "slice",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
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
		i, ok := values[1].(Int)
		if !ok {
			return nil, fmt.Errorf("second argument must be integer")
		}
		j, ok := values[2].(Int)
		if !ok {
			return nil, fmt.Errorf("third argument must be integer")
		}
		length := Int(len(l))
		if i < 1 || i > length || j < 1 || j > length {
			return nil, fmt.Errorf("list is out of range")
		}
		return l[i-1 : j], nil
	},
	Man: "module: (slice l 2 3) - make a slice of a list l[2, 3] (list is 1-indexing and slice is a closed interval)",
}

var peekExtension = Extension{
	Name: "peek",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		if len(values) < 2 {
			return nil, fmt.Errorf("peak requires at least 2 arguments")
		}
		l, ok := values[0].(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		length := Int(len(l))
		if length < 1 {
			return nil, fmt.Errorf("empty list")
		}
		var outputs List
		for j := 1; j < len(values); j++ {
			i, ok := values[j].(Int)
			if !ok {
				return nil, fmt.Errorf("second argument must be integer")
			}
			if i < 1 || i > length {
				return nil, fmt.Errorf("list is out of range")
			}
			outputs = append(outputs, l[i-1])
		}
		if len(outputs) == 1 {
			return outputs[0], nil
		}
		return outputs, nil
	},
	Man: "module: (peek l 3 2) - get elem from list (can get multiple elements) (list is 1-indexing)",
}

var lenExtension = Extension{
	Name: "len",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		if len(values) != 1 {
			return nil, fmt.Errorf("len requires 1 argument")
		}
		switch v := values[0].(type) {
		case List:
			return Int(len(v)), nil
		case Dict:
			return Int(len(v)), nil
		default:
			return nil, fmt.Errorf("first argument must be list or dict")
		}
	},
	Man: "module: (len l) - get length of a list of dict",
}

// mapModule - TODO make map parallel by make a copy of the latest frame, reuse other frames, call in parallel
var mapModule = Module{
	Name: "map",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		if len(expr.Args) != 2 {
			return nil, fmt.Errorf("map requires 2 arguments")
		}
		l1, err := r.Step(expr.Args[0], interruptCh)
		if err != nil {
			return nil, err
		}
		l, ok := l1.(List)
		if !ok {
			return nil, fmt.Errorf("first argument must be list")
		}
		f1, err := r.Step(expr.Args[1], interruptCh)
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
				o, err := r.Step(f.Impl, interruptCh)
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
				}, interruptCh)
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
	Man: "module: (map l (lambda y (add 1 y))) - map",
}

// TODO - implement map filter reduce

var typeExtension = Extension{
	Name: "type",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		var types List
		for _, v := range values {
			types = append(types, getType(v))
		}
		if len(types) == 1 {
			return types[0], nil
		}
		return types, nil
	},
	Man: "module: (type x 1 (lambda y (add 1 y))) - get types of objects (can get multiple ones)",
}

var stackModule = Module{
	Name: "stack",
	Exec: func(r *Runtime, expr LambdaExpr, interruptCh <-chan struct{}) (Object, error) {
		var stack List
		for _, f := range r.Stack {
			frame := make(Dict)
			for k, v := range f {
				frame[String(k)] = v
			}
			stack = append(stack, frame)
		}
		return stack, nil
	},
	Man: "module: (stack) - get stack",
}

var printExtension = Extension{
	Name: "print",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {
		for _, v := range values {
			fmt.Printf("%v ", v)
		}
		fmt.Println()
		return Int(len(values)), nil
	},
	Man: "module: (print 1 x (lambda 3)) - print values",
}

var unicodeExtension = Extension{
	Name: "unicode",
	Exec: func(interruptCh <-chan struct{}, values ...Object) (Object, error) {

		var output String = ""
		for _, v := range values {
			if v, ok := v.(Int); ok {
				output += String(fmt.Sprintf("%c", rune(v)))
			} else {
				return nil, fmt.Errorf("argument must be int")
			}
		}
		return output, nil
	},
	Man: "module: (unicode 72 101 108 108 111 44 32 87 111 114 108 100 33) - convert a list of integers into string - this is just for hello world",
}
