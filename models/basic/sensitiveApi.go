package basic

import (
	"fmt"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/parser"
)

var (
	sensitiveKey = "sensitives"
)

var SensitiveApi = models.Model{
	Name: "SensitiveApi",
	Tags: []string{"basic"},
	Options: map[string]any{
		sensitiveKey: []string{
			"foo", "bar",
		},
	},
	ModelFunc: func(model *models.Model, ctx parser.Context) {

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
				ctx.AddDefect(models.NewDefect(fmt.Sprintf("sensitive api: %q", call.Function.Name)))
			}
		}
	},
}
