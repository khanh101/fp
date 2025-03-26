package fp

import (
	"fmt"
	"os"
	"strconv"
)

const DETECT_NONPURE = true

type Runtime interface {
	Step(expr Expr) Value
}

func NewRuntime() Runtime {
	return &newRuntime{stack: []frame{
		make(frame),
	}}
}

// Value : union of int and lambda - TODO : introduce new data types
type Value interface{}
type lambda struct {
	params []string
	impl   Expr
	frame  frame
}

type frame map[string]Value

func (f frame) update(otherFrame frame) frame {
	for k, v := range otherFrame {
		f[k] = v
	}
	return f
}

type newRuntime struct {
	stack []frame
}

func (r *newRuntime) Step(expr Expr) Value {
	switch expr := expr.(type) {
	case string:
		var v Value
		// convert to number
		v, err := strconv.Atoi(expr)
		if err == nil {
			return v
		}
		for i := len(r.stack) - 1; i >= 0; i-- {
			if v, ok := r.stack[i][expr]; ok {
				if DETECT_NONPURE && i != 0 && i < len(r.stack)-1 {
					_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
				}
				return v
			}
		}
		panic("runtime error")
	case LambdaExpr:
		switch expr.Name {
		case "output":
			for _, arg := range expr.Args {
				v := r.Step(arg)
				fmt.Printf("%v ", v)
			}
			fmt.Println()
			return len(expr.Args)
		case "let":
			name := expr.Args[0].(string)
			v := r.Step(expr.Args[1])
			r.stack[len(r.stack)-1][name] = v
			return v
		case "input":
			name := expr.Args[0].(string)
			var v int
			_, err := fmt.Scanf("%d", &v)
			if err != nil {
				panic(err)
			}
			r.stack[len(r.stack)-1][name] = v
			return v
		case "lambda":
			v := lambda{
				params: nil,
				impl:   nil,
				frame:  nil,
			}
			for i := 0; i < len(expr.Args)-1; i++ {
				paramName := expr.Args[i].(string)
				v.params = append(v.params, paramName)
			}
			v.impl = expr.Args[len(expr.Args)-1]
			v.frame = make(frame).update(r.stack[len(r.stack)-1])
			return v
		case "case":
			cond := r.Step(expr.Args[0])
			i := func() int {
				for i := 1; i < len(expr.Args); i += 2 {
					if arg, ok := expr.Args[i].(string); ok && arg == "_" {
						return i
					}
					if r.Step(expr.Args[i]) == cond {
						return i
					}
				}
				panic("runtime error")
			}()
			return r.Step(expr.Args[i+1])
		case "sign":
			v := r.Step(expr.Args[0]).(int)
			switch {
			case v > 0:
				return +1
			case v < 0:
				return -1
			case v == 0:
				return 0
			}
		case "sub":
			a := r.Step(expr.Args[0]).(int)
			b := r.Step(expr.Args[1]).(int)
			return a - b
		case "add":
			v := 0
			for _, arg := range expr.Args {
				v += r.Step(arg).(int)
			}
			return v
		case "tail":
			var v Value
			for _, arg := range expr.Args {
				v = r.Step(arg)
			}
			return v
		default: // function application
			// 1. get func recursively
			f := func() lambda {
				for i := len(r.stack) - 1; i >= 0; i-- {
					if f, ok := r.stack[i][expr.Name]; ok {
						if DETECT_NONPURE && i != 0 && i < len(r.stack)-1 {
							_, _ = fmt.Fprintf(os.Stderr, "non-pure function")
						}
						return f.(lambda)
					}
				}
				panic("runtime error")
			}()
			// 1. evaluate arguments
			var args []Value
			for _, arg := range expr.Args {
				args = append(args, r.Step(arg))
			}
			// 2. add argument to local frame
			localFrame := make(frame).update(f.frame)
			for i := 0; i < len(f.params); i++ {
				localFrame[f.params[i]] = args[i]
			}
			// 3. push frame to stack
			r.stack = append(r.stack, localFrame)
			// 4. exec function
			v := r.Step(f.impl)
			// 5. pop frame from stack
			r.stack = r.stack[:len(r.stack)-1]
			return v
		}
	default:
		panic("runtime error")
	}
	panic("runtime error")
}
