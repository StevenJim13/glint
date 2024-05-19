/*
glint
j
*/
package glint

import (
	"os"
	"sync"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/tool"
)

// // Context file content
// type Context interface {
// 	// Content returns file content
// 	Content() []byte
// 	// LinesInfo returns line information of file
// 	LinesInfo() LinesInfo
// 	// Functions returns all function define AST node(s) of file
// 	File() string
// 	// Lint()
// 	Lint(Outputer, Context)
// 	// AddDefect
// 	AddDefect(*Model, int, int, string, ...any)
// 	Defects() []*Defect
// 	ASTTree() *sitter.Tree
// 	FileRoot() Context
// 	Parent() Context
// 	SetParent(Context)
// 	AddChild(Context)
// 	RangeChildren(func(ctx Context))
// 	IsDir() bool
// 	// Functions
// 	Functions() []*Function
// 	// CallExpress
// 	CallExpresses() []*CallExpress
// 	//
// 	Classes() map[string]*Class
// 	Variables() []*Var
// 	Consts() []*Const
// }

// // BaseContext ...
// type BaseContext struct {
// 	// 文件的语言和检查规则
// 	*Linter
// 	// 文件路径
// 	file string
// 	// 子节点
// 	children []Context
// 	// 文件的内容
// 	content []byte
// 	// 文件的行信息
// 	info [][2]int
// 	// defects
// 	defects []*Defect
// 	//
// 	ast *sitter.Tree
// 	// parent
// 	parent Context
// 	// root file node
// 	root Context
// 	//
// 	funcs []*Function
// 	//
// 	classes map[string]*Class
// 	//
// 	vars   []*Var
// 	consts []*Const

// 	calls []*CallExpress
// 	sync.Once
// }

// // IsDir implements Context.
// func (b *BaseContext) IsDir() bool {
// 	return b.Linter == nil
// }

// // RangeChildren implements Context.
// func (b *BaseContext) RangeChildren(fn func(ctx Context)) {
// 	for index := range b.children {
// 		fn(b.children[index])
// 	}
// }

// // CreateBaseContext 为了测试需要
// func CreateTestBaseContext(
// 	lang utils.Language,
// 	file string,
// 	content []byte,
// ) *BaseContext {
// 	node := &BaseContext{
// 		file:    file,
// 		content: content,
// 		Linter:  &Linter{Lang: lang},
// 		classes: make(map[string]*Class),
// 	}
// 	return node
// }

// // ASTTree implements Context.
// func (b *BaseContext) ASTTree() *sitter.Tree {
// 	if b.ast == nil {
// 		parser := sitter.NewParser()
// 		parser.SetLanguage(b.Lang.Lang())
// 		tree, err := parser.ParseCtx(context.TODO(), nil, b.Content())
// 		if err != nil {
// 			log.Errorf("failed to parse file %q ast, err:%s", b.file, err)
// 		} else {
// 			b.ast = tree
// 		}
// 	}
// 	return b.ast
// }

// // FileTree implements Context.
// func (b *BaseContext) FileRoot() Context {
// 	if b.root == nil {
// 		b.root = b.getFileRoot()
// 	}
// 	return b.root
// }

// func (b *BaseContext) getFileRoot() Context {
// 	var tmp Context
// 	tmp = b
// 	for tmp.Parent() != nil {
// 		tmp = tmp.Parent()
// 	}
// 	return tmp
// }

// // Parent implements Context.
// func (b *BaseContext) Parent() Context {
// 	return b.parent
// }

// func (b *BaseContext) SetParent(parent Context) {
// 	b.parent = parent
// }

// // DefectSet implements Context.
// func (b *BaseContext) DefectSet() []*Defect {
// 	return b.defects
// }

// // Defect implements Context.
// func (b *BaseContext) Defect(model *Model, row int, col int, s string, args ...any) {
// 	def := &Defect{
// 		Model: model,
// 		Desc:  fmt.Sprintf(s, args...),
// 		Row:   row,
// 		Col:   col,
// 	}
// 	b.defects = append(b.defects, def)
// }

