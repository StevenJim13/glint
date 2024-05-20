package golang

import (
	"github.com/stkali/glint/config"
	"github.com/stkali/glint/glint"
	"github.com/stkali/utility/log"
)

func PreHandle(conf *config.Config, ctx glint.Context) error {
	log.Infof("apply golang pre handle!")
	return nil
}
