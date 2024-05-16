package golang

import (
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/log"
)

var AnnotateModel = glint.Model{
	Name:    "Annotate",
	Options: map[string]any{},
	Tags:    []string{"basic"},
	ModelFunc: func(model *glint.Model, ctx glint.Context) {
		root := ctx.ASTTree()
		log.Infof("success get ast root")
		if root == nil {
			return
		}
		content := ctx.Content()
		// log.Infof("successfully parse ast: %s", ctx.File())
		// ast.InspectTree(root.RootNode(), content, "    ")
		query, _ := sitter.NewQuery([]byte("(function_declaration) @comment"), golang.GetLanguage())
		qc := sitter.NewQueryCursor()
		qc.Exec(query, ctx.ASTTree().RootNode())
		for {
			m, ok := qc.NextMatch()
			if !ok {
				break
			}
			// Apply predicates filtering
			m = qc.FilterPredicates(m, content)
			for _, c := range m.Captures {
				fmt.Println("--------", c.Node.Content(content))
			}
		}
	},
}
