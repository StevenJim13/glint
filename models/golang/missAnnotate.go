/*
导出的函数、接口、变量、常量 、结构体 需要有注释
注释格式 注释需要使用单行并且要以 可导出的名称开头。

当函数很短很小，很容易理解时，可忽略注释
*/

package golang

import (
	"reflect"
	"strings"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/models/ast"
	"github.com/stkali/utility/errors"
)

const (
	noAnnotateFuncLinesKey = "noAnnotateFuncLines"
)

var MissAnnotateModel = glint.Model{
	Name: "MissAnnotate",
	Options: map[string]any{
		noAnnotateFuncLinesKey: 6,
	},
	Tags: []string{"basic"},
	GenerateModelFunc: func(model *glint.Model) (glint.ModelFuncType, error) {
		var linesNumber *int
		value, ok := model.Options[noAnnotateFuncLinesKey]
		if ok {
			if lines, ok := value.(int); !ok {
				return nil, errors.Newf("%s expected int but get %s", reflect.TypeOf(value))
			} else {
				linesNumber = &lines
			}
		}
		return func(ctx glint.Context) {
			content := ctx.Content()
			for _, function := range ctx.Functions() {

				if ast.NodeLines(function.Node()) < *linesNumber {
					continue
				}
				if pre := function.Node().PrevSibling(); pre.Type() == "comment" && annotateBy(pre.Content(content), function.Name()) {
					continue
				} else {
					ctx.Defect(model, function.Row(), function.Row(), "missing function annotate")
				}
			}
		}, nil
	},
}

func annotateBy(text, name string) bool {
	return strings.HasPrefix(strings.TrimLeft(text[2:], " "), name)
}
