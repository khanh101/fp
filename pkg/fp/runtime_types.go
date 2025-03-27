package fp

import "fmt"

// types - TODO implement custom data types like Int, List, Dict

// Object : object union of int, Lambda, Module, List - TODO : introduce new data types
type Object interface{}

func getType(o Object) String {
	switch o.(type) {
	case int:
		return "int"
	case String:
		return "String"
	case Lambda:
		return "Lambda"
	case Module:
		return "Module"
	case List:
		return "List"
	default:
		return "unknown"
	}
}

type String string

func (s String) String() string {
	return string(s)
}

type Lambda struct {
	Params []Name `json:"params,omitempty"`
	Impl   Expr   `json:"impl,omitempty"`
	Frame  Frame  `json:"frame,omitempty"`
}

func (l Lambda) String() string {
	return l.Impl.String()
}

type Module func(r *Runtime, expr LambdaExpr) (Object, error)

func (m Module) String() string {
	return fmt.Sprintf("[module %p]", m)
}

type List []Object

func (l List) String() string {
	s := ""
	s += "["
	for _, obj := range l {
		s += fmt.Sprintf("%v,", obj)
	}
	s += "]"
	return s
}