// // CallExpresses implements Context.
// func (b *BaseContext) CallExpresses() []*CallExpress {
// 	errors.Warning("baseContext no callexpresses")
// 	return nil
// }

// // Functions implements Context.
// func (b *BaseContext) Functions() []*Function {
// 	errors.Warning("baseContext no functions")
// 	return nil
// }

// func (b *BaseContext) Classes() map[string]*Class {
// 	errors.Warning("baseContext no classes")
// 	return nil
// }

// func (b *BaseContext) Variables() []*Var {
// 	errors.Warning("baseContext no Variables")
// 	return nil
// }

// func (b *BaseContext) Consts() []*Const {
// 	errors.Warning("baseContext no Consts")
// 	return nil
// }

// // Lint implements Context.
// func (b *BaseContext) Lint(output Outputer, ctx Context) {
// 	if b.Linter != nil && b.LintFunc != nil {
// 		b.LintFunc(output, ctx)
// 	}
// }

// // loadContent TODO
// func (b *BaseContext) loadContent() (err error) {
// 	b.content, err = os.ReadFile(b.file)
// 	return
// }

// // Name implements Context.
// func (b *BaseContext) File() string {
// 	return b.file
// }

// func (b *BaseContext) AddChild(node Context) {
// 	node.SetParent(b)
// 	b.children = append(b.children, node)
// }

// func (b *BaseContext) String() string {
// 	return fmt.Sprintf("<Node: %s>", b.File())
// }

// var _ Context = (*BaseContext)(nil)

// var _ Context = (*GolangContext)(nil)

// // BaseContext

// // Golang Context

// type GolangContext struct {
// 	BaseContext
// }

// // CallExpresses implements Context.
// // Subtle: this method shadows the method (*BaseContext).CallExpresses of GolangContext.BaseContext.
// func (g *GolangContext) CallExpresses() []*CallExpress {
// 	return g.calls
// }

// // Functions implements Context.
// // Subtle: this method shadows the method (*BaseContext).Functions of GolangContext.BaseContext.
// func (g *GolangContext) Functions() []*Function {
// 	g.parseElements()
// 	return g.funcs
// }

// func (g *GolangContext) Classes() map[string]*Class {
// 	g.parseElements()
// 	return g.classes
// }

// // Functions implements Context.
// // Subtle: this method shadows the method (*BaseContext).Functions of GolangContext.BaseContext.
// func (g *GolangContext) Variables() []*Var {
// 	g.parseElements()
// 	return g.vars
// }

// // Functions implements Context.
// // Subtle: this method shadows the method (*BaseContext).Functions of GolangContext.BaseContext.
// func (g *GolangContext) Consts() []*Const {
// 	g.parseElements()
// 	return g.consts
// }

// func (g *GolangContext) parseMethod(node *sitter.Node, content []byte) *Method {

// 	method := &Method{}
// 	if nameNode := node.ChildByFieldName("name"); nameNode != nil {
// 		method.name = nameNode.Content(content)
// 		method.node = node
// 	}
// 	if paramNode := node.ChildByFieldName("receiver"); paramNode.Type() == "parameter_list" {
// 		if decl := paramNode.Child(1); decl.Type() == "parameter_declaration" {
// 			if reciverNode := decl.ChildByFieldName("name"); reciverNode != nil {
// 				method.reciver = reciverNode.Content(content)
// 			}
// 			typeNode := decl.ChildByFieldName("type")
// 			if typeNode == nil {
// 				return nil
// 			}
// 			switch typeNode.Type() {
// 			case "type_identifier": // 绑定到了实例上
// 				method.instance = NewClass(typeNode.Content(content), nil)
// 			case "pointer_type": // 绑定到地址
// 				if !handlePointType(method, typeNode, content) {
// 					return nil
// 				}
// 			case "generic_type":
// 				if !handleGenericType(method, typeNode, content) {
// 					return nil
// 				}
// 			default:
// 				return nil
// 			}
// 		}
// 	}
// 	return method
// }

