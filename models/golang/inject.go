package golang

import (
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/models/basic"
	"github.com/stkali/glint/utils"
)

func init() {

	language := &glint.Langauge{
		Names:      []string{"golang", "go"},
		Extends:    []string{".go"},
		NewConext:  NewContext,
		Prehandler: PreHandle,
		Models: []*glint.Model{
			&basic.SensitiveApiModel,
			&basic.FileBasicModel,
			&AnnotateStyleModel,
			&MissAnnotateModel,
			&TestModel,
		},
	}
	err := glint.RegisterLangauge(language)
	errors.CheckErr(err)
	log.Infof("successfully injected %s models", utils.GoLang)
}
