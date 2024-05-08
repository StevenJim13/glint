package util

import "github.com/stkali/utility/errors"

type Language int

func(l Language) String () string {
	if l > Unknown || l < CCpp {
		return langLabels[Unknown]
	}
	return langLabels[l]
}

const (
	CCpp Language = iota
	Rust
	Go 
	Python
	Java
	JavaScript
	CSharp
	Ruby
	Perl
	Unknown
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
}


func Extends(lang Language) ([]string, error) {
	if lang < CCpp || lang >= Unknown {
		return nil, errors.Newf("invalid language: %s", lang)
	}
	return extendsTable[lang], nil
}

