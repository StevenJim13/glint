package basic

import (
	"reflect"

	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

var (
	sensitiveKey = "sensitives"
)

var SensitiveApiModel = glint.Model{
	Name: "SensitiveApi",
	Tags: []string{"basic"},
	Options: map[string]any{
		sensitiveKey: []string{
			"foo", "bar", "function",
		},
	},
	GenerateModelFunc: func(model *glint.Model) (glint.CheckFuncType, error) {

		value, ok := model.Options[sensitiveKey]
		if !ok {
			return nil, nil
		}
		sensitives, ok := value.([]any)
		if !ok {
			return nil, errors.Newf("%s expected []any{} but get %s", sensitiveKey, reflect.TypeOf(value))
		}
		sensTable := make(map[string]struct{}, len(sensitives))
		for index := range sensitives {
			if sen, ok := sensitives[index].(string); ok {
				sensTable[sen] = struct{}{}
			} else {
				return nil, errors.Newf("%s expected item type is string but get %s", reflect.TypeOf(sensitives[index]))
			}
		}
		return func(ctx glint.Context) error {
			log.Infof("apply %s model", model.Name)
			source := ctx.Source()
			if source == nil {
				return errors.Newf("failed to get %q source", ctx.Path())
			}

			for _, call := range source.CallExpresses() {
				if _, ok = sensTable[call.Function.Name()]; ok {
					glint.AddDefect(ctx, model, call.Function.Row(), call.Function.Col(),
						"sensitive api: %q", call.Function.Name)
				}
			}
			return nil
		}, nil
	},
}
