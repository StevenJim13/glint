package python

import (
	"os"

	"github.com/stkali/glint/parser"
	"github.com/stkali/utility/log"
)

type FileContext struct {
	file    string
	content []byte
}

// AddDefect implements parser.Context.
func (f *FileContext) AddDefect(defect parser.Defect) {
	panic("unimplemented")
}

// CallExpresses implements parser.Context.
func (f *FileContext) CallExpresses() []*parser.CallExpress {
	panic("unimplemented")
}

// Content implements parser.Context.
func (f *FileContext) Content() []byte {
	var err error
	if f.content == nil {
		if f.content, err = os.ReadFile(f.file); err != nil {
			// errors.Is(err, fs.Err)
			log.Errorf("failed to read file: %q, err: %s", f.file, err)
		}
	}
	return f.content
}

// Functions implements parser.Context.
func (f *FileContext) Functions() []*parser.Function {
	panic("unimplemented")
}

// LinesInfo implements parser.Context.
func (f *FileContext) LinesInfo() parser.LinesInfo {
	panic("unimplemented")
}

var _ parser.Context = (*FileContext)(nil)
