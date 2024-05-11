package parser

import (
	"fmt"
	"io/fs"
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
}
type Defect interface {
}
type Function struct {
	Name   string
	Return string
}

type CallExpress struct {
	Function *Function
}

type LinesInfo [][2]uint

func (l LinesInfo) String() string {
	return fmt.Sprintf("<LinesInfo(%d)>", len(l))
}

func (l LinesInfo) Lines() int {
	return len(l)
}

func (l LinesInfo) Range(f func(line [2]uint) bool) {
	for index := range l {
		if !f(l[index]) {
			return
		}
	}
}

// DispatchLang 返回对应的语言
func DispatchLang(path string) utils.Language {
	ext := filepath.Ext(path)
	switch ext {
	case ".c", ".h", ".cpp", ".hpp":
		return utils.CCpp
	case ".rs":
		return utils.Rust
	case ".go":
		return utils.GoLang
	case ".py":
		return utils.Python
	case ".java":
		return utils.Java
	case ".js":
		return utils.JavaScript
	case ".cs":
		return utils.CSharp
	case ".rb":
		return utils.Ruby
	case ".pl":
		return utils.Perl
	}
	return utils.Unknown
}

var langauges = map[utils.Language]*sitter.Language{
	utils.CCpp:       cpp.GetLanguage(),
	utils.Rust:       rust.GetLanguage(),
	utils.GoLang:     golang.GetLanguage(),
	utils.Python:     python.GetLanguage(),
	utils.Java:       java.GetLanguage(),
	utils.JavaScript: javascript.GetLanguage(),
}

type LangParser struct {
}

type Parser interface {
	Parse(root string, ch Context) error
}

func NewParser() (*LangParser, error) {

	return &LangParser{}, nil
}

func (l *LangParser) Parse(root string, ch Context) error {

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {

		log.Infof("path: %s, d: %s, err: %s", path, d.Name(), err)
		return nil
	})
	return nil
}
