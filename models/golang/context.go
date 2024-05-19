package golang

import "github.com/stkali/glint/glint"

type Context struct {
	glint.FileContext
}

func (c *Context) Check(output glint.Outputer) {

}

func NewContext(path string) {

}

var model = glint.Model{
	GenerateModelFunc: func(model *glint.Model) (glint.ModelFuncType, error) {

		return func(ctx glint.FileContext) {
			ctx.Range(func(ctx glint.BaseContext) {

			})

		}, nil
	},
}
