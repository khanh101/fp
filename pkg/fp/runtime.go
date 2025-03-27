package fp

import (
	"fmt"
	"os"
	"strconv"
)

const DETECT_NONPURE = true

func NewRuntime() *Runtime {
	return (&Runtime{
		Stack: []Frame{
			make(Frame),
		},
		systemExtension: make(map[Name]func(r *Runtime, expr LambdaExpr) Value),
		userExtension:   make(map[Name]func(...Value) Value),
	}).WithSystemExtension("let", func(r *Runtime, expr LambdaExpr) Value {
		name := expr.Args[0].(Name)
		var v Value
		for i := 1; i < len(expr.Args); i++ {
			if i == len(expr.Args)-1 {
				v = r.Step(expr.Args[i], WithTailCall)
			} else {
				v = r.Step(expr.Args[i])
			}
		}
		r.Stack[len(r.Stack)-1][name] = v
		return v
	}).WithSystemExtension("lambda", func(r *Runtime, expr LambdaExpr) Value {
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
	}).WithSystemExtension("case", func(r *Runtime, expr LambdaExpr) Value {
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
		return r.Step(expr.Args[i+1], WithTailCall)
	}).WithSystemExtension("sign", func(r *Runtime, expr LambdaExpr) Value {
		v := r.Step(expr.Args[0], WithTailCall).(int)
		switch {
		case v > 0:
			return +1
		case v < 0:
			return -1
		case v == 0:
			return 0
		}
		panic("runtime error")
	}).WithSystemExtension("sub", func(r *Runtime, expr LambdaExpr) Value {
		a := r.Step(expr.Args[0]).(int)
		b := r.Step(expr.Args[1], WithTailCall).(int)
		return a - b
	}).WithSystemExtension("add", func(r *Runtime, expr LambdaExpr) Value {
		var v int
		for i := 0; i < len(expr.Args); i++ {
			if i == len(expr.Args)-1 {
				v += r.Step(expr.Args[i], WithTailCall).(int)
			} else {
				v += r.Step(expr.Args[i]).(int)
			}
		}
		return v
	}).WithSystemExtension("tail", func(r *Runtime, expr LambdaExpr) Value {
		var v Value
		for i := 0; i < len(expr.Args); i++ {
			if i == len(expr.Args)-1 {
				v = r.Step(expr.Args[i], WithTailCall)
			} else {
				v = r.Step(expr.Args[i])
			}
		}
		return v
	})
}

func (r *Runtime) WithExtension(name Name, f func(...Value) Value) *Runtime {
	r.userExtension[name] = f
	return r
}

func (r *Runtime) WithSystemExtension(name Name, f func(r *Runtime, expr LambdaExpr) Value) *Runtime {
	r.systemExtension[name] = f
	return r
}

type Runtime struct {
	Stack           []Frame
	systemExtension map[Name]func(r *Runtime, expr LambdaExpr) Value
	userExtension   map[Name]func(...Value) Value
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

type callOption struct {
	tailCall bool
}

func WithTailCall(o *callOption) *callOption {
	o.tailCall = false // TODO - debug tailcall
	return o
}

// Step - implement minimal set of instructions for the language to be Turing complete
// let, Lambda, case, sign, sub, add, tail
func (r *Runtime) Step(expr Expr, callOptions ...func(*callOption) *callOption) Value {
	o := &callOption{
		tailCall: false,
	}
	for _, opt := range callOptions {
		o = opt(o)
	}
	switch expr := expr.(type) {
	case Name:
		var v Value
		// convert to number
		v, err := strconv.Atoi(string(expr))
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
		// check for systemExtension
		if f, ok := r.systemExtension[expr.Name]; ok {
			return f(r, expr)
		}
		// check for userExtension
		if f, ok := r.userExtension[expr.Name]; ok {
			var args []Value
			for _, arg := range expr.Args {
				args = append(args, r.Step(arg))
			}
			return f(args...)
		}
		// user-defined function application
		// 1. get func recursively
		f := func() Lambda {
			for i := len(r.Stack) - 1; i >= 0; i-- {
				if f, ok := r.Stack[i][expr.Name]; ok {
					if DETECT_NONPURE && i != 0 && i < len(r.Stack)-1 {
						_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
					}
					return f.(Lambda)
				}
			}
			panic("runtime error")
		}()
		// 1. evaluate arguments
		var args []Value
		for _, arg := range expr.Args {
			args = append(args, r.Step(arg))
		}
		if o.tailCall {
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
		if o.tailCall {
			// pass
		} else {
			// 5. pop Frame from Stack
			r.Stack = r.Stack[:len(r.Stack)-1]
		}
		return v
	default:
		panic("runtime error")
	}
}
