package c

import (
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/models/basic"
	"github.com/stkali/glint/utils"
)

func init() {
	if err := models.InjectModels(utils.CCpp,
		&basic.SensitiveApi,
		&basic.FileBasic,
	); err != nil {
		panic(err)
	}
	log.Debugf("successfully injected %s models", utils.CCpp)
}
