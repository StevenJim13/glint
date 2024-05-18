package glint

import (
	"context"
	"fmt"
	"os"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stkali/glint/ast"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
)

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

type Var struct {
	BaseElem
}

// Context file content
type Context interface {
	// Content returns file content
	Content() []byte
	// LinesInfo returns line information of file
	LinesInfo() LinesInfo
	// Functions returns all function define AST node(s) of file
	File() string
	// Lint()
	Lint(Outputer, Context)
	// AddDefect
	Defect(*Model, int, int, string, ...any)
	DefectSet() []*Defect
	ASTTree() *sitter.Tree
	FileRoot() Context
	Parent() Context
	SetParent(Context)
	AddChild(Context)
	RangeChildren(func(ctx Context))
	IsDir() bool
	// Functions
	Functions() []*Function
	// CallExpress
	CallExpresses() []*CallExpress
	//
	Classes() map[string]*Class
	Variables() []*Var
	Consts() []*Const
}

// BaseContext ...
type BaseContext struct {
	// 文件的语言和检查规则
	*Linter
	// 文件路径
	file string
	// 子节点
	children []Context
	// 文件的内容
	content []byte
	// 文件的行信息
	info [][2]int
	// defects
	defects []*Defect
	//
	ast *sitter.Tree
	// parent
	parent Context
	// root file node
	root Context
	//
	funcs []*Function
	//
	classes map[string]*Class
	//
	vars   []*Var
	consts []*Const

	calls []*CallExpress
	sync.Once
}

// IsDir implements Context.
func (b *BaseContext) IsDir() bool {
	return b.Linter == nil
}

// RangeChildren implements Context.
func (b *BaseContext) RangeChildren(fn func(ctx Context)) {
	for index := range b.children {
		fn(b.children[index])
	}
}

// CreateBaseContext 为了测试需要
func CreateTestBaseContext(
	lang utils.Language,
	file string,
	content []byte,
) *BaseContext {
	node := &BaseContext{
		file:    file,
		content: content,
		Linter:  &Linter{Lang: lang},
		classes: make(map[string]*Class),
	}
	return node
}

// ASTTree implements Context.
func (b *BaseContext) ASTTree() *sitter.Tree {
	if b.ast == nil {
		parser := sitter.NewParser()
		parser.SetLanguage(b.Lang.Lang())
		tree, err := parser.ParseCtx(context.TODO(), nil, b.Content())
		if err != nil {
			log.Errorf("failed to parse file %q ast, err:%s", b.file, err)
		} else {
			b.ast = tree
		}
	}
	return b.ast
}

// FileTree implements Context.
func (b *BaseContext) FileRoot() Context {
	if b.root == nil {
		b.root = b.getFileRoot()
	}
	return b.root
}

func (b *BaseContext) getFileRoot() Context {
	var tmp Context
	tmp = b
	for tmp.Parent() != nil {
		tmp = tmp.Parent()
	}
	return tmp
}

// Parent implements Context.
func (b *BaseContext) Parent() Context {
	return b.parent
}

func (b *BaseContext) SetParent(parent Context) {
	b.parent = parent
}

// DefectSet implements Context.
func (b *BaseContext) DefectSet() []*Defect {
	return b.defects
}

// Defect implements Context.
func (b *BaseContext) Defect(model *Model, row int, col int, s string, args ...any) {
	def := &Defect{
		Model: model,
		Desc:  fmt.Sprintf(s, args...),
		Row:   row,
		Col:   col,
	}
	b.defects = append(b.defects, def)
}

// CallExpresses implements Context.
func (b *BaseContext) CallExpresses() []*CallExpress {
	errors.Warning("baseContext no callexpresses")
	return nil
}

// Functions implements Context.
func (b *BaseContext) Functions() []*Function {
	errors.Warning("baseContext no functions")
	return nil
}

func (b *BaseContext) Classes() map[string]*Class {
	errors.Warning("baseContext no classes")
	return nil
}

func (b *BaseContext) Variables() []*Var {
	errors.Warning("baseContext no Variables")
	return nil
}

