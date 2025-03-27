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

type Extension = func(r *Runtime, expr LambdaExpr) Object
type Runtime struct {
	debug        bool
	parseLiteral func(lit Name) (Object, error)
	Stack        []Frame
	extension    map[Name]Extension
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

func (r *Runtime) WithExtension(name Name, f Extension) *Runtime {
	r.extension[name] = f
	return r
}

func (r *Runtime) WithParseLiteral(f func(lit Name) (Object, error)) *Runtime {
	r.parseLiteral = f
	return r
}

func (r *Runtime) WithDebug(debug bool) *Runtime {
	r.debug = debug
	return r
}
