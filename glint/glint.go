package glint

import (
	"io"
	"os"
	"strconv"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/models"
	_ "github.com/stkali/glint/models/c"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

type Context interface {
}

// setEnv ...
func setEnv(conf *config.Config) error {

	// validte config
	if err := config.Validate(conf); err != nil {
		return err
	}

	// set log
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

// Lint ...
func Lint(conf *config.Config, project string) error {
	if err := setEnv(conf); err != nil {
		return err
	}
	log.Infof("config: %+v", conf)

	
	if len(conf.ExcludeNames) != 0 {

	}

	exclude := func(model config.Model) bool {
		if model.Name 
	}

	// 加载检查模型
	for _, lang := range conf.Languages {
		if _, err := models.LoadModels(lang); err != nil {
			return err
		} else {

		}
		conf.ExcludeNames

	}
	// 逐个执行模型

	//

	return nil
}
