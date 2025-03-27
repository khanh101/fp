package fp

import (
	"fmt"
)

const DETECT_NONPURE = true

// Object : object union of int, string, Lambda - TODO : introduce new data types
type Object interface{}
type Lambda struct {
	Params []Name
	Impl   Expr
	Frame  Frame
}

func (l Lambda) String() string {
	return l.Impl.String()
}

type Frame map[Name]Object

func (f Frame) Update(otherFrame Frame) Frame {
	for k, v := range otherFrame {
		f[k] = v
	}
	return f
}

type Module = func(r *Runtime, expr LambdaExpr) (Object, error)
type Runtime struct {
	debug        bool
	parseLiteral func(lit Name) (Object, error)
	Stack        []Frame
	Module       map[Name]Module
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

func (r *Runtime) LoadModule(name Name, f Module) *Runtime {
	r.Module[name] = f
	return r
}

func (r *Runtime) LoadParseLiteral(f func(lit Name) (Object, error)) *Runtime {
	r.parseLiteral = f
	return r
}

func (r *Runtime) SetDebug(debug bool) *Runtime {
	r.debug = debug
	return r
}
