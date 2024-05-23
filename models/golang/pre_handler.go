package golang

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/errors"
)

type Source struct {
	tree *sitter.Tree
	root *sitter.Node
	ts   map[string]*glint.Type
	fs   map[string]*glint.Function
	vs   map[string]*glint.Value
	cs   map[string]*glint.Value
}

// Calls implements glint.Sourcer.
func (s *Source) Calls() map[string]*glint.Call {
	panic("unimplemented")
}

// Consts implements glint.Sourcer.
func (s *Source) Consts() map[string]*glint.Value {
	panic("unimplemented")
}

// Functions implements glint.Sourcer.
func (s *Source) Functions() map[string]*glint.Function {
	panic("unimplemented")
}

// Root implements glint.Sourcer.
func (s *Source) Root() *sitter.Node {
	return s.root
}

// Tree implements glint.Sourcer.
func (s *Source) Tree() *sitter.Tree {
	return s.tree
}

// Types implements glint.Sourcer.
func (s *Source) Types() map[string]*glint.Type {
	panic("unimplemented")
}

// Varibales implements glint.Sourcer.
func (s *Source) Varibales() map[string]*glint.Value {
	panic("unimplemented")
}

func NewSource(ctx *Context) (*Source, error) {
	content := ctx.Content()
	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())
	tree, err := parser.ParseCtx(context.TODO(), nil, content)
	if err != nil {
		return nil, err
	}
	src := Source{tree: tree, root: tree.RootNode()}

	return &src, nil
}

var _ glint.Sourcer = (*Source)(nil)

// func parseMethod(node *sitter.Node, content []byte) *Method {

// 	method := &Method{}
// 	if nameNode := node.ChildByFieldName("name"); nameNode != nil {
// 		method.name = nameNode.Content(content)
// 		method.node = node
// 	}
// 	if paramNode := node.ChildByFieldName("receiver"); paramNode.Type() == "parameter_list" {
// 		if decl := paramNode.Child(1); decl.Type() == "parameter_declaration" {
// 			if reciverNode := decl.ChildByFieldName("name"); reciverNode != nil {
// 				method.reciver = reciverNode.Content(content)
// 			}
// 			typeNode := decl.ChildByFieldName("type")
// 			if typeNode == nil {
// 				return nil
// 			}
// 			switch typeNode.Type() {
// 			case "type_identifier": // 绑定到了实例上
// 				method.instance = NewClass(typeNode.Content(content), nil)
// 			case "pointer_type": // 绑定到地址
// 				if !handlePointType(method, typeNode, content) {
// 					return nil
// 				}
// 			case "generic_type":
// 				if !handleGenericType(method, typeNode, content) {
// 					return nil
// 				}
// 			default:
// 				return nil
// 			}
// 		}
// 	}
// 	return method
// }

// func handlePointType(method *Method, node *sitter.Node, content []byte) bool {
// 	method.prt = true
// 	first := node.Child(1)
// 	if first == nil {
// 		return false
// 	}
// 	switch first.Type() {
// 	case "generic_type":
// 		return handleGenericType(method, first, content)
// 	case "type_identifier":
// 		method.instance = NewClass(first.Content(content), nil)
// 	default:
// 		return false
// 	}
// 	return true
// }
// func handleGenericType(method *Method, node *sitter.Node, content []byte) bool {
// 	if identifier := node.ChildByFieldName("type"); identifier != nil && identifier.Type() == "type_identifier" {
// 		method.instance = NewClass(identifier.Content(content), nil)
// 		return true
// 	}
// 	return false
// }

// 解析整个包，然后构建出包级别的关联关系
// 解析单个文件，
// func N1ewSource(ctx *Context) (*Source, error) {
// 	content := ctx.Content()
// 	parser := sitter.NewParser()
// 	parser.SetLanguage(golang.GetLanguage())
// 	tree, err := parser.ParseCtx(context.TODO(), nil, content)
// 	if err != nil {
// 		return nil, err
// 	}
// 	root := tree.RootNode()
// 	ast.ApplyChildrenNodes(root, func(sub *sitter.Node) {
// 		switch sub.Type() {
// 		case "method_declaration":
// 			method := parseMethod(sub, content)
// 			if method == nil {
// 				utils.Bugf("cannot parse method: %s", sub.Content(content))
// 			} else {
// 				if class, ok := g.classes[method.instance.name]; !ok {
// 					g.classes[method.instance.name] = glint.NewClass(method.instance.name, nil)
// 				} else {
// 					class.AddMethod(method)
// 				}
// 			}

// 		case "const_declaration":
// 			ast.ApplyChildrenNodes(sub, func(spec *sitter.Node) {
// 				if spec.Type() == "const_spec" {
// 					nameNode := spec.ChildByFieldName("name")
// 					if nameNode != nil {
// 						g.consts = append(g.consts, &Const{BaseElem{name: nameNode.Content(content), node: spec}})
// 					} else {
// 						utils.Bugf("cannot parse const desclare: %s", sub.Content(content))
// 					}
// 				}
// 			})

// 		case "var_declaration":
// 			ast.ApplyChildrenNodes(sub, func(spec *sitter.Node) {
// 				if spec.Type() == "var_spec" {
// 					nameNode := spec.ChildByFieldName("name")
// 					if nameNode != nil {
// 						g.consts = append(g.consts, &Const{BaseElem{name: nameNode.Content(content), node: spec}})
// 					} else {
// 						utils.Bugf("cannot parse var desclare: %s", sub.Content(content))
// 					}
// 				}
// 			})
// 		case "function_declaration":
// 			if nameNode := sub.ChildByFieldName("name"); nameNode.Type() == "identifier" {
// 				g.funcs = append(g.funcs, &Function{BaseElem{name: nameNode.Content(content), node: sub}})
// 			} else {
// 				utils.Bugf("cannot parse function desclare: %s", sub.Content(content))
// 			}
// 		case "type_declaration":
// 			if spec := sub.Child(1); spec.Type() == "type_spec" {
// 				if nameNode := spec.ChildByFieldName("name"); nameNode.Type() == "type_identifier" {
// 					name := nameNode.Content(content)
// 					if class, ok := g.classes[name]; !ok {
// 						g.classes[name] = NewClass(name, sub)
// 					} else {
// 						class.node = sub
// 					}
// 				} else {
// 					utils.Bugf("cannot parse struct declare: %s", sub.Content(content))
// 				}
// 			}
// 		}
// 	})

// 	return nil, nil
// }

var _ glint.Sourcer = (*Source)(nil)

// PreHandle ...
func PreHandle(conf *config.Config, pkg glint.Packager) error {

	// var err error
	pkg.Range(func(ctx glint.Context) {
		gctx, ok := ctx.(*Context)
		if !ok {
			errors.Join(err, errors.Newf("failed to assert %s to *golang.Context"), ctx)
			return
		}
		source, nErr := NewSource(gctx)
		if err != nil {
			err = errors.Join(err, nErr)
			return 
		}
		gctx.source = source
	})

	// 单文件的解析结束，开始合并
		// 1 不在当前文件实现的方法
		// 2 绑定到当前实现上





	return nil
}
