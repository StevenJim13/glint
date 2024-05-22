package golang

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/log"
)

// PreHandle ...
func PreHandle(conf *config.Config, pkg glint.Packager) error {
	ctx := NewContext("glint/glint.go", nil)
	content := ctx.Content()

	parser := sitter.NewParser()
	parser.SetLanguage(golang.GetLanguage())
	tree, err := parser.ParseCtx(context.TODO(), nil, content)
	if err != nil {
		fmt.Println("----------", err)
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(tree)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	v := gob.NewDecoder(&buf)
	var newT sitter.BaseTree
	err = v.Decode(&newT)
	// fmt.Println("xxxxxxxx", buf.Bytes())
	fmt.Println("xxxx11111x", err)
	log.Infof("apply golang pre handle!")
	return nil
}
