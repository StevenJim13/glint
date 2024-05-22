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
	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
)

type NewContextType func(string, CheckFuncType) Context

type PreHandlerType func(*config.Config, Packager) error

type Packager interface {
	Path() string
	Name() string
	SetName(name string)
	AddContext(Context)
	AddPackage(Packager)
	Range(func(ctx Context))
	RangePackages(func(pkg Packager))
	Walk(fn func(ctx Context))
	fmt.Stringer
}

type Sourcer interface {
	Tree() *sitter.Tree
	Root() *sitter.Node
	Functions() map[string]*Function
	Classes() map[string]*Class
	CallExpresses() map[string]*CallExpress
	Consts() map[string]*Const
	Varibales() map[string]*Variable
}

type Filer interface {
	Path() string
	Lines() Lines
	Content() []byte
}

type Context interface {
	DefectSeter
	Filer
	fmt.Stringer
	Check() error
	HandleErr(error)
	Source() Sourcer
	Package() Packager
	LinkPackage(Packager)
}

type Package struct {
	name     string
	path     string
	children []Context
	packages []Packager
}

// AddPackage implements Packager.
func (p *Package) AddPackage(pkg Packager) {
	p.packages = append(p.packages, pkg)
}

// RangePackages implements Packager.
func (p *Package) RangePackages(fn func(pkg Packager)) {
	for index := range p.packages {
		fn(p.packages[index])
	}
}

// AddContext implements Packager.
func (p *Package) AddContext(ctx Context) {
	log.Info(ctx)
	ctx.LinkPackage(p)
	p.children = append(p.children, ctx)
}

// Range implements Packager.
func (p *Package) Range(fn func(ctx Context)) {
	for index := range p.children {
		fn(p.children[index])
	}
}

// SetName implements Packager.
func (p *Package) Walk(fn func(ctx Context)) {
	p.Range(fn)
	p.RangePackages(func(pkg Packager) {
		pkg.Walk(fn)
	})
}

// Name implements Packager.
func (p *Package) Name() string {
	return p.name
}

// Path implements Packager.
func (p *Package) Path() string {
	return p.path
}

// SetName implements Packager.
func (p *Package) SetName(name string) {
	p.name = name
}

// SetName implements Packager.
func (p *Package) String() string {
	if p.name != "" {
		return fmt.Sprintf("<Package: %s>", p.name)
	}
	return fmt.Sprintf("<Package: %s>", p.path)
}

func NewPackage(path string) Packager {
	return &Package{
		path: path,
	}
}

var _ Packager = (*Package)(nil)

// -------------------------------------
type File struct {
	path    string
	lines   Lines
	content []byte
	pkg     Packager
	DefectSet
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
func (f *File) Lines() Lines {
	if f.lines == nil {
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
						f.lines = append(f.lines, [2]int{lineLength, 2})
						index += 1
						gap = index + 1
					} else {
						// \r
						f.lines = append(f.lines, [2]int{lineLength, 0})
						gap = index + 1
					}
				} else {
					// EOF
					f.lines = append(f.lines, [2]int{lineLength, 0})
				}
			case '\n':
				lineLength := len(tool.ToString(ctt[gap:index]))
				f.lines = append(f.lines, [2]int{lineLength, 1})
				gap = index + 1
			}
			index += 1
		}
		if index > gap {
			lineLength := len(tool.ToString(ctt[gap:index]))
			f.lines = append(f.lines, [2]int{lineLength, -1})
		}
	}
	return f.lines
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

func (f *File) HandleErr(err error) {
	fmt.Fprintln(os.Stderr, err)
	panic(err)
}

func (f *File) Package() Packager {
	return f.pkg
}

func (f *File) LinkPackage(pkg Packager) {
	f.pkg = pkg
}

func (f *File) String() string {
	return fmt.Sprintf("<FileContext: %s>", f.path)
}

type Lines [][2]int

func (l Lines) Range(fn func(index int, item [2]int) bool) {
	for index, item := range l {
		if !fn(index, item) {
			return
		}
	}
}

type FileContext struct {
	File
	check     CheckFuncType
	handleErr func(error)
}

var _ Context = (*FileContext)(nil)

func NewFileContext(file string, check CheckFuncType) Context {
	filepath := File{path: file}
	ctx := &FileContext{
		File:  filepath,
		check: check,
	}
	return ctx
}

func (f *FileContext) HandleErr(err error) {
	f.handleErr(err)
}

func (f *FileContext) Check() error {
	return f.check(f)
}

func (f *FileContext) Source() Sourcer {
	return nil
}
