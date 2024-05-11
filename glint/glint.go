package glint

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/models"
	_ "github.com/stkali/glint/models/c"
	_ "github.com/stkali/glint/models/golang"
	_ "github.com/stkali/glint/models/python"
	"github.com/stkali/glint/parser"
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
	// init environment
	if err := setEnv(conf); err != nil {
		return err
	}
	// 清理 exclude 规则
	cleanModels(conf)
	// 生成规则集
	manager, err := models.NewModelManager(conf.Languages)
	if err != nil {
		return errors.Newf("failed to create ModelManager, err: %s", err)
	}

	linter, err := NewLinter(manager)
	if err != nil {
		return err
	}
	return linter.Lint(project)
}

// Linter ...
type Linter struct {
	models []*models.Model
	parser parser.Parser
	ctxCh chan parser.Context
}

// Lint 
func (l *Linter) Lint(project string) error {

	// 解析文件
	if err := l.parser.Parse(project, <-l.ctxCh); err != nil {
		return errors.Newf("failed to parse project: %q, err: %s", project)
	}
	return nil
}

func skipPath(path string) bool {
	return strings.HasPrefix(path, ".git")
}

func NewLinter(manager *models.ModelManager) (*Linter, error) {
	
	linter := &Linter{
		ctxCh: make(chan parser.Context, 1),
		parser: &parser.LangParser{},
	}
	
	return linter, nil
}

// cleanModels 清除那些无用的数据源
func cleanModels(conf *config.Config) {

	existTag := func(tags []string) bool {
		for index := range tags {
			if exists(conf.ExcludeTags, tags[index]) {
				return true
			}
		}
		return false
	}

	for _, lang := range conf.Languages {
		real := 0
		for index := range lang.Models {
			model := lang.Models[index]
			if !exists(conf.ExcludeNames, model.Name) && !existTag(model.Tags) {
				if index != real {
					lang.Models[real] = model
				}
				real += 1
			}
		}
		lang.Models = lang.Models[:real]
	}
}

func exists(s []string, v string) bool {
	for index := range s {
		if s[index] == v {
			return true
		}
	}
	return false
}
