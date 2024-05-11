package utils

import (
	"strings"

	"github.com/stkali/utility/errors"
)

type Language int

func(l Language) String () string {
	if l > Unknown || l < CCpp {
		return langLabels[Unknown]
	}
	return langLabels[l]
}

const (
	Unknown Language = iota
	CCpp
	Rust
	GoLang
	Python
	Java
	JavaScript
	CSharp
	Ruby
	Perl
	Shell
	Markdown
	Dockerfile
	YAML 
	Makefile
	maxLang
)

var langLabels = []string{
	"C/C++",
	"Rust",
	"Golang",
	"Python",
	"Java",
	"JavaScript",
	"CSharp",
	"Ruby",
	"Perl",
	"Shell",
	"Markdown",
	"Dockerfile",
	"YAML",
	"Makefile",
}

var extendsTable = [][]string{
	{".c", ".h", ".cpp", ".hpp", ".cxx"},
	{".rs"},
	{".go"},
	{".py"},
	{".java"},
	{".js"},
	{".cs"},
	{".rb"},
	{".pl"},
	{".sh"},
	{".md"},
	{""},
	{".yaml"},
	{""},
}


func Extends(lang Language) ([]string, error) {
	if lang < CCpp || lang >= maxLang {
		return nil, errors.Newf("invalid language: %s", lang)
	}
	return extendsTable[lang], nil
}

func ToLanguage(name string) Language {
	
	switch strings.ToLower(name) {
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
	case "js", "javascript":
		return JavaScript
	case "c#", "csharp":
		return CSharp
	case "ruby", "rb":
		return Ruby
	case "perl", "pl":
		return Perl
	}
	return Unknown
}