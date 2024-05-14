package basic

import (
	"github.com/stkali/glint/glint"
)

var (
	sensitiveKey = "sensitives"
)

var SensitiveApi = glint.Model{
	Name: "SensitiveApi",
	Tags: []string{"basic"},
	Options: map[string]any{
		sensitiveKey: []string{
			"foo", "bar",
		},
	},
	ModelFunc: func(model *glint.Model, ctx glint.Context) {

		sensitiveFuncs, ok := model.Options[sensitiveKey]
		if !ok {
			return
		}
		sensList, ok := sensitiveFuncs.([]string)
		if !ok {
			return
		}
		// build sensitive api hash table
		sensHashTable := make(map[string]struct{}, len(sensList))
		for index := range sensList {
			sensHashTable[sensList[index]] = struct{}{}
		}

		for _, call := range ctx.CallExpresses() {
			if _, ok = sensHashTable[call.Function.Name]; ok {
				p := call.Function.Position
				ctx.Defect(model, p[0], p[1],
					"sensitive api: %q", call.Function.Name,
				)
			}
		}
	},
}
