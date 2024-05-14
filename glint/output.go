package glint

import (
	"fmt"
	"io"
	"os"

	"github.com/stkali/utility/errors"
)

type Outputer interface {
	Write(ctx Context)
	Flush()
}

func CreateOutput(file, format string) (Outputer, error) {

	var writer io.Writer
	var outer Outputer
	if file == "" {
		writer = os.Stdout
	} else {
		if fd, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.ModePerm); err != nil {
			writer = fd
		} else {
			return nil, errors.Newf("cannot create result outputer, err: %s", err)
		}
	}
	switch format {
	case "json":
		outer = NewJsonOutput(writer)
	default:
		outer = NewTextOutput(writer)
	}
	return outer, nil
}

type JsonOutput struct {
	output io.Writer
}

func NewJsonOutput(fd io.Writer) Outputer {
	return &JsonOutput{output: fd}
}

// Flush implements Outputer.
func (j *JsonOutput) Flush() {
	if closer, ok := j.output.(io.Closer); ok {
		closer.Close()
	}
}

func (j *JsonOutput) Write(ctx Context) {
	// j.output.Write()
}

type TextOutput struct {
	output io.Writer
}

// Flush implements Outputer.
func (c *TextOutput) Flush() {
	panic("unimplemented")
}

func NewTextOutput(fd io.Writer) Outputer {
	return &TextOutput{output: fd}
}

func (c *TextOutput) Write(ctx Context) {
	if len(ctx.DefectSet()) == 0 {
		return
	}
	fmt.Println(ctx.File())
	for id, d := range ctx.DefectSet() {
		fmt.Printf("  %d. model:%s position:(%d,%d) desc:%s\n", id, d.Row, d.Col, *d.Model, d.Desc)
	}
}
