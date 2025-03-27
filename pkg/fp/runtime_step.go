package fp

import (
	"fmt"
	"os"
)

type StepOption func(*stepOption) *stepOption
type stepOption struct {
	tailCallOptimization bool
}

func TCOStepOption(tco bool) StepOption {
	return func(o *stepOption) *stepOption {
		o.tailCallOptimization = false // TODO - debug tail call optimization
		return o
	}
}

// Step - implement minimal set of instructions for the language to be Turing complete
// let, Lambda, case, sign, sub, add, tail
func (r *Runtime) Step(expr Expr, opts ...StepOption) Object {
	if r.debug {
		defer func() {
			logDebug("%v\n", r)
		}()
	}
	o := &stepOption{
		tailCallOptimization: false,
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		o = opt(o)
	}
	switch expr := expr.(type) {
	case Name:
		var v Object
		// parse name
		v, err := r.parseLiteral(expr)
		if err == nil {
			return v
		}
		// find in stack for variable
		for i := len(r.Stack) - 1; i >= 0; i-- {
			if v, ok := r.Stack[i][expr]; ok {
				if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
					_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
				}
				return v
			}
		}
		panicError("runtime error: variable %s not found", expr.String())
	case LambdaExpr:
		// find in stack for user-defined function
		if f, ok := func() (Lambda, bool) {
			// 1. get func recursively
			for i := len(r.Stack) - 1; i >= 0; i-- {
				if f, ok := r.Stack[i][expr.Name]; ok {
					if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
						_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
					}
					return f.(Lambda), true
				}
			}
			return Lambda{}, false
		}(); ok {
			// 1. evaluate arguments
			args := r.stepWithTailOption(nil, expr.Args...)
			if o.tailCallOptimization {
				// 2. reuse last frame
				for i := 0; i < len(f.Params); i++ {
					r.Stack[len(r.Stack)-1][f.Params[i]] = args[i]
				}
			} else {
				// 2. add argument to local Frame
				localFrame := make(Frame).Update(f.Frame)
				for i := 0; i < len(f.Params); i++ {
					localFrame[f.Params[i]] = args[i]
				}
				// 3. push Frame to Stack
				r.Stack = append(r.Stack, localFrame)
			}
			// 4. exec function
			v := r.Step(f.Impl)
			if o.tailCallOptimization {
				// pass
			} else {
				// 5. pop Frame from Stack
				r.Stack = r.Stack[:len(r.Stack)-1]
			}
			return v
		}
		// check for extension
		if f, ok := r.extension[expr.Name]; ok {
			return f(r, expr)
		}
		panicError("runtime error: function %s not found", expr.Name.String())
	default:
		panicError("runtime error: unknown expression type")
	}
	panicError("unreachable")
	return nil
}

func (r *Runtime) stepWithTailOption(opt StepOption, exprList ...Expr) []Object {
	var outputs []Object
	if len(exprList) > 0 {
		for i := 0; i < len(exprList)-1; i++ {
			outputs = append(outputs, r.Step(exprList[i]))
		}
		outputs = append(outputs, r.Step(exprList[len(exprList)-1], opt))
	}
	return outputs
}
