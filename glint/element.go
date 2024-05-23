package glint

import sitter "github.com/smacker/go-tree-sitter"

// type BaseElem struct {
// 	Name string
// 	Node *sitter.Node
// 	Row int
// 	Col int
// }

// type Function struct {
// 	BaseElem
// 	Return
// }

// type CallExpress struct {
// 	BaseElem
// 	Function *Function
// 	Method   *Method
// }

// type Class struct {
// 	BaseElem
// 	attributes []Elementer
// 	methods    []*Method
// }

// func NewClass(name string, node *sitter.Node) *Class {
// 	return &Class{
// 		BaseElem: BaseElem{
// 			name: name,
// 			node: node,
// 		},
// 	}
// }

// func (c *Class) Methods() []*Method {
// 	return c.Methods()
// }

// func (c *Class) AddMethod(method *Method) {
// 	if method.instance != c {
// 		method.instance = c
// 	}
// 	c.methods = append(c.methods, method)
// }

// func (c *Class) RangeMethod(fn func(method *Method)) {
// 	for index := range c.methods {
// 		fn(c.methods[index])
// 	}
// }

// var _ Elementer = (*Class)(nil)

// type Method struct {
// 	Instance *Class
// 	Ptr      bool
// 	Reciver  string
// 	Function
// }

// func NewMethod(name string, node *sitter.Node, class *Class) *Method {
// 	return &Method{
// 		Instance: class,
// 		Function: Function{
// 			BaseElem: BaseElem{
// 				name: name,
// 				node: node,
// 			},
// 		},
// 	}
// }

// func (m *Method) String() string {
// 	return fmt.Sprintf("<Method: %s loc:%d,%d>", m.name, m.Row(), m.Col())
// }

// type Const struct {
// 	BaseElem
// }

// type Variable struct {
// 	BaseElem
// 	global bool
// }

type InnerNode struct {
	Node *sitter.Node
}

func (i *InnerNode) Row() int {
	return int(i.Node.StartPoint().Row)
}

func (i *InnerNode) Col() int {
	return int(i.Node.StartPoint().Column)
}

type Type struct {
	Name    string
	Fields  []*Field
	Methods []*Method
	File    string
	InnerNode
}

func (t *Type) RangeMethod(fn func(method *Method)) {
	for index := range t.Methods {
		fn(t.Methods[index])
	}
}

func (t *Type) RangeField(fn func(field *Field)) {
	for index := range t.Fields {
		fn(t.Fields[index])
	}
}

type Method struct {
	InnerNode
	Ptr     bool
	Reciver string
	Class   Type
	Name    string
	File    string
	Args    []*Argument
	Return  []*Type
}

type Argument struct {
	InnerNode
	// 形参
	Symbol string
	// 形参类型
	Type Type
}

type Function struct {
	InnerNode
	Name   string
	File   string
	Args   []*Argument
	Return []*Type
	Call   []*Call
}

type Value struct {
	InnerNode
	Name  string
	Type  *Type
	Value string
	File  string
}

type Field struct {
	InnerNode
	Name string
	Type *Type
}

type Call struct {
	InnerNode
	Function *Function
	Args     map[string]string
}
