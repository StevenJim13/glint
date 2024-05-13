package golang

import (
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/models/basic"
	"github.com/stkali/glint/utils"
)

func init() {
	err := models.InjectModels(utils.GoLang,
		&basic.SensitiveApi,
		&basic.FileBasic,
	)
	errors.CheckErr(err)
	log.Infof("successfully injected %s models", utils.GoLang)
}
