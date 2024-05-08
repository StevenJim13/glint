package c

import (
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/util"
)

func init() {
	if err := models.InjectModels(util.CCpp,
		&SensitiveApi,
		&FileBasic,
	); err != nil {
		panic(err)
	}
	log.Infof("successfully injected %s models", util.CCpp)
}
