package fp

import (
	"fmt"
	"os"
)

const DETECT_NONPURE = true

func (r *Runtime) getFromStack(name Name) (Object, error) {
	for i := len(r.Stack) - 1; i >= 0; i-- {
		if o, ok := r.Stack[i][name]; ok {
			if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
				_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
			}
			return o, nil
		}
	}
	return nil, fmt.Errorf("object not found %s", name)
}

func (r *Runtime) getFuncOrModule(name Name) (Object, error) {
	// find in stack for user-defined function or module
	f, ok, err := func() (Object, bool, error) {
		// 1. get func recursively
		for i := len(r.Stack) - 1; i >= 0; i-- {
			if f, ok := r.Stack[i][name]; ok {
				if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
					_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
				}
				switch f := f.(type) {
				case Lambda, Module:
					return f, true, nil
				default:
					return nil, false, fmt.Errorf("unexpected type %T", f)
				}
			}
		}
		return nil, false, nil
	}()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("function not found for %s", name)
	}
	return f, nil
}

func (r *Runtime) getVar(name Name) (Object, error) {
	for i := len(r.Stack) - 1; i >= 0; i-- {
		if v, ok := r.Stack[i][name]; ok {
			if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
				_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
			}
			return v, nil
		}
	}
	return nil, fmt.Errorf("runtime error: variable %s not found", name.String())
}

// Step - implement minimal set of instructions for the language to be Turing complete
// let, Lambda, case, sign, sub, add, tail
func (r *Runtime) Step(expr Expr) (Object, error) {
	switch expr := expr.(type) {
	case Name:
		var v Object
		// parse name
		v, err := r.parseLiteral(expr)
		if err == nil {
			return v, nil
		}
		// find in stack for variable
		return r.getFromStack(expr)

	case LambdaExpr:
		f, err := r.getFromStack(expr.Name)
		if err != nil {
			return nil, err
		}
		switch f := f.(type) {
		case Lambda:
			// 1. evaluate arguments
			args, err := r.stepMany(expr.Args...)
			if err != nil {
				return nil, err
			}
			// 2. add argument to local Frame
			localFrame := make(Frame).Update(f.Frame)
			for i := 0; i < len(f.Params); i++ {
				localFrame[f.Params[i]] = args[i]
			}
			// 3. push Frame to Stack
			r.Stack = append(r.Stack, localFrame)
			// 4. exec function
			v, err := r.Step(f.Impl)
			if err != nil {
				return nil, err
			}
			// 5. pop Frame from Stack
			r.Stack = r.Stack[:len(r.Stack)-1]
			return v, nil
		case Module:
			return f(r, expr)
		default:
			return nil, fmt.Errorf("function or module %s found but wrong type", expr.Name.String())
		}
	default:
		return nil, fmt.Errorf("runtime error: unknown expression type")
	}
}

func (r *Runtime) stepMany(exprList ...Expr) ([]Object, error) {
	var outputs []Object
	for _, expr := range exprList {
		v, err := r.Step(expr)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, v)
	}
	return outputs, nil
}
