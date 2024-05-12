package models

import (
	"fmt"

	"github.com/stkali/glint/parser"
)

type Defect struct {
	desc string
}

func (d *Defect) String() string {
	return fmt.Sprintf("<Defect: %s>", d.desc)
}

func NewDefect(desc string) parser.Defect {
	return &Defect{desc: desc}
}

var _ parser.Defect = (*Defect)(nil)
