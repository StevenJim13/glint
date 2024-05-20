/*
glint
j
*/
package glint

import (
	"fmt"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
)

type Packager interface {
	Package() string
	SetPkgName(name string)
	AddSubContext(Context)
	Range(func(ctx Context))
}

type Sourcer interface {
	AST() *sitter.Tree
	Functions() map[string]*Function
	Classes() map[string]*Class
	CallExpresses() map[string]*CallExpress
	Consts() map[string]*Const
	Varibales() map[string]*Variable
}

type Filer interface {
	Path() string
	Lines() LinesInfo
	Content() []byte
	// Check() error
	Defects() []*Defect
	AddDefect(*Defect)
}

type Context interface {
	Packager
	Filer
	Sourcer
	Check() error
	HandleErr(error)
	IsPackage() bool
}

type Pkg struct {
	children []Context
	pkg      string
}

// SetName implements Packager.
func (p *Pkg) SetPkgName(name string) {
	p.pkg = name
}

func NewPackage() *Pkg {
	return &Pkg{}
}

// AddSubContext implements Packager.
func (p *Pkg) AddSubContext(ctx Context) {
	p.children = append(p.children, ctx)
}

// Package implements Packager.
func (p *Pkg) Package() string {
	return p.pkg
}

// Range implements Packager.
func (p *Pkg) Range(fn func(ctx Context)) {
	for index := range p.children {
		fn(p.children[index])
	}
}

var _ Packager = (*Pkg)(nil)

// -------------------------------------
type File struct {
	path      string
	linesInfo LinesInfo
	filetype  utils.Language
	dir       bool
	content   []byte
	defects   []*Defect
}

func NewFile(path string) *File {
	return &File{path: path}
}

func (f *File) Path() string {
	return f.path
}

// Lines implements Context.
//
//	0  -   \r
//	1  -   \n
//	2  -   \r\n
func (f *File) Lines() LinesInfo {
	if f.linesInfo == nil {
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
						f.linesInfo = append(f.linesInfo, [2]int{lineLength, 2})
						index += 1
						gap = index + 1
					} else {
						// \r
						f.linesInfo = append(f.linesInfo, [2]int{lineLength, 0})
						gap = index + 1
					}
				} else {
					// EOF
					f.linesInfo = append(f.linesInfo, [2]int{lineLength, 0})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				f.linesInfo = append(f.linesInfo, [2]int{lineLength, 1})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			f.linesInfo = append(f.linesInfo, [2]int{lineLength, -1})
		}
	}
	return f.linesInfo
}

func (f *File) Content() []byte {
	var err error
	if f.content == nil {
		f.content, err = os.ReadFile(f.path)
		if err != nil {
			f.HandleErr(err)
		}
	}
	return f.content
}

// AddDefect implements Context.
func (f *File) AddDefect(defect *Defect) {
	f.defects = append(f.defects, defect)
}

// Defects implements Context.
func (f *File) Defects() []*Defect {
	return f.defects
}

