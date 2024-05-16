package golang

import (
	"testing"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
)

func createTestContext() glint.Context {

	return nil
}

func TestAnnotateModel(t *testing.T) {

	content := `
// Outputer ...
type Outputer interface {
	// Write ....
	Write(ctx Context)
	// Flush ...
	Flush()
}
// Function ..
func GetName() string{
}`
	ctx := glint.CreateTestFileNode(utils.GoLang, "annotate.go", []byte(content))
	AnnotateModel.ModelFunc(&AnnotateModel, ctx)

}
