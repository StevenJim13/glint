package golang

import "github.com/stkali/glint/glint"

var FunctionModel = glint.Model{
	Name: "Function",
	Tags: []string{"basic"},
	Options: map[string]any{
		"": "",
	},
	GenerateModelFunc: func(model *glint.Model) (glint.CheckFuncType, error) {
		return func(ctx glint.Context) error {
			return nil
		}, nil
	},
}
