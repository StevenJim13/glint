package c

import (
	"os"

	"github.com/stkali/glint/models"
	"github.com/stkali/utility/log"
)

type FileContext struct {
	file    string
	content []byte
}

// AddDefect implements models.Context.
func (f *FileContext) AddDefect(defect models.Defect) {
	panic("unimplemented")
}

// CallExpresses implements models.Context.
func (f *FileContext) CallExpresses() []*models.CallExpress {
	panic("unimplemented")
}

// Content implements models.Context.
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

// Functions implements models.Context.
func (f *FileContext) Functions() []*models.Function {
	panic("unimplemented")
}

// LinesInfo implements models.Context.
func (f *FileContext) LinesInfo() models.LinesInfo {
	panic("unimplemented")
}

var _ models.Context = (*FileContext)(nil)