func (f *File) HandleErr(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func (f *File) String() string {
	return fmt.Sprintf("<File Context: %s>", f.path)
}

var _ Filer = (*File)(nil)

type PackageContext struct {
	Pkg
}

// IsPackage implements Context.
func (f *PackageContext) IsPackage() bool {
	return true
}

// AddDefect implements Context.
func (f *PackageContext) AddDefect(*Defect) {
	panic("unimplemented")
}

// Content implements Context.
func (f *PackageContext) Content() []byte {
	panic("unimplemented")
}

// Defects implements Context.
func (f *PackageContext) Defects() []*Defect {
	panic("unimplemented")
}

// Lines implements Context.
func (f *PackageContext) Lines() LinesInfo {
	panic("unimplemented")
}

// Path implements Context.
func (f *PackageContext) Path() string {
	panic("unimplemented")
}

func NewPackageContext(path string) *PackageContext {
	pkgCtx := PackageContext{
		Pkg: *NewPackage(),
	}
	return &pkgCtx
}

// Check implements Context.
func (f *PackageContext) Check() error {
	panic("PackageContext should not call the 'Check' method")
}

// AddDefect implements Context.
func (f *PackageContext) AST() *sitter.Tree {
	panic("PackageContext should not call the 'AST' method")
}

// AddDefect implements Context.
func (f *PackageContext) CallExpresses() map[string]*CallExpress {
	panic("PackageContext should not call the 'CallExpresses' method")
}

// HandleErr implements Context.
func (f *PackageContext) HandleErr(err error) {
	log.Error(err)
}

// Classes implements BaseContext.
func (f *PackageContext) Classes() map[string]*Class {
	panic("PackageContext should not call the 'Classes' method")
}

// Consts implements BaseContext.
func (f *PackageContext) Consts() map[string]*Const {
	panic("PackageContext should not call the 'Consts' method")
}

// Functions implements BaseContext.
func (f *PackageContext) Functions() map[string]*Function {
	panic("PackageContext should not call the 'Functions' method")
}

// Varibales implements BaseContext.
func (f *PackageContext) Varibales() map[string]*Variable {
	panic("PackageContext should not call the 'Varibales' method")
}

var _ Context = (*PackageContext)(nil)

type FileContext struct {
	File
	check CheckFuncType
}

// IsPackage implements Context.
func (f *FileContext) IsPackage() bool {
	return false
}

func NewFileContext(path string, check CheckFuncType) *FileContext {
	file := File{path: path}
	return &FileContext{File: file, check: check}
}

// AST implements Context.
func (f *FileContext) AST() *sitter.Tree {
	panic("unimplemented")
}

// AddDefect implements Context.
// Subtle: this method shadows the method (File).AddDefect of FileContext.File.
func (f *FileContext) AddDefect(*Defect) {
	panic("unimplemented")
}

// AddSubContext implements Context.
func (f *FileContext) AddSubContext(Context) {
	panic("unimplemented")
}

// CallExpresses implements Context.
func (f *FileContext) CallExpresses() map[string]*CallExpress {
	panic("unimplemented")
}

// Check implements Context.
func (f *FileContext) Check() error {
	return f.check(f)
}

// Classes implements Context.
func (f *FileContext) Classes() map[string]*Class {
	panic("unimplemented")
}

// Consts implements Context.
func (f *FileContext) Consts() map[string]*Const {
	panic("unimplemented")
}

// Content implements Context.
// Subtle: this method shadows the method (File).Content of FileContext.File.
func (f *FileContext) Content() []byte {
	panic("unimplemented")
}

// Defects implements Context.
// Subtle: this method shadows the method (File).Defects of FileContext.File.
func (f *FileContext) Defects() []*Defect {
	panic("unimplemented")
}

// Functions implements Context.
func (f *FileContext) Functions() map[string]*Function {
	panic("unimplemented")
}

// HandleErr implements Context.
// Subtle: this method shadows the method (File).HandleErr of FileContext.File.
func (f *FileContext) HandleErr(error) {
	panic("unimplemented")
}

// Lines implements Context.
// Subtle: this method shadows the method (File).Lines of FileContext.File.
func (f *FileContext) Lines() LinesInfo {
	panic("unimplemented")
}

// Package implements Context.
func (f *FileContext) Package() string {
	panic("unimplemented")
}

// Path implements Context.
// Subtle: this method shadows the method (File).Path of FileContext.File.
func (f *FileContext) Path() string {
	panic("unimplemented")
}

// Range implements Context.
func (f *FileContext) Range(func(ctx Context)) {
	panic("unimplemented")
}

// SetName implements Context.
func (f *FileContext) SetPkgName(name string) {
	panic("unimplemented")
}

// Varibales implements Context.
func (f *FileContext) Varibales() map[string]*Variable {
	panic("unimplemented")
}

var _ Context = (*FileContext)(nil)

type NewContextType func(string, CheckFuncType) Context

type PreHandlerType func(*config.Config, Context) error