// func handleGenericType(method *Method, node *sitter.Node, content []byte) bool {
// 	if identifier := node.ChildByFieldName("type"); identifier != nil && identifier.Type() == "type_identifier" {
// 		method.instance = NewClass(identifier.Content(content), nil)
// 		return true
// 	}
// 	return false
// }

// func handlePointType(method *Method, node *sitter.Node, content []byte) bool {
// 	method.prt = true
// 	first := node.Child(1)
// 	if first == nil {
// 		return false
// 	}
// 	switch first.Type() {
// 	case "generic_type":
// 		return handleGenericType(method, first, content)
// 	case "type_identifier":
// 		method.instance = NewClass(first.Content(content), nil)
// 	default:
// 		return false
// 	}
// 	return true
// }

// func (g *GolangContext) parseElements() {

// 	g.Do(func() {
// 		content := g.Content()
// 		ast.ApplyChildrenNodes(g.ASTTree().RootNode(), func(sub *sitter.Node) {
// 			switch sub.Type() {
// 			case "method_declaration":
// 				method := g.parseMethod(sub, content)
// 				if method == nil {
// 					utils.Bugf("cannot parse method: %s", sub.Content(content))
// 				} else {
// 					if class, ok := g.classes[method.instance.name]; !ok {
// 						g.classes[method.instance.name] = NewClass(method.instance.name, nil)
// 					} else {
// 						class.AddMethod(method)
// 					}
// 				}

// 			case "const_declaration":
// 				ast.ApplyChildrenNodes(sub, func(spec *sitter.Node) {
// 					if spec.Type() == "const_spec" {
// 						nameNode := spec.ChildByFieldName("name")
// 						if nameNode != nil {
// 							g.consts = append(g.consts, &Const{BaseElem{name: nameNode.Content(content), node: spec}})
// 						} else {
// 							utils.Bugf("cannot parse const desclare: %s", sub.Content(content))
// 						}
// 					}
// 				})

// 			case "var_declaration":
// 				ast.ApplyChildrenNodes(sub, func(spec *sitter.Node) {
// 					if spec.Type() == "var_spec" {
// 						nameNode := spec.ChildByFieldName("name")
// 						if nameNode != nil {
// 							g.consts = append(g.consts, &Const{BaseElem{name: nameNode.Content(content), node: spec}})
// 						} else {
// 							utils.Bugf("cannot parse var desclare: %s", sub.Content(content))
// 						}
// 					}
// 				})
// 			case "function_declaration":
// 				if nameNode := sub.ChildByFieldName("name"); nameNode.Type() == "identifier" {
// 					g.funcs = append(g.funcs, &Function{BaseElem{name: nameNode.Content(content), node: sub}})
// 				} else {
// 					utils.Bugf("cannot parse function desclare: %s", sub.Content(content))
// 				}
// 			case "type_declaration":
// 				if spec := sub.Child(1); spec.Type() == "type_spec" {
// 					if nameNode := spec.ChildByFieldName("name"); nameNode.Type() == "type_identifier" {
// 						name := nameNode.Content(content)
// 						if class, ok := g.classes[name]; !ok {
// 							g.classes[name] = NewClass(name, sub)
// 						} else {
// 							class.node = sub
// 						}
// 					} else {
// 						utils.Bugf("cannot parse struct declare: %s", sub.Content(content))
// 					}
// 				}
// 			}
// 		})
// 	})

// }

// func NewGlangContext(file string, linter *Linter) Context {
// 	baseContext := BaseContext{
// 		file:    file,
// 		Linter:  linter,
// 		classes: make(map[string]*Class),
// 	}
// 	return &GolangContext{
// 		BaseContext: baseContext,
// 	}
// }

// // NewPythonContext create a context of python script file
// func NewPythonContext(file string, linter *Linter) Context {
// 	return nil
// }

// func NewCCppContext(file string, linter *Linter) Context {
// 	return nil
// }

// func CreateContext(file string, linter *Linter) Context {
// 	if linter == nil {
// 		return &BaseContext{
// 			file:    file,
// 			classes: make(map[string]*Class),
// 		}
// 	}
// 	switch linter.Lang {
// 	case utils.GoLang:
// 		return NewGlangContext(file, linter)
// 	case utils.Python:
// 		return NewPythonContext(file, linter)
// 	case utils.CCpp:
// 		return NewCCppContext(file, linter)
// 	}
// 	return nil
// }

