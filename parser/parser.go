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
	Name() string
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

func NewContext(path string) Context {
	return &ASTContext{
		file: path,
	}
}

type FileNode struct {
	// 所属语言
	Language utils.Language
	// 文件路径
	File string
	// 子节点
	Children []*FileNode
	// LintFunc
	Linter *Linter
}

func (n *FileNode) AddChild(node *FileNode) {
	n.Children = append(n.Children, node)
}

func (n *FileNode) String() string {
	return fmt.Sprintf("<Node: %s>", n.File)
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
				File: path,
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
				File:   path,
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

type Linter struct {
	Lang     utils.Language
	LintFunc LintModels
}

func (l *Linter) String() string {
	return fmt.Sprintf("<Linter: %s>", l.Lang)
}
