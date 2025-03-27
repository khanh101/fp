package fp

import (
	"fmt"
	"os"
)

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
		for i := len(r.Stack) - 1; i >= 0; i-- {
			if v, ok := r.Stack[i][expr]; ok {
				if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
					_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
				}
				return v, nil
			}
		}
		return nil, fmt.Errorf("runtime error: variable %s not found", expr.String())
	case LambdaExpr:
		// find in stack for user-defined function
		f, ok, err := func() (Lambda, bool, error) {
			// 1. get func recursively
			for i := len(r.Stack) - 1; i >= 0; i-- {
				if f, ok := r.Stack[i][expr.Name]; ok {
					if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
						_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
					}
					f, ok := f.(Lambda)
					if !ok {
						return Lambda{}, false, fmt.Errorf("first argument in S-expression is not a Lambda")
					}
					return f, true, nil
				}
			}
			return Lambda{}, false, nil
		}()
		if err != nil {
			return nil, err
		}
		if ok {
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
		}
		// check for Module
		if f, ok := r.Module[expr.Name]; ok {
			return f(r, expr)
		}
		return nil, fmt.Errorf("runtime error: function %s not found", expr.Name.String())
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
