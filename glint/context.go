package glint

import (
	"context"
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
)

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
	// Functions
	Functions() []*Function
	// CallExpress
	CallExpresses() []*CallExpress
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
}

// FileNode ...
type FileNode struct {
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
	calls []*CallExpress
}

// IsDir implements Context.
func (f *FileNode) IsDir() bool {
	return f.Linter == nil
}

// RangeChildren implements Context.
func (f *FileNode) RangeChildren(fn func(ctx Context)) {
	for index := range f.children {
		fn(f.children[index])
	}
}

type Function struct {
	name string
	ret  string
	node *sitter.Node
}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Row() int {
	return int(f.node.StartPoint().Row)
}

func (f *Function) Col() int {
	return int(f.node.StartPoint().Column)
}

func (f *Function) Node() *sitter.Node {
	return f.node
}

type CallExpress struct {
	Function *Function
}

// CreateFileNode 为了测试需要
func CreateTestFileNode(
	lang utils.Language,
	file string,
	content []byte,
) *FileNode {
	node := &FileNode{
		file:    file,
		content: content,
		Linter:  &Linter{Lang: lang},
	}
	return node
}

// ASTTree implements Context.
func (f *FileNode) ASTTree() *sitter.Tree {
	if f.ast == nil {
		parser := sitter.NewParser()
		parser.SetLanguage(f.Lang.Lang())
		tree, err := parser.ParseCtx(context.TODO(), nil, f.Content())
		if err != nil {
			log.Errorf("failed to parse file %q ast, err:%s", f.file, err)
		} else {
			f.ast = tree
		}
	}
	return f.ast
}

// FileTree implements Context.
func (f *FileNode) FileRoot() Context {
	if f.root == nil {
		f.root = f.getFileRoot()
	}
	return f.root
}

func (f *FileNode) getFileRoot() Context {
	var tmp Context
	tmp = f
	for tmp.Parent() != nil {
		tmp = tmp.Parent()
	}
	return tmp
}

// Parent implements Context.
func (f *FileNode) Parent() Context {
	return f.parent
}

func (f *FileNode) SetParent(parent Context) {
	f.parent = parent
}

// DefectSet implements Context.
func (f *FileNode) DefectSet() []*Defect {
	return f.defects
}

// Defect implements Context.
func (f *FileNode) Defect(model *Model, row int, col int, s string, args ...any) {
	def := &Defect{
		Model: model,
		Desc:  fmt.Sprintf(s, args...),
		Row:   row,
		Col:   col,
	}
	f.defects = append(f.defects, def)
}

// CallExpresses implements Context.
func (f *FileNode) CallExpresses() []*CallExpress {
	errors.Warning("baseContext no callexpresses")
	return nil
}

// Functions implements Context.
func (f *FileNode) Functions() []*Function {
	errors.Warning("baseContext no functions")
	return nil
}

// Lint implements Context.
func (f *FileNode) Lint(output Outputer, ctx Context) {
	if f.Linter != nil && f.LintFunc != nil {
		f.LintFunc(output, ctx)
	}
}

// Content implements Context.
func (f *FileNode) Content() []byte {
	if f.content == nil {
		if err := f.loadContent(); err != nil {
			errors.Warningf("failed to get file: %q content, err: %s", f.file, err)
		}
	}
	return f.content
}

// loadContent TODO
func (f *FileNode) loadContent() (err error) {
	f.content, err = os.ReadFile(f.file)
	return
}

// LinesInfo implements Context.
//
//	0  -   \r
//	1  -   \n
//	2  -   \r\n
func (f *FileNode) LinesInfo() LinesInfo {
	if f.info == nil {
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
						f.info = append(f.info, [2]int{lineLength, 2})
						index += 1
						gap = index + 1
					} else {
						// \r
						f.info = append(f.info, [2]int{lineLength, 0})
						gap = index + 1
					}
				} else {
					// EOF
					f.info = append(f.info, [2]int{lineLength, 0})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				f.info = append(f.info, [2]int{lineLength, 1})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			f.info = append(f.info, [2]int{lineLength, -1})
		}
	}
	return f.info
}

// Name implements Context.
func (f *FileNode) File() string {
	return f.file
}

func (f *FileNode) AddChild(node Context) {
	node.SetParent(f)
	f.children = append(f.children, node)
}

func (f *FileNode) String() string {
	return fmt.Sprintf("<Node: %s>", f.File())
}

var _ Context = (*FileNode)(nil)

var _ Context = (*GolangContext)(nil)

// BaseContext

// Golang Context

type GolangContext struct {
	FileNode
}

// CallExpresses implements Context.
// Subtle: this method shadows the method (*FileNode).CallExpresses of GolangContext.FileNode.
func (g *GolangContext) CallExpresses() []*CallExpress {

	return g.calls
}

// Functions implements Context.
// Subtle: this method shadows the method (*FileNode).Functions of GolangContext.FileNode.
func (g *GolangContext) Functions() []*Function {
	return g.funcs
}

func NewGlangContext(file string, linter *Linter) Context {
	baseContext := FileNode{
		file:   file,
		Linter: linter,
	}
	return &GolangContext{
		FileNode: baseContext,
	}
}

func NewPythonContext(file string, linter *Linter) Context {
	return nil
}

func NewCCppContext(file string, linter *Linter) Context {
	return nil
}

func CreateContext(file string, linter *Linter) Context {
	if linter == nil {
		return &FileNode{file: file}
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
