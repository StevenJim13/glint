/*
导出的函数、接口、变量、常量 、结构体 需要有注释
注释格式 注释需要使用单行并且要以 可导出的名称开头。

当函数很短很小，很容易理解时，可忽略注释
*/

package golang

import (
	"reflect"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stkali/glint/ast"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

const (
	// 没有
	noAnnotateFuncLinesKey = "noAnnotateFuncLines"
	// 是否检查全局变量
	checkVariableKey = "checkVariable"
	// 是否检查全局变量
	checkConstKey = "checkVariable"
	// 是否检查结构体
	checkClassKey = "checkClass"
	// 是否检查方法
	checkMethodKey = "checkMethod"
)

var MissAnnotateModel = glint.Model{
	Name: "MissAnnotate",
	Options: map[string]any{
		noAnnotateFuncLinesKey: 6,
		checkVariableKey:       true,
		checkClassKey:          true,
		checkMethodKey:         true,
	},
	Tags: []string{"basic"},
	GenerateModelFunc: func(model *glint.Model) (glint.CheckFuncType, error) {

		var linesNumber int = 0
		value, ok := model.Options[noAnnotateFuncLinesKey]
		if ok {
			if lines, ok := value.(int); !ok {
				return nil, errors.Newf("%s expected int but get %s", reflect.TypeOf(value))
			} else {
				linesNumber = lines
			}
		}

		var checkVariable bool = true
		value, ok = model.Options[checkVariableKey]
		if ok {
			if real, ok := value.(bool); !ok {
				return nil, errors.Newf("%s expected bool but get %s", reflect.TypeOf(value))
			} else {
				checkVariable = real
			}
		}

		var checkConst bool = true
		value, ok = model.Options[checkConstKey]
		if ok {
			if real, ok := value.(bool); !ok {
				return nil, errors.Newf("%s expected bool but get %s", reflect.TypeOf(value))
			} else {
				checkConst = real
			}
		}

		var checkClass bool = true
		value, ok = model.Options[checkClassKey]
		if ok {
			if real, ok := value.(bool); !ok {
				return nil, errors.Newf("%s expected bool but get %s", reflect.TypeOf(value))
			} else {
				checkClass = real
			}
		}

		return func(ctx glint.Context) error {
			log.Infof("apply %s model", model.Name)
			// exportable functions
			content := ctx.Content()
			for _, function := range ctx.Functions() {
				if linesNumber != 0 && ast.NodeLines(function.Node()) > linesNumber {
					continue
				}
				if !isAnnodate(function, content) {
					glint.AddDefect(ctx, model, function.Row(), function.Row(),
						"%q missing function annotate", function.Name())
				}
			}

			if checkClass {
				// exportable classes
				for name, class := range ctx.Classes() {

					if class.Node() != nil && !isAnnodate(class, content) {
						glint.AddDefect(ctx, model, class.Row(), class.Col(),
							"%q missing const annotate", name)
					}
					class.RangeMethod(func(method *glint.Method) {
						if !isAnnodate(method, content) {
							glint.AddDefect(ctx, model, method.Row(), method.Col(),
								"%q missing method annotate", method.Name())
						}
					})
				}
			}

			if checkConst {

				// exportable consts
				for _, cst := range ctx.Consts() {
					if !isAnnodate(cst, content) {
						glint.AddDefect(ctx, model, cst.Row(), cst.Col(),
							"%q missing const annotate", cst.Name())
					}
				}
			}

			if checkVariable {
				// exportable variables
				for _, variable := range ctx.Varibales() {
					if !isAnnodate(variable, content) {
						glint.AddDefect(ctx, model, variable.Row(), variable.Col(),
							"%q missing const annotate", variable.Name())
					}
				}
			}
			return nil
		}, nil
	},
}

func annotateBy(text, name string) bool {
	return strings.HasPrefix(strings.TrimLeft(text[2:], " "), name)
}

const (
	Name string = "111"
)

func isAnnodate(elem glint.Elementer, content []byte) bool {
	name := elem.Name()

	if name[0] < 'A' || name[0] > 'Z' {
		return true
	}
	if pre := elem.Node().PrevSibling(); pre.Type() == "comment" {
		header := ast.QueryCommentHeader(pre, func(n *sitter.Node) bool { return n.Type() == "comment" })
		if annotateBy(header.Content(content), name) {
			return true
		}
	}
	return false
}
