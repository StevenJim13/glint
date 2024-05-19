package glint

import sitter "github.com/smacker/go-tree-sitter"

type Elementer interface {
	Name() string
	Node() *sitter.Node
	Row() int
	Col() int
}

type BaseElem struct {
	name string
	node *sitter.Node
}

// Col implements Elementer.
func (b *BaseElem) Col() int {
	return int(b.node.StartPoint().Column)
}

// Row implements Elementer.
func (b *BaseElem) Row() int {
	return int(b.node.StartPoint().Row)
}

// Name implements Elementer.
func (b *BaseElem) Name() string {
	return b.name
}

// Node implements Elementer.
func (b *BaseElem) Node() *sitter.Node {
	return b.node
}

var _ Elementer = (*BaseElem)(nil)

type Function struct {
	BaseElem
}

type CallExpress struct {
	BaseElem
	Function *Function
	Method   *Method
}

type Class struct {
	BaseElem
	attributes []Elementer
	methods    []*Method
}

func NewClass(name string, node *sitter.Node) *Class {
	return &Class{
		BaseElem: BaseElem{
			name: name,
			node: node,
		},
	}
}

func (c *Class) Methods() []*Method {
	return c.Methods()
}

func (c *Class) AddMethod(method *Method) {
	if method.instance != c {
		method.instance = c
	}
	c.methods = append(c.methods, method)
}

func (c *Class) RangeMethod(fn func(method *Method)) {
	for index := range c.methods {
		fn(c.methods[index])
	}
}

var _ Elementer = (*Class)(nil)

type Method struct {
	instance *Class
	prt      bool
	reciver  string
	Function
}

func NewMethod(name string, node *sitter.Node, class *Class) *Method {
	return &Method{
		instance: class,
		Function: Function{
			BaseElem: BaseElem{
				name: name,
				node: node,
			},
		},
	}
}

func (m *Method) String() string {
	return fmt.Sprintf("<Method: %s loc:%d,%d>", m.name, m.Row(), m.Col())
}

type Const struct {
	BaseElem
}

type Variable struct {
	BaseElem
	global bool
}
