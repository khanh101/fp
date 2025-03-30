package fp

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type Runtime struct {
	parseLiteral func(lit String) (Object, error)
	Stack        []Frame `json:"stack,omitempty"`
}
type Frame map[String]Object

func (f Frame) Update(otherFrame Frame) Frame {
	for k, v := range otherFrame {
		f[k] = v
	}
	return f
}

func (r *Runtime) String() string {
	s := ""
	for i, f := range r.Stack {
		s += "["
		for k, v := range f {
			s += fmt.Sprintf("%s -> %v, ", k, v)
		}
		if i != len(r.Stack)-1 {
			s += "]\n"
		} else {
			s += "]"
		}
	}
	return s
}

func (r *Runtime) LoadModule(m Module) *Runtime {
	r.Stack[0][m.Name] = m
	return r
}

func (r *Runtime) LoadParseLiteral(f func(lit String) (Object, error)) *Runtime {
	r.parseLiteral = f
	return r
}

type Extension struct {
	Name String
	Exec func(ctx context.Context, values ...Object) (Object, error)
	Man  string
}

func (r *Runtime) LoadExtension(e Extension) *Runtime {
	return r.LoadModule(Module{
		Name: e.Name,
		Exec: func(ctx context.Context, r *Runtime, expr LambdaExpr) (Object, error) {
			args, err := r.stepMany(ctx, expr.Args...)
			if err != nil {
				return nil, err
			}
			var unwrappedArgs []Object
			i := 0
			for i < len(args) {
				if _, ok := args[i].(Unwrap); ok {
					if i+1 >= len(args) {
						return nil, errors.New("unwrapping arguments must be a list")
					}
					argsList, ok := args[i+1].(List)
					if !ok {
						return nil, errors.New("unwrapping arguments must be a list")
					}
					for _, elem := range argsList {
						unwrappedArgs = append(unwrappedArgs, elem)
					}
					i += 2
				} else {
					unwrappedArgs = append(unwrappedArgs, args[i])
					i++
				}
			}
			return e.Exec(ctx, unwrappedArgs...)
		},
		Man: e.Man,
	})
}

const DETECT_NONPURE = true

func (r *Runtime) searchOnStack(name String) (Object, error) {
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

type Interrupt struct{}

func (i Interrupt) Error() string {
	return "interrupt"
}

var InterruptError = Interrupt{}

// Step - implement minimal set of instructions for the language to be Turing complete
// let, Lambda, case, sign, sub, add, tail
func (r *Runtime) Step(ctx context.Context, expr Expr) (Object, error) {
	// TODO - get step option from context here -
	// TODO - something is like - parallel, tail_call_optimization, error, or deadline, or implement my own context class
	select {
	case <-ctx.Done():
		return nil, InterruptError
	default:
		switch expr := expr.(type) {
		case Name:
			var v Object
			// parse name
			v, err := r.parseLiteral(String(expr))
			if err == nil {
				return v, nil
			}
			// find in stack for variable
			return r.searchOnStack(String(expr))

		case LambdaExpr:
			f, err := r.searchOnStack(String(expr.Name))
			if err != nil {
				return nil, err
			}
			switch f := f.(type) {
			case Lambda:
				// 1. evaluate arguments
				args, err := r.stepMany(ctx, expr.Args...)
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
				v, err := r.Step(ctx, f.Impl)
				if err != nil {
					return nil, err
				}
				// 5. pop Frame from Stack
				r.Stack = r.Stack[:len(r.Stack)-1]
				return v, nil
			case Module:
				return f.Exec(ctx, r, expr)
			default:
				return nil, fmt.Errorf("function or module %s found but wrong type", expr.Name.String())
			}
		default:
			return nil, fmt.Errorf("runtime error: unknown expression type")
		}
	}
}

func (r *Runtime) stepMany(ctx context.Context, exprList ...Expr) ([]Object, error) {
	var outputs []Object
	for _, expr := range exprList {
		v, err := r.Step(ctx, expr)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, v)
	}
	return outputs, nil
}
