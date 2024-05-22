package golang

import (
	"fmt"

	"github.com/stkali/glint/glint"
)

type Context struct {
	glint.File
	check  glint.CheckFuncType
	source glint.Sourcer
}

// Source implements glint.Context.
func (c *Context) Source() glint.Sourcer {
	return c.source
}

func (c *Context) Check() error {
	return c.check(c)
}

func (c *Context) String() string {
	return fmt.Sprintf("<GolangContent: %s>", c.Path())
}

func (c *Context) HandleErr(err error) {
	fmt.Println(err)
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