type Context interface {
	// file
	File() string
	IsDir() bool
	Lines() LinesInfo
	Content() []byte
	AddSubContext(Context)
	Range(func(ctx Context))
	Check()
	// language
	Package() string
	Functions() map[string]*Function
	Classes() map[string]*Class
	Consts() map[string]*Const
	Varibales() map[string]*Variable
	SetErrHandler(func(error))
}

type FileContext struct {
	file      string
	linesInfo LinesInfo
	filetype  utils.Language
	dir       bool
	children  []Context
	content   []byte
	check     ModelFuncType
	// 文件所属的 package name
	pkg        string
	functions  map[string]*Function
	classes    map[string]*Class
	varibales  map[string]*Variable
	consts     map[string]*Const
	errHandler func(err error)
}

// Check implements Context.
func (f *FileContext) Check() {
	if f.check != nil {
		f.check(f)
	}
}

type emptyCtxType struct {
	FileContext
	sync.Mutex
}

var emptyCtx = &emptyCtxType{}

func (e *emptyCtxType) acquire() {
	e.Lock()
}
func (e *emptyCtxType) release() {
	e.children = e.children[:1]
	e.Unlock()
}

func (f *FileContext) File() string {
	return f.file
}

func (f *FileContext) IsDir() bool {
	return f.dir
}

// Lines implements Context.
//
//	0  -   \r
//	1  -   \n
//	2  -   \r\n
func (f *FileContext) Lines() LinesInfo {
	if f.linesInfo == nil {
		ctt := f.Content()
		gap := 0
		index, length := 0, len(ctt)
		for index < length {
			switch ctt[index] {
			case '\r':
				lineLength := len(tool.ToString(ctt[gap:index]))
				if index+1 < length {
					if ctt[index+1] == '\n' {
						// \r\n
						f.linesInfo = append(f.linesInfo, [2]int{lineLength, 2})
						index += 1
						gap = index + 1
					} else {
						// \r
						f.linesInfo = append(f.linesInfo, [2]int{lineLength, 0})
						gap = index + 1
					}
				} else {
					// EOF
					f.linesInfo = append(f.linesInfo, [2]int{lineLength, 0})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				f.linesInfo = append(f.linesInfo, [2]int{lineLength, 1})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			f.linesInfo = append(f.linesInfo, [2]int{lineLength, -1})
		}
	}
	return f.linesInfo
}

func (f *FileContext) Content() []byte {
	var err error
	if f.content == nil {
		f.content, err = os.ReadFile(f.file)
		f.errHandler(err)
	}
	return f.content
}

func (f *FileContext) AddSubContext(ctx Context) {
	f.children = append(f.children, ctx)
}

func (f *FileContext) Range(fn func(ctx Context)) {
	for index := range f.children {
		fn(f.children[index])
	}
}

// Package 在pre中创建
func (f *FileContext) Package() string {
	return f.pkg
}

// Classes implements BaseContext.
func (f *FileContext) Classes() map[string]*Class {
	return f.classes
}

// Consts implements BaseContext.
func (f *FileContext) Consts() map[string]*Const {
	return f.consts
}

// Functions implements BaseContext.
func (f *FileContext) Functions() map[string]*Function {
	return f.functions
}

// Varibales implements BaseContext.
func (f *FileContext) Varibales() map[string]*Variable {
	return f.varibales
}

func (f *FileContext) SetErrHandler(fn func(err error)) {
	f.errHandler = fn
}

var _ Context = (*FileContext)(nil)

type ContextMaker struct {
	cxtCreaterMap map[string]MakeContextFunc
}

type MakeContextFunc func(string) Context

type PreHandlerFunc func(*config.Config, Context) error

func NewContextMaker(languages []*config.Language) (*ContextMaker, error) {
	return &ContextMaker{}, nil
}

func (c *ContextMaker) New(path string) Context {
	return nil
}
