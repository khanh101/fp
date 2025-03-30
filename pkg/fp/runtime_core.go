package fp

import (
	"context"
	"errors"
	"fmt"
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
