package output

import (
	"fmt"

	"github.com/bytedance/sonic"
)

type Outputer interface {
	AddDefectSet(file string, ds *DefectSet)
}

type DefectSet []*Defect


func (s DefectSet) Add(defect *Defect) {
	s = append(s, defect)
}

func (s DefectSet) Json() {

}

type Defect struct {
	// Model 规则的名称
	Model *string
	// 缺陷的描述
	Desc string
	// Line
	Row int
	Col  int
}

func (d *Defect) Json() ([]byte, error) {
	return sonic.Marshal(d)
}

func (d *Defect) Cmd() string {
	return fmt.Sprintf("[%d, %d]\ttype: %s desc: %s", d.Row, d.Col, *d.Model, d.Desc)
}
