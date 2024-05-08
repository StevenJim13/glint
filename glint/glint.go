package glint

import (
	"io"
	"os"

	"github.com/stkali/glint/config"
	_ "github.com/stkali/glint/models/c"
	"github.com/stkali/glint/util"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)


type Context interface {
}

func setEnv(conf *config.Config) error {
	if err := conf.Validate(); err != nil {
		return err
	}

	if !util.IsCompatible(config.Version, conf.Version) {
		return errors.Newf("")
	}

	var logWriter io.Writer
	log.SetLevel(conf.LogLevel)
	if conf.LogFile == "" {
		logWriter = os.Stderr
	} else {
		if f, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_APPEND, os.ModePerm); err != nil {
			return errors.Newf("cannot open log file: %q, err: %s", conf.LogFile, err)
		} else {
			logWriter = f
		}
	}
	log.SetOutput(logWriter)
	return nil
}

func Lint(conf *config.Config, project string) error {
	log.Infof("config: %+v", conf)
	setEnv(conf)
	return nil
}
