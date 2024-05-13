package utils

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/bash"
	"github.com/smacker/go-tree-sitter/cpp"
	"github.com/smacker/go-tree-sitter/csharp"
	"github.com/smacker/go-tree-sitter/css"
	"github.com/smacker/go-tree-sitter/dockerfile"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/html"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/ruby"
	"github.com/smacker/go-tree-sitter/rust"
	"github.com/smacker/go-tree-sitter/yaml"
	"github.com/stkali/utility/errors"
)

type Language int

const (
	Unknown Language = iota
	Any
	CCpp
	Rust
	GoLang
	Python
	Java
	CSharp
	JavaScript
	HTML
	CSS
	Ruby
	Bash
	Dockerfile
	YAML
	maxLang
)

var labelTable = []string{
	"Unknown",
	"Any",
	"C/C++",
	"Rust",
	"Golang",
	"Python",
	"Java",
	"Csharp",
	"JavaScript",
	"HTML",
	"CSS",
	"Ruby",
	"Bash",
	"Dockerfile",
	"YAML",
}

func (l Language) String() string {
	if l >= maxLang || l <= Unknown {
		return fmt.Sprintf("Unknown(%d)", l)
	}
	return labelTable[l]
}

var extendsTable = [][]string{
	{""},
	{""},
	{".c", ".h", ".cpp", ".hpp", ".cxx"},
	{".rs"},
	{".go"},
	{".py"},
	{".java"},
	{".cs"},
	{".js"},
	{".html"},
	{".css"},
	{".rb"},
	{".sh"},
	{""},
	{".yaml"},
}

func Extends(lang Language) ([]string, error) {
	if lang >= maxLang || lang <= Unknown {
		return nil, errors.Newf("invalid language: %s", lang)
	}
	return extendsTable[lang], nil
}

func ToLanguage(name string) Language {

	switch strings.ToLower(name) {
	case "*":
		return Any
	case "c", "c++", "c/c++", "cpp":
		return CCpp
	case "rust", "rt":
		return Rust
	case "go", "golang":
		return GoLang
	case "py", "python", "python3":
		return Python
	case "java":
		return Java
	case "c#", "csharp":
		return CSharp
	case "js", "javascript":
		return JavaScript
	case "html":
		return HTML
	case "css":
		return CSS
	case "ruby", "rb":
		return Ruby
	case "bash", "shell", "sh":
		return Bash
	case "dockerfile", "docker", "df":
		return Dockerfile
	case "yaml", "yml":
		return YAML
	}
	return Unknown
}

var langTable = []*sitter.Language{
	nil,
	nil,
	cpp.GetLanguage(),
	rust.GetLanguage(),
	golang.GetLanguage(),
	python.GetLanguage(),
	java.GetLanguage(),
	csharp.GetLanguage(),
	javascript.GetLanguage(),
	html.GetLanguage(),
	css.GetLanguage(),
	ruby.GetLanguage(),
	bash.GetLanguage(),
	dockerfile.GetLanguage(),
	yaml.GetLanguage(),
}

func (l Language) Lang() *sitter.Language {
	if l >= maxLang || l <= Unknown {
		return nil
	}
	return langTable[l]
}
