package glint

import (
	"fmt"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/tool"
)

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
	rootNode Context
}

func NewFileTree(root string) *FileTree {
	tree := &FileTree{
		Root: tool.ToAbsPath(root),
	}
	return tree
}

func (f *FileTree) RootNode() Context {
	return f.rootNode
}

func (f *FileTree) AddSubContext(node Context) {
	if f.rootNode == nil {
		f.rootNode = node
		return
	}
	panic("Tree head node has been set")
}

// func (f *FileTree) Parse(excFiles, excDirs []string, maker *ContextMaker) error {
// 	exclude, err := getExclude(excFiles, excDirs)
// 	if err != nil {
// 		return err
// 	}
// 	err = buildFileTree(f.Root, f, exclude, maker)
// 	if err != nil {
// 		return err
// 	}
// 	log.Infof("successfully to build file tree: %s", f)
// 	return nil
// }

// func (f *FileTree) Walk(fn func(ctx Context) error) {
// 	walk(f.rootNode, fn)
// }

// func (f *FileTree) String() string {
// 	return fmt.Sprintf("<FileTree: %s>", f.Root)
// }

// func walk(ctx Context, fn func(ctx Context) error) {
// 	fn(ctx)
// 	ctx.Range(func(ctx Context) {
// 		walk(ctx, fn)
// 	})
// }

// func buildFileTree(
// 	path string,
// 	root interface{ AddSubContext(Context) },
// 	exclude func(string, bool) bool,
// 	maker *ContextMaker,
// ) error {

// 	info, err := os.Lstat(path)
// 	if err != nil {
// 		return err
// 	}
// 	if info.IsDir() {
// 		if exclude(info.Name(), false) {
// 			return nil
// 		} else {
// 			ctx := maker.New(path)
// 			root.AddSubContext(ctx)
// 			dirs, err := os.ReadDir(path)
// 			if err != nil {
// 				return err
// 			}
// 			for index := range dirs {
// 				subPath := filepath.Join(path, dirs[index].Name())
// 				if err := buildFileTree(subPath, ctx, exclude, maker); err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	} else {
// 		// file
// 		// 是否需要被排除
// 		if exclude(info.Name(), true) {
// 			return nil
// 		} else {
// 			ctx := maker.New(info.Name())
// 			root.AddSubContext(ctx)
// 		}
// 	}
// 	return nil
// }

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
