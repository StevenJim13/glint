package golang

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/stkali/glint/glint"
)

type Context struct {
	glint.File
	ast       *sitter.Tree
	functions map[string]*glint.Function
	classes   map[string]*glint.Class
	varibales map[string]*glint.Variable
	consts    map[string]*glint.Const
	calls     []*glint.CallExpress
	check     glint.CheckFuncType
}

// CallExpresses implements glint.Context.
func (c *Context) CallExpresses() map[string]*glint.CallExpress {
	return nil
}

// Classes implements glint.Context.
func (c *Context) Classes() map[string]*glint.Class {
	return nil
}

// Consts implements glint.Context.
func (c *Context) Consts() map[string]*glint.Const {
	return nil
}

// Functions implements glint.Context.
func (c *Context) Functions() map[string]*glint.Function {
	return nil
}

// Varibales implements glint.Context.
func (c *Context) Varibales() map[string]*glint.Variable {
	return nil
}

// Package implements glint.Context.
func (c *Context) Package() string {
	panic("PackageContext should not call the 'Package' method")
}

// AddSubContext implements glint.Context.
func (c *Context) AddSubContext(glint.Context) {
	panic("PackageContext should not call the 'AddSubContext' method")
}

// Range implements glint.Context.
func (c *Context) Range(func(ctx glint.Context)) {
	panic("PackageContext should not call the 'Range' method")
}

func (c *Context) Check() error {
	return c.check(c)
}

func (c *Context) String() string {
	return fmt.Sprintf("<GolangContent: %s>", c.Path())
}
func (c *Context) IsPackage() bool {
	return false
}
func (c *Context) SetPkgName(name string) {

}

func (c *Context) AST() *sitter.Tree {
	if c.ast != nil {
		return c.ast
	}
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())
	if tree, err := parser.ParseCtx(context.Background(), nil, c.Content()); err != nil {
		c.HandleErr(err)
	} else {
		c.ast = tree
	}

	return c.ast
}

var _ glint.Context = (*Context)(nil)

func NewContext(path string, check glint.CheckFuncType) glint.Context {
	file := glint.NewFile(path)
	ctx := Context{
		File:  *file,
		check: check,
	}
	return &ctx
}
