package basic

import (
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/models"
	"github.com/stkali/glint/utils"
)

func init() {

	err := models.InjectModels(utils.Any,
		&FileBasic,
	)
	errors.CheckErr(err)
	log.Infof("successfully injected %s models", utils.CCpp)
}
