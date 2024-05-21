package glint

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
)

type DefectSeter interface {
	AddDefect(*Defect)
	Defects() []*Defect
}

type DefectSet []*Defect

// AddDefect implements DefectSeter.
func (d DefectSet) AddDefect(defect *Defect) {
	d = append(d, defect)
}

// Defects implements DefectSeter.
func (d DefectSet) Defects() []*Defect {
	return d
}

var _ DefectSeter = (*DefectSet)(nil)

type Defect struct {
	Model *Model
	Desc  string
	Row   int
	Col   int
}

func NewDefect(model *Model, row, col int, desc string, args ...any) *Defect {
	return &Defect{
		Model: model,
		Desc:  fmt.Sprintf(desc, args...),
		Row:   row,
		Col:   col,
	}
}

func (d *Defect) String() string {
	return fmt.Sprintf("model: %q, desc: %s, position:(%d,%d)", *&d.Model.Name, d.Desc, d.Row, d.Col)
}

func AddDefect(ctx Context, model *Model, row, col int, desc string, args ...any) {
	def := NewDefect(model, row, col, desc, args...)
	ctx.AddDefect(def)
}

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
	output    io.Writer
	bufWriter *bufio.Writer
	encoder   *json.Encoder
	sync.Mutex
}

// String implements Outputer.
func (j *JsonOutput) String() string {
	return fmt.Sprintf("<JsonOutput: %p>", j.output)
}

func NewJsonOutput(fd io.Writer) Outputer {
	writer := bufio.NewWriter(fd)
	return &JsonOutput{
		bufWriter: bufio.NewWriter(fd),
		encoder:   json.NewEncoder(writer),
		output:    fd,
	}
}

// Flush implements Outputer.
func (j *JsonOutput) Close() {
	j.bufWriter.Flush()
	utils.Close(j.output)
}

func (j *JsonOutput) Write(ctx Context) {
	if len(ctx.Defects()) == 0 {
		return
	}
	j.Lock()
	defer j.Unlock()
	v := map[string][]*Defect{
		ctx.Path(): ctx.Defects(),
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
	if len(ctx.Defects()) == 0 {
		return
	}
	t.Lock()
	defer t.Unlock()
	fmt.Fprintln(t.output, ctx.Path())
	for id, d := range ctx.Defects() {
		fmt.Fprintf(t.output, "%6d|(%4d,%4d) model:%s desc:%s\n", id, d.Row, d.Col, d.Model.Name, d.Desc)
	}
}

func (t *TextOutput) String() string {
	return fmt.Sprintf("<TextOutput: %p>", t.output)
}
