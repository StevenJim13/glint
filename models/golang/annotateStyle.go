package golang

import (
	"reflect"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
)

const (
	disableMultiKey = "disableMulti"
	spaceInHeadKey  = "spaceInHead"
	allAsciiKey     = "allAscii"
)

type lintCommentType func(model *glint.Model, text string, point sitter.Point, ctx glint.Context)

var AnnotateStyleModel = glint.Model{
	Name: "AnnotateStyle",
	Options: map[string]any{
		// 禁用多行注释
		disableMultiKey: false,
		// 检查头有一个空格
		spaceInHeadKey: true,
		// 注释必须都是ascii码
		allAsciiKey: true,
	},
	Tags: []string{"basic"},
	GenerateModelFunc: func(model *glint.Model) (glint.ModelFuncType, error) {
		var lints = []lintCommentType{}

		if value, ok := model.Options[disableMultiKey]; ok {
			if disMulti, ok := value.(bool); !ok {
				return nil, errors.Newf("%q expected bool(true or false) but get %s", disableMultiKey, reflect.TypeOf(value))
			} else if disMulti {
				lints = append(lints, disableMultilineCommentLint)

			}
		}

		if value, ok := model.Options[spaceInHeadKey]; ok {
			if space, ok := value.(bool); !ok {
				return nil, errors.Newf("%q expected bool(true or false) but get %s", spaceInHeadKey, reflect.TypeOf(value))
			} else if space {
				lints = append(lints, spaceInCommentHeadLint)
			}
		}

		if value, ok := model.Options[allAsciiKey]; ok {
			if isAscii, ok := value.(bool); !ok {
				return nil, errors.Newf("%q expected bool(true or false) but get %s", allAsciiKey, reflect.TypeOf(value))
			} else if isAscii {
				lints = append(lints, pureAsciiLint)
			}
		}
		if len(lints) == 0 {
			return nil, nil
		}

		return func(ctx glint.Context) {

			root := ctx.ASTTree()
			if root == nil {
				return
			}
			content := ctx.Content()
			qc := sitter.NewQueryCursor()
			qc.Exec(queryCommentStmt, ctx.ASTTree().RootNode())
			for {
				m, ok := qc.NextMatch()
				if !ok {
					break
				}
				m = qc.FilterPredicates(m, content)
				for _, c := range m.Captures {
					node := c.Node
					text := node.Content(content)
					for index := range lints {
						lints[index](model, text, node.StartPoint(), ctx)
					}
				}
			}

		}, nil
	},
}

func disableMultilineCommentLint(model *glint.Model, text string, point sitter.Point, ctx glint.Context) {
	if text[1] == '*' {
		ctx.Defect(model, int(point.Row), int(point.Column),
			"Multi-line comments cannot be used, single-line comments are recommended.")
	}
}
func spaceInCommentHeadLint(model *glint.Model, text string, point sitter.Point, ctx glint.Context) {
	if len(text) > 2 {
		if char := text[2]; char != ' ' && char != '\n' && char != '\r' {
			ctx.Defect(model, int(point.Row), int(point.Column),
				"There must be at least one space between the comment character and the content")
		}
	}
}

func pureAsciiLint(model *glint.Model, text string, point sitter.Point, ctx glint.Context) {
	if !utils.IsPureAscii(text) {
		ctx.Defect(model, int(point.Row), int(point.Column),
			"There are non-ascii characters in the comment")
	}
}
