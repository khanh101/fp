package fp

import (
	"fmt"
)

type Runtime struct {
	parseLiteral func(lit Name) (Object, error)
	Stack        []Frame         `json:"stack,omitempty"`
	Module       map[Name]Module `json:"module,omitempty"`
}
type Frame map[Name]Object

func (f Frame) Update(otherFrame Frame) Frame {
	for k, v := range otherFrame {
		f[k] = v
	}
	return f
}

type Module = func(r *Runtime, expr LambdaExpr) (Object, error)

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

func (r *Runtime) LoadModule(name Name, f Module) *Runtime {
	r.Module[name] = f
	return r
}

func (r *Runtime) LoadParseLiteral(f func(lit Name) (Object, error)) *Runtime {
	r.parseLiteral = f
	return r
}
