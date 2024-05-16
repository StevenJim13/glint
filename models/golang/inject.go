package golang

import (
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/models/basic"
	"github.com/stkali/glint/utils"
)

func init() {
	err := glint.InjectModels(utils.GoLang,
		&basic.SensitiveApiModel,
		&basic.FileBasicModel,
		&AnnotateStyleModel,
	)
	errors.CheckErr(err)
	log.Infof("successfully injected %s models", utils.GoLang)
}
