package glint

import (
	"io"
	"os"
	"path/filepath"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/models"
	_ "github.com/stkali/glint/models/c"
	_ "github.com/stkali/glint/models/golang"
	_ "github.com/stkali/glint/models/python"
	"github.com/stkali/glint/parser"
	"github.com/stkali/glint/utils"
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

/*
	ext   -> 	parser
				language

*/

// Lint ...
func Lint(conf *config.Config, project string) error {

	// init environment
	if err := setEnv(conf); err != nil {
		return err
	}
	log.Info("successfully set lint env")

	// 清理 exclude 规则
	cleanModels(conf)
	log.Info("successfully clean models")

	// 生成规则集
	// manager, err := models.NewModelManager(conf.Languages)
	// if err != nil {
	// 	return errors.Newf("failed to create ModelManager, err: %s", err)
	// }
	log.Info("create models manager")

	linter, err := NewLinter(conf)
	if err != nil {
		return err
	}
	log.Info("successfully create linter")
	return linter.Lint(project)
}

// Linter ...
type Linter struct {
	conf   *config.Config
	macher *modelsMatcher
}

func PreHandleProject(tree *parser.FileTree) {

}

func VisitLint(tree *parser.FileTree) {
	for _, child := range tree.RootNode().Children {
		if child.Linter != nil && child.Linter.LintFunc != nil {
			log.Info(child.File)
			child.Linter.LintFunc(parser.NewContext(child.File))
		}
	}
}

// Lint
func (l *Linter) Lint(project string) error {
	log.Infof("start lint project: %q", project)
	// 解析文件
	tree := parser.NewFileTree(project)
	err := tree.Parse(l.conf.ExcludeFiles, l.conf.ExcludeDirs, l.macher)
	if err != nil {
		return err
	}
	PreHandleProject(tree)
	VisitLint(tree)

	return nil
}

func (l *Linter) ApplyModels(tree *parser.FileNode) {

}

func NewLinter(conf *config.Config) (*Linter, error) {
	modelsMatcher, err := NewModelsMatcher(conf.Languages...)
	if err != nil {
		return nil, err
	}
	log.Infof("successfully created models matcher: %s", modelsMatcher)
	linter := &Linter{
		conf:   conf,
		macher: modelsMatcher,
	}
	return linter, err
}

// cleanModels 清除那些无用的规则
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

func makeLintModels(models ...*models.Model) parser.LintModels {
	return func(ctx parser.Context) {
		for index := range models {
			model := models[index]
			models[index].ModelFunc(model, ctx)
		}
	}
}

type modelsMatcher struct {
	models map[string]*parser.Linter
}

func NewModelsMatcher(langs ...*config.Language) (*modelsMatcher, error) {

	matcher := &modelsMatcher{
		models: make(map[string]*parser.Linter, len(langs)),
	}

	for index := range langs {
		lang := langs[index]
		linter := &parser.Linter{
			Lang: utils.ToLanguage(lang.Name),
		}
		modelList, err := models.LoadModels(lang)
		if err != nil {
			return nil, err
		}
		linter.LintFunc = makeLintModels(modelList...)
		for i := range lang.Extends {
			if _, ok := matcher.models[lang.Extends[i]]; ok {
				return nil, errors.Newf("conflict extends: %q", lang.Extends[i])
			} else {
				matcher.models[lang.Extends[i]] = linter
			}
		}
	}

	return matcher, nil
}

func (m *modelsMatcher) Match(file string) *parser.Linter {
	ext := filepath.Ext(file)
	if lint, ok := m.models[ext]; ok {
		return lint
	}
	return nil
}
