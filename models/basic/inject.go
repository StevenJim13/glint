package basic

import (
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"

	"github.com/stkali/glint/glint"
	"github.com/stkali/glint/utils"
)

func init() {

	// A collection of rules for anonymous languages without corresponding syntax parsing function.
	// It can only analyze file content, file name, and other information, but cannot parse syntax
	// trees.
	err := glint.InjectModels(utils.Any,
		&FileBasicModel,
	)

	errors.CheckErr(err)
	log.Infof("successfully injected %s models", utils.Any)
}
