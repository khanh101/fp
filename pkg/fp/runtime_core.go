package fp

import (
	"fmt"
	"os"
)

const DETECT_NONPURE = true

func (r *Runtime) WithExtension(name Name, f Extension) *Runtime {
	r.extension[name] = f
	return r
}

type Extension = func(r *Runtime, expr LambdaExpr) Value
type Runtime struct {
	Stack      []Frame
	parseToken func(Token) (interface{}, error)
	extension  map[Name]Extension
}

// Value : union of int, string, Lambda - TODO : introduce new data types
type Value interface{}
type Lambda struct {
	Params []Name
	Impl   Expr
	Frame  Frame
}

func (l Lambda) String() string {
	return l.Impl.String()
}

type Frame map[Name]Value

func (f Frame) Update(otherFrame Frame) Frame {
	for k, v := range otherFrame {
		f[k] = v
	}
	return f
}

// Step - implement minimal set of instructions for the language to be Turing complete
// let, Lambda, case, sign, sub, add, tail
func (r *Runtime) Step(expr Expr, stepOptions ...StepOption) Value {
	o := &stepOption{
		tailCallOptimization: false,
	}
	for _, opt := range stepOptions {
		o = opt(o)
	}
	switch expr := expr.(type) {
	case Name:
		var v Value
		// parse token
		v, err := r.parseToken(string(expr))
		if err == nil {
			return v
		}
		for i := len(r.Stack) - 1; i >= 0; i-- {
			if v, ok := r.Stack[i][expr]; ok {
				if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
					_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
				}
				return v
			}
		}
		panic("runtime error")
	case LambdaExpr:
		// check for user-defined function
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
			var args []Value
			for _, arg := range expr.Args {
				args = append(args, r.Step(arg))
			}
			if o.tailCallOptimization {
				// tail call - use last frame
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
		panic("runtime error")

	default:
		panic("runtime error")
	}
}
