package parser

import (
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/stkali/glint/util"
)

// DispatchLang 返回对应的语言
func DispatchLang(path string) util.Language {
	ext := filepath.Ext(path)
	switch ext {
	case ".c", ".h", ".cpp", ".hpp":
		return util.CCpp
	case ".rs":
		return util.Rust
	case ".go":
		return util.Go
	case ".py":
		return util.Python
	case ".java":
		return util.Java
	case ".js":
		return util.JavaScript
	case ".cs":
		return util.CSharp
	case ".rb":
		return util.Ruby
	case ".pl":
		return util.Perl
	}
	return util.Unknown
}

var langauges = map[util.Language]*sitter.Language{
	util.CCpp: cpp.GetLanguage(),
	util.Rust: rust.GetLanguage(),
	util.Go: golang.GetLanguage(),
	util.Python: python.GetLanguage(),
	util.Java:java.GetLanguage(),
	util.JavaScript: javascript.GetLanguage(),
}
