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
	check glint.CheckFuncType
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
