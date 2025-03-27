package fp

import (
	"fmt"
	"os"
	"strconv"
)

const DETECT_NONPURE = true

// NewPlainRuntime - language specification
func NewPlainRuntime() *Runtime {
	return (&Runtime{
		Stack: []Frame{
			make(Frame),
		},
		parseToken: func(expr string) (interface{}, error) {
			return strconv.Atoi(expr)
		},
		extension: make(map[Name]func(r *Runtime, expr LambdaExpr) Value),
	}).WithExtension("let", func(r *Runtime, expr LambdaExpr) Value {
		name := expr.Args[0].(Name)
		var v Value
		for i := 1; i < len(expr.Args); i++ {
			if i == len(expr.Args)-1 {
				v = r.Step(expr.Args[i], WithTailCallOptimization)
			} else {
				v = r.Step(expr.Args[i])
			}
		}
		r.Stack[len(r.Stack)-1][name] = v
		return v
	}).WithExtension("lambda", func(r *Runtime, expr LambdaExpr) Value {
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
	}).WithExtension("case", func(r *Runtime, expr LambdaExpr) Value {
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
	})
}

// NewBasicRuntime : minimal set of extensions for Turing completeness
func NewBasicRuntime() *Runtime {
	return NewPlainRuntime().WithArithmeticExtension("sign", func(value ...Value) Value {
		v := value[len(value)-1].(int)
		switch {
		case v > 0:
			return +1
		case v < 0:
			return -1
		case v == 0:
			return 0
		}
		panic("runtime error")
	}).WithArithmeticExtension("tail", func(value ...Value) Value {
		return value[len(value)-1]
	}).WithArithmeticExtension("sub", func(value ...Value) Value {
		if len(value) != 2 {
			panic("runtime error")
		}
		return value[0].(int) - value[1].(int)
	}).WithArithmeticExtension("add", func(value ...Value) Value {
		v := 0
		for i := 0; i < len(value); i++ {
			v += value[i].(int)
		}
		return v
	})
}

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

func (r *Runtime) WithExtension(name Name, f func(r *Runtime, expr LambdaExpr) Value) *Runtime {
	r.extension[name] = f
	return r
}

type Runtime struct {
	Stack      []Frame
	parseToken func(Token) (interface{}, error)
	extension  map[Name]func(r *Runtime, expr LambdaExpr) Value
}

// Value : union of int, string, Lambda - TODO : introduce new data types
type Value interface{}
type Lambda struct {
	Params []Name
	Impl   Expr
	Frame  Frame
}

type Frame map[Name]Value

func (f Frame) Update(otherFrame Frame) Frame {
	for k, v := range otherFrame {
		f[k] = v
	}
	return f
}

type stepOption struct {
	tailCallOptimization bool
}

func WithTailCallOptimization(o *stepOption) *stepOption {
	o.tailCallOptimization = false // TODO - debug tailcall
	return o
}

// Step - implement minimal set of instructions for the language to be Turing complete
// let, Lambda, case, sign, sub, add, tail
func (r *Runtime) Step(expr Expr, stepOptions ...func(*stepOption) *stepOption) Value {
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
	panic("runtime error")
}
