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
	"github.com/stkali/glint/utils"
)

// DispatchLang 返回对应的语言
func DispatchLang(path string) utils.Language {
	ext := filepath.Ext(path)
	switch ext {
	case ".c", ".h", ".cpp", ".hpp":
		return utils.CCpp
	case ".rs":
		return utils.Rust
	case ".go":
		return utils.Go
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
	utils.CCpp: cpp.GetLanguage(),
	utils.Rust: rust.GetLanguage(),
	utils.Go: golang.GetLanguage(),
	utils.Python: python.GetLanguage(),
	utils.Java:java.GetLanguage(),
	utils.JavaScript: javascript.GetLanguage(),
}
