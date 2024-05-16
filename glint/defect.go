package glint

import "fmt"

type Defect struct {
	Model *Model
	Desc  string
	Row   int
	Col   int
}

func (d *Defect) String() string {
	return fmt.Sprintf("model: %q, desc: %s, position:(%d,%d)", *&d.Model.Name, d.Desc, d.Row, d.Col)
}
