package parser

import (
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
	// AddDefect adds a defect to
	AddDefect(defect Defect)
	// LinesInfo returns line information of file
	LinesInfo() LinesInfo
	// Functions returns all function define AST node(s) of file
	Functions() []*Function
	// CallExpresses returns all callexpression node(s) of file
	CallExpresses() []*CallExpress
	//
	File() string
	// Lint()
	Lint(ctx Context)
}

type Defect interface {
	String() string
}

type Function struct {
	Name   string
	Return string
}

type CallExpress struct {
	Function *Function
}

type LinesInfo [][2]int

func (l LinesInfo) String() string {
	return fmt.Sprintf("<LinesInfo(%d)>", len(l))
}

func (l LinesInfo) Lines() int {
	return len(l)
}

func (l LinesInfo) Range(f func(line [2]int) bool) {
	for index := range l {
		if !f(l[index]) {
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

type Linter struct {
	Lang     utils.Language
	LintFunc LintModels
}

type FileNode struct {
	// 所属语言
	Language utils.Language
	// 文件路径
	file string
	// 子节点
	Children []*FileNode
	content  []byte
	info     [][2]int
	*Linter
}

// Lint implements Context.
func (n *FileNode) Lint(ctx Context) {
	if n.LintFunc != nil {
		n.LintFunc(ctx)
	}
}

// AddDefect implements Context.
func (n *FileNode) AddDefect(defect Defect) {
	return
}

// CallExpresses implements Context.
func (n *FileNode) CallExpresses() []*CallExpress {
	return nil
}

// Content implements Context.
func (n *FileNode) Content() []byte {
	if n.content == nil {
		if err := n.loadContent(); err != nil {
			errors.Warningf("failed to get file: %q content, err: %s", n.file, err)
		}
	}
	return n.content
}

// loadContent TODO
func (n *FileNode) loadContent() (err error) {
	n.content, err = os.ReadFile(n.file)
	return
}

// Functions implements Context.
func (n *FileNode) Functions() []*Function {
	return nil
}

// LinesInfo implements Context.
func (n *FileNode) LinesInfo() LinesInfo {
	if n.info == nil {
		ctt := n.Content()
		gap := 0
		index, length := 0, len(ctt)
		for index < length {
			switch ctt[index] {
			case '\r':
				lineLength := len(tool.ToString(ctt[gap:index]))
				if index+1 < length {
					if ctt[index+1] == '\n' {
						// \r\n
						n.info = append(n.info, [2]int{lineLength, 3})
						index += 1
						gap = index + 1
					} else {
						// \r
						n.info = append(n.info, [2]int{lineLength, 1})
						gap = index + 1
					}
				} else {
					// EOF
					n.info = append(n.info, [2]int{lineLength, 1})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				n.info = append(n.info, [2]int{lineLength, 2})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			n.info = append(n.info, [2]int{lineLength, 2})
		}
	}
	return n.info
}

// Name implements Context.
func (n *FileNode) File() string {
	return n.file
}

func (n *FileNode) AddChild(node *FileNode) {
	n.Children = append(n.Children, node)
}

func (n *FileNode) String() string {
	return fmt.Sprintf("<Node: %s>", n.File)
}

var _ Context = (*FileNode)(nil)

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

func (f *FileTree) Parse(excFiles, excDirs []string, matcher Matcher) error {
	exclude, err := getExclude(excFiles, excDirs)
	if err != nil {
		return err
	}
	err = buildFileTree(f.Root, f, exclude, matcher)
	if err != nil {
		return err
	}
	log.Infof("successfully to build file tree: %s", f)
	return nil
}

func (f *FileTree) Walk(fn func(node *FileNode) error) {
	walk(f.rootNode, fn)
}

func walk(node *FileNode, fn func(node *FileNode) error) {
	fn(node)
	if len(node.Children) != 0 {
		for index := range node.Children {
			walk(node.Children[index], fn)
		}
	}
}

func (f *FileTree) String() string {
	return fmt.Sprintf("<FileTree: %s>", f.Root)
}

func buildFileTree(
	path string,
	root interface{ AddChild(*FileNode) },
	exclude func(path string, file bool) bool,
	matcher Matcher,
) error {

	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		//
		if exclude(info.Name(), false) {
			log.Infof("exclude: %s", info.Name())
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
				if err := buildFileTree(subPath, node, exclude, matcher); err != nil {
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
			linter := matcher.Match(info.Name())
			if linter == nil {
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

type LintModels func(ctx Context)

type Matcher interface {
	Match(file string) *Linter
}
