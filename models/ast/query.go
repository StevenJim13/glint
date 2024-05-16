package ast

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/stkali/utility/tool"
)

var (
	Language         = golang.GetLanguage()
	queryCommentStmt *sitter.Query
	queryFuncDefine  *sitter.Query
	err              error
)

func init() {

	//
	queryCommentStmt, err = sitter.NewQuery(tool.ToBytes("(comment) @comment"), Language)
	tool.CheckError("failed to build comment query", err)

	//
	queryFuncDefine, err = sitter.NewQuery(tool.ToBytes(""), Language)
	tool.CheckError("failed to build function define query", err)

}
