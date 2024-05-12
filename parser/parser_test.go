package parser

import (
	"fmt"
	"os"
	"testing"
)

func TestGetFileInfo(t *testing.T) {
	fmt.Println(os.Getwd())
	ctx := CContext{file: "../testdata/newlines/windows.txt"}
	ctx.LinesInfo()
	ctx = CContext{file: "../testdata/newlines/mac.txt"}
	ctx.LinesInfo()
}
