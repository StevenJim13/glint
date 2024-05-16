// ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
package glint

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
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
	Functions() []Function
	// CallExpress
	CallExpresses() []CallExpress
	// AddDefect
	Defect(model *Model, row, col int, s string, args ...any)
	DefectSet() []*Defect
	ASTTree() *sitter.Tree
	FileRoot() *FileNode
	Parent() *FileNode
}

type Defect struct {
	Model *Model
	Desc  string
	Row   int
	Col   int
}

func (d *Defect) String() string {
	return fmt.Sprintf("model: %q, desc: %s, position:(%d,%d)", *&d.Model.Name, d.Desc, d.Row, d.Col)
}

type Function struct {
	Name     string
	Return   string
	Position [2]int
}

type CallExpress struct {
	Function *Function
}

// FileNode
type FileNode struct {
	// 文件的语言和检查规则
	*Linter
	// 文件路径
	file string
	// 子节点
	Children []*FileNode
	// 文件的内容
	content []byte
	// 文件的行信息
	info [][2]int
	// defects
	defects []*Defect
	//
	ast *sitter.Tree
	// parent
	parent *FileNode
	// root file node
	root *FileNode
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
func (f *FileNode) FileRoot() *FileNode {
	if f.root == nil {
		f.root = f.getFileRoot()
	}
	return f.root
}

func (f *FileNode) getFileRoot() *FileNode {
	tmp := f
	for tmp.parent != nil {
		tmp = tmp.parent
	}
	return tmp
}

// Parent implements Context.
func (f *FileNode) Parent() *FileNode {
	return f.parent
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
func (f *FileNode) CallExpresses() []CallExpress {
	return nil
}

// Functions implements Context.
func (f *FileNode) Functions() []Function {
	return nil
}

// Lint implements Context.
func (f *FileNode) Lint(output Outputer, ctx Context) {
	if f.LintFunc != nil {
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
						f.info = append(f.info, [2]int{lineLength, 3})
						index += 1
						gap = index + 1
					} else {
						// \r
						f.info = append(f.info, [2]int{lineLength, 1})
						gap = index + 1
					}
				} else {
					// EOF
					f.info = append(f.info, [2]int{lineLength, 1})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				f.info = append(f.info, [2]int{lineLength, 2})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			f.info = append(f.info, [2]int{lineLength, 2})
		}
	}
	return f.info
}

// Name implements Context.
func (f *FileNode) File() string {
	return f.file
}

func (f *FileNode) AddChild(node *FileNode) {
	node.parent = f
	f.Children = append(f.Children, node)
}

func (f *FileNode) String() string {
	return fmt.Sprintf("<Node: %s>", f.File())
}

var _ Context = (*FileNode)(nil)

type LinesInfo [][2]int

func (l LinesInfo) String() string {
	return fmt.Sprintf("<LinesInfo(%d)>", len(l))
}

func (l LinesInfo) Lines() int {
	return len(l)
}

func (l LinesInfo) Range(f func(index int, line [2]int) bool) {
	for index := range l {
		if !f(index, l[index]) {
			return
		}
	}
}

var langauges = map[utils.Language]*sitter.Language{
	utils.CCpp:       cpp.GetLanguage(),
	utils.Rust:       rust.GetLanguage(),
	utils.GoLang:     golang.GetLanguage(),
	utils.Python:     python.GetLanguage(),
	utils.Java:       java.GetLanguage(),
	utils.JavaScript: javascript.GetLanguage(),
}

type FileTree struct {
	Root     string
	rootNode *FileNode
}

func NewFileTree(root string) *FileTree {
	tree := &FileTree{
		Root: tool.ToAbsPath(root),
	}
	return tree
}

func (f *FileTree) RootNode() *FileNode {
	return f.rootNode
}

func (f *FileTree) AddChild(node *FileNode) {
	if f.rootNode == nil {
		f.rootNode = node
		return
	}
	panic("Tree head node has been set")
}

func (f *FileTree) Parse(excFiles, excDirs []string, dispatch func(string) *Linter) error {
	exclude, err := getExclude(excFiles, excDirs)
	if err != nil {
		return err
	}
	err = buildFileTree(f.Root, f, exclude, dispatch)
	if err != nil {
		return err
	}
	log.Infof("successfully to build file tree: %s", f)
	return nil
}

func (f *FileTree) Walk(fn func(node *FileNode) error) {
	walk(f.rootNode, fn)
}

func (f *FileTree) String() string {
	return fmt.Sprintf("<FileTree: %s>", f.Root)
}

func walk(node *FileNode, fn func(node *FileNode) error) {
	fn(node)
	if len(node.Children) != 0 {
		for index := range node.Children {
			walk(node.Children[index], fn)
		}
	}
}

func buildFileTree(
	path string,
	root interface{ AddChild(*FileNode) },
	exclude func(string, bool) bool,
	dispatch func(string) *Linter,
) error {

	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		if exclude(info.Name(), false) {
			return nil
		} else {
			node := &FileNode{
				file: path,
			}
			root.AddChild(node)
			dirs, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for index := range dirs {
				subPath := filepath.Join(path, dirs[index].Name())
				if err := buildFileTree(subPath, node, exclude, dispatch); err != nil {
					return err
				}
			}
		}
	} else {
		// file
		// 是否需要被排除
		if exclude(info.Name(), true) {
			return nil
		} else {
			linter := dispatch(info.Name())
			if linter == nil {
				log.Errorf("exclude: %s", info.Name())
				return nil
			}
			node := &FileNode{
				file:   path,
				Linter: linter,
			}
			root.AddChild(node)
		}
	}
	return nil
}

func getExclude(excFiles, excDirs []string) (func(path string, file bool) bool, error) {
	veriryFile, err := makeExcludeFunc(excFiles...)
	if err != nil {
		return nil, err
	}
	verifyDir, err := makeExcludeFunc(excDirs...)
	if err != nil {
		return nil, err
	}
	if verifyDir == nil && veriryFile == nil {
		return func(path string, file bool) bool { return false }, nil
	}
	return func(path string, file bool) bool {
		if file && veriryFile != nil {
			return veriryFile(path)
		} else if verifyDir != nil {
			return verifyDir(path)
		}
		return false
	}, nil
}

type VerifyFileFunc func(path string) bool

func makeExcludeFunc(excludes ...string) (VerifyFileFunc, error) {
	var verify VerifyFileFunc
	switch len(excludes) {
	case 0:
	case 1:
		verify = func(path string) bool {
			if ok, err := filepath.Match(excludes[0], path); err != nil {
				panic(err)
			} else {
				return ok
			}
		}

	default:
		verify = func(path string) bool {
			for index := range excludes {
				if ok, err := filepath.Match(excludes[index], path); err != nil {
					panic(err)
				} else if ok {
					return true
				}
			}
			return false
		}
	}
	return verify, nil
}
