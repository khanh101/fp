package fp

// types - TODO implement custom data types like Int, List, Dict

// Object : object union of int, string, Lambda - TODO : introduce new data types
type Object interface{}
type Lambda struct {
	Params []Name `json:"params,omitempty"`
	Impl   Expr   `json:"impl,omitempty"`
	Frame  Frame  `json:"frame,omitempty"`
}

func (l Lambda) String() string {
	return l.Impl.String()
}

type List = []Object