func (b *BaseContext) Consts() []*Const {
	errors.Warning("baseContext no Consts")
	return nil
}

// Lint implements Context.
func (b *BaseContext) Lint(output Outputer, ctx Context) {
	if b.Linter != nil && b.LintFunc != nil {
		b.LintFunc(output, ctx)
	}
}

// Content implements Context.
func (b *BaseContext) Content() []byte {
	if b.content == nil {
		if err := b.loadContent(); err != nil {
			errors.Warningf("failed to get file: %q content, err: %s", b.file, err)
		}
	}
	return b.content
}

// loadContent TODO
func (b *BaseContext) loadContent() (err error) {
	b.content, err = os.ReadFile(b.file)
	return
}

// LinesInfo implements Context.
//
//	0  -   \r
//	1  -   \n
//	2  -   \r\n
func (b *BaseContext) LinesInfo() LinesInfo {
	if b.info == nil {
		ctt := b.Content()
		gap := 0
		index, length := 0, len(ctt)
		for index < length {
			switch ctt[index] {
			case '\r':
				lineLength := len(tool.ToString(ctt[gap:index]))
				if index+1 < length {
					if ctt[index+1] == '\n' {
						// \r\n
						b.info = append(b.info, [2]int{lineLength, 2})
						index += 1
						gap = index + 1
					} else {
						// \r
						b.info = append(b.info, [2]int{lineLength, 0})
						gap = index + 1
					}
				} else {
					// EOF
					b.info = append(b.info, [2]int{lineLength, 0})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				b.info = append(b.info, [2]int{lineLength, 1})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			b.info = append(b.info, [2]int{lineLength, -1})
		}
	}
	return b.info
}

// Name implements Context.
func (b *BaseContext) File() string {
	return b.file
}

func (b *BaseContext) AddChild(node Context) {
	node.SetParent(b)
	b.children = append(b.children, node)
}

func (b *BaseContext) String() string {
	return fmt.Sprintf("<Node: %s>", b.File())
}

var _ Context = (*BaseContext)(nil)

var _ Context = (*GolangContext)(nil)

// BaseContext

// Golang Context

type GolangContext struct {
	BaseContext
}

// CallExpresses implements Context.
// Subtle: this method shadows the method (*BaseContext).CallExpresses of GolangContext.BaseContext.
func (g *GolangContext) CallExpresses() []*CallExpress {
	return g.calls
}

// Functions implements Context.
// Subtle: this method shadows the method (*BaseContext).Functions of GolangContext.BaseContext.
func (g *GolangContext) Functions() []*Function {
	g.parseElements()
	return g.funcs
}

func (g *GolangContext) Classes() map[string]*Class {
	g.parseElements()
	return g.classes
}

// Functions implements Context.
// Subtle: this method shadows the method (*BaseContext).Functions of GolangContext.BaseContext.
func (g *GolangContext) Variables() []*Var {
	g.parseElements()
	return g.vars
}

// Functions implements Context.
// Subtle: this method shadows the method (*BaseContext).Functions of GolangContext.BaseContext.
func (g *GolangContext) Consts() []*Const {
	g.parseElements()
	return g.consts
}

func (g *GolangContext) parseMethod(node *sitter.Node, content []byte) *Method {

	method := &Method{}
	if nameNode := node.ChildByFieldName("name"); nameNode != nil {
		method.name = nameNode.Content(content)
		method.node = node
	}
	if paramNode := node.ChildByFieldName("receiver"); paramNode.Type() == "parameter_list" {
		if decl := paramNode.Child(1); decl.Type() == "parameter_declaration" {
			if reciverNode := decl.ChildByFieldName("name"); reciverNode != nil {
				method.reciver = reciverNode.Content(content)
			}
			typeNode := decl.ChildByFieldName("type")
			if typeNode == nil {
				return nil
			}
			switch typeNode.Type() {
			case "type_identifier": // 绑定到了实例上
				method.instance = NewClass(typeNode.Content(content), nil)
			case "pointer_type": // 绑定到地址
				if !handlePointType(method, typeNode, content) {
					return nil
				}
			case "generic_type":
				if !handleGenericType(method, typeNode, content) {
					return nil
				}
			default:
				return nil
			}
		}
	}
	return method
}

func handleGenericType(method *Method, node *sitter.Node, content []byte) bool {
	if identifier := node.ChildByFieldName("type"); identifier != nil && identifier.Type() == "type_identifier" {
		method.instance = NewClass(identifier.Content(content), nil)
		return true
	}
	return false
}

func handlePointType(method *Method, node *sitter.Node, content []byte) bool {
	method.prt = true
	first := node.Child(1)
	if first == nil {
		return false
	}
	switch first.Type() {
	case "generic_type":
		return handleGenericType(method, first, content)
	case "type_identifier":
		method.instance = NewClass(first.Content(content), nil)
	default:
		return false
	}
	return true
}

func (g *GolangContext) parseElements() {

	g.Do(func() {
		content := g.Content()
		ast.ApplyChildrenNodes(g.ASTTree().RootNode(), func(sub *sitter.Node) {
			switch sub.Type() {
			case "method_declaration":
				method := g.parseMethod(sub, content)
				if method == nil {
					utils.Bugf("cannot parse method: %s", sub.Content(content))
				} else {
					if class, ok := g.classes[method.instance.name]; !ok {
						g.classes[method.instance.name] = NewClass(method.instance.name, nil)
					} else {
						class.AddMethod(method)
					}
				}

			case "const_declaration":
				ast.ApplyChildrenNodes(sub, func(spec *sitter.Node) {
					if spec.Type() == "const_spec" {
						nameNode := spec.ChildByFieldName("name")
						if nameNode != nil {
							g.consts = append(g.consts, &Const{BaseElem{name: nameNode.Content(content), node: spec}})
						} else {
							utils.Bugf("cannot parse const desclare: %s", sub.Content(content))
						}
					}
				})

			case "var_declaration":
				ast.ApplyChildrenNodes(sub, func(spec *sitter.Node) {
					if spec.Type() == "var_spec" {
						nameNode := spec.ChildByFieldName("name")
						if nameNode != nil {
							g.consts = append(g.consts, &Const{BaseElem{name: nameNode.Content(content), node: spec}})
						} else {
							utils.Bugf("cannot parse var desclare: %s", sub.Content(content))
						}
					}
				})
			case "function_declaration":
				if nameNode := sub.ChildByFieldName("name"); nameNode.Type() == "identifier" {
					g.funcs = append(g.funcs, &Function{BaseElem{name: nameNode.Content(content), node: sub}})
				} else {
					utils.Bugf("cannot parse function desclare: %s", sub.Content(content))
				}
			case "type_declaration":
				if spec := sub.Child(1); spec.Type() == "type_spec" {
					if nameNode := spec.ChildByFieldName("name"); nameNode.Type() == "type_identifier" {
						name := nameNode.Content(content)
						if class, ok := g.classes[name]; !ok {
							g.classes[name] = NewClass(name, sub)
						} else {
							class.node = sub
						}
					} else {
						utils.Bugf("cannot parse struct declare: %s", sub.Content(content))
					}
				}
			}
		})
	})

}

func NewGlangContext(file string, linter *Linter) Context {
	baseContext := BaseContext{
		file:    file,
		Linter:  linter,
		classes: make(map[string]*Class),
	}
	return &GolangContext{
		BaseContext: baseContext,
	}
}

// NewPythonContext create a context of python script file
func NewPythonContext(file string, linter *Linter) Context {
	return nil
}

func NewCCppContext(file string, linter *Linter) Context {
	return nil
}

func CreateContext(file string, linter *Linter) Context {
	if linter == nil {
		return &BaseContext{
			file:    file,
			classes: make(map[string]*Class),
		}
	}
	switch linter.Lang {
	case utils.GoLang:
		return NewGlangContext(file, linter)
	case utils.Python:
		return NewPythonContext(file, linter)
	case utils.CCpp:
		return NewCCppContext(file, linter)
	}
	return nil
}
