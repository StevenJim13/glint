package glint

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
)

type Outputer interface {
	Write(ctx Context)
	Flush()
}

// CreateOutput ...
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
	sync.Mutex
}

func NewTextOutput(fd io.Writer) Outputer {
	return &TextOutput{output: fd}
}

// Flush implements Outputer.
func (c *TextOutput) Flush() {
	utils.Close(c.output)
}

func (t *TextOutput) Write(ctx Context) {
	if len(ctx.DefectSet()) == 0 {
		return
	}
	t.Lock()
	defer t.Unlock()
	fmt.Println(ctx.File())
	for id, d := range ctx.DefectSet() {
		fmt.Printf("%6d|(%4d,%4d) model:%s desc:%s\n", id, d.Row, d.Col, d.Model.Name, d.Desc)
	}
}

func (t *TextOutput) String() string {
	return fmt.Sprintf("<TextOutput: %p>", t.output)
}
