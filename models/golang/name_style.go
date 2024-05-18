package golang

import (
	"fmt"
	"os"

	"github.com/stkali/glint/glint"
	"golang.org/x/tools/go/packages"
)

const (
	filenameCheckKey  = "filenameCheck"
	directoryCheckKey = "directoryCheck"
	funcName
)

var TestModel = glint.Model{
	Name: "Test",
	GenerateModelFunc: func(model *glint.Model) (glint.ModelFuncType, error) {
		pkgs, err := packages.Load(&packages.Config{
			Mode: packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
		}, "/Users/kali/develop/glint")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for index := range pkgs {
			fmt.Println("package: ", pkgs[index])
		}

		return func(ctx glint.Context) {

		}, nil
	},
}

// var NameStyleModel = glint.Model{
// 	Name: "NameStyle",
// 	Options: map[string]any{
// 		constNameStyleKey: "upper",
// 		globalNameStyleKey: ""
// 		// 禁用多行注释
// 		funcNameStyleKey: false,
// 		// 检查头有一个空格
// 		globalNameStyleKey: true,
// 		// 注释必须都是ascii码
// 		constNameStyleKey: true
// 		methodNameStyleKey: ,
// 	},
// 	Tags: []string{"basic"},
// 	GenerateModelFunc: func(model *glint.Model) (glint.ModelFuncType, error) {
// 		var lints = []lintCommentType{}

// 		if value, ok := model.Options[disableMultiKey]; ok {
// 			if disMulti, ok := value.(bool); !ok {
// 				return nil, errors.Newf("%q expected bool(true or false) but get %s", disableMultiKey, reflect.TypeOf(value))
// 			} else if disMulti {
// 				lints = append(lints, disableMultilineCommentLint)

// 			}
// 		}

// 		if value, ok := model.Options[spaceInHeadKey]; ok {
// 			if space, ok := value.(bool); !ok {
// 				return nil, errors.Newf("%q expected bool(true or false) but get %s", spaceInHeadKey, reflect.TypeOf(value))
// 			} else if space {
// 				lints = append(lints, spaceInCommentHeadLint)
// 			}
// 		}

// 		if value, ok := model.Options[allAsciiKey]; ok {
// 			if isAscii, ok := value.(bool); !ok {
// 				return nil, errors.Newf("%q expected bool(true or false) but get %s", allAsciiKey, reflect.TypeOf(value))
// 			} else if isAscii {
// 				lints = append(lints, pureAsciiLint)
// 			}
// 		}
// 		if len(lints) == 0 {
// 			return nil, nil
// 		}

// 		return func(ctx glint.Context) {

// 			root := ctx.ASTTree()
// 			if root == nil {
// 				return
// 			}
// 			content := ctx.Content()
// 			qc := sitter.NewQueryCursor()
// 			qc.Exec(queryCommentStmt, ctx.ASTTree().RootNode())
// 			for {
// 				m, ok := qc.NextMatch()
// 				if !ok {
// 					break
// 				}
// 				m = qc.FilterPredicates(m, content)
// 				for _, c := range m.Captures {
// 					node := c.Node
// 					text := node.Content(content)
// 					for index := range lints {
// 						lints[index](model, text, node.StartPoint(), ctx)
// 					}
// 				}
// 			}

// 		}, nil
// 	},
// }
