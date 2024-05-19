package glint

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
)

type Outputer interface {
	Write(ctx Context)
	Close()
	fmt.Stringer
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
	bufWriter *bufio.Writer
	sync.Mutex
	encoder sonic.Encoder
}

func NewJsonOutput(fd io.Writer) Outputer {
	writer := bufio.NewWriter(fd)
	return &JsonOutput{
		bufWriter: bufio.NewWriter(fd),
		encoder: sonic.ConfigDefault.NewEncoder(writer),
		output:  fd,
	}
}

// Flush implements Outputer.
func (j *JsonOutput) Close() {
	j.bufWriter.Flush()
	utils.Close(j.output)
}

func (j *JsonOutput) Write(ctx Context) {
	if len(ctx.DefectSet()) == 0 {
		return
	}
	j.Lock()
	defer j.Unlock()
	v := map[string][]*Defect{
		ctx.File(): ctx.DefectSet(),
	}
	j.encoder.Encode(v)
}

type TextOutput struct {
	output io.Writer
	sync.Mutex
}

func NewTextOutput(fd io.Writer) Outputer {
	return &TextOutput{output: fd}
}

// Flush implements Outputer.
func (c *TextOutput) Close() {
	utils.Close(c.output)
}

func (t *TextOutput) Write(ctx Context) {
	if len(ctx.DefectSet()) == 0 {
		return
	}
	t.Lock()
	defer t.Unlock()
	fmt.Fprintln(t.output, ctx.File())
	for id, d := range ctx.DefectSet() {
		fmt.Fprintf(t.output, "%6d|(%4d,%4d) model:%s desc:%s\n", id, d.Row, d.Col, d.Model.Name, d.Desc)
	}
}

func (t *TextOutput) String() string {
	return fmt.Sprintf("<TextOutput: %p>", t.output)
}
