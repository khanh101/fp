package fp

import "fmt"

// types - TODO implement custom data types like Int, List, Dict

// Object : object union of int, Lambda, Module, List, Wildcard - TODO : introduce new data types
type Object interface {
	String() string
}

func getType(o Object) String {
	switch o.(type) {
	case Int:
		return "Int"
	case String:
		return "String"
	case Lambda:
		return "Lambda"
	case Module:
		return "Module"
	case List:
		return "List"
	case Wildcard:
		return "wildcard"
	default:
		return "unknown"
	}
}

type Int int

func (i Int) String() string {
	return fmt.Sprintf("%d", i)
}

type Dict map[Object]Object

func (d Dict) String() string {

	s := ""
	s += "{"
	for k, v := range d {
		s += fmt.Sprintf("%s -> %s,", k.String(), v.String())
	}
	s += "}"
	return s
}

type Wildcard struct{}

func (w Wildcard) String() string {
	return "_"
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
	s := "(lambda "
	for _, param := range l.Params {
		s += param.String() + " "
	}
	s += l.Impl.String()
	s += ")"
	return s
}

type Module struct {
	Exec func(r *Runtime, expr LambdaExpr) (Object, error)
	Man  string `json:"man,omitempty"`
}

func (m Module) String() string {
	return m.Man
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
