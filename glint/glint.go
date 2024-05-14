package glint

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

// setEnv ...
func setEnv(conf *config.Config) error {

	// validte config
	if err := config.Validate(conf); err != nil {
		return err
	}

	// set warning
	if conf.WarningDisable {
		errors.DisableWarning()
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
	linter, err := NewGLinter(conf)
	if err != nil {
		return err
	}
	log.Infof("successfully create linter: %s", linter)
	return linter.Lint(project)
}

// Linter ...
type GLinter struct {
	conf   *config.Config
	models map[string]*Linter
	output Outputer
}

func NewGLinter(conf *config.Config) (*GLinter, error) {

	// create outputer
	outpuer, err := CreateOutput(conf.ResultFile, conf.ResultFormat)
	if err != nil {
		return nil, err
	}

	// create glinter
	glinter := &GLinter{
		conf:   conf,
		output: outpuer,
		models: make(map[string]*Linter, len(conf.Languages)),
	}

	// build linter(s) mapping
	for _, lang := range conf.Languages {
		linter := &Linter{
			Lang: utils.ToLanguage(lang.Name),
		}
		modelList, err := LoadModels(lang)
		if err != nil {
			return nil, err
		}
		linter.LintFunc = makeLintModels(outpuer, modelList...)
		for i := range lang.Extends {
			if _, ok := glinter.models[lang.Extends[i]]; ok {
				return nil, errors.Newf("conflict extends: %q", lang.Extends[i])
			} else {
				glinter.models[lang.Extends[i]] = linter
			}
		}
	}
	return glinter, err
}

// Lint
func (l *GLinter) Lint(project string) error {

	log.Infof("start lint project: %q", project)
	tree := NewFileTree(project)
	err := tree.Parse(l.conf.ExcludeFiles, l.conf.ExcludeDirs, l.DispatchLinter)
	if err != nil {
		return err
	}
	log.Infof("successfully parsed filetree: %s", tree)
	PreHandleProject(tree)
	log.Infof("successfully pre handle file tree")
	VisitLint(tree, l.conf.Concurrecy)
	log.Infof("visit file tree")
	return nil
}

func (l *GLinter) String() string {
	return "<Linter>"
}

func (l *GLinter) DispatchLinter(file string) *Linter {
	ext := filepath.Ext(file)
	if lint, ok := l.models[ext]; ok {
		return lint
	}
	return nil
}

func PreHandleProject(tree *FileTree) {

}

// VisitLint ...
func VisitLint(tree *FileTree, concurrecy int) {

	ctxCh := make(chan Context, 1)
	// 遍历文件树
	go func() {
		tree.Walk(func(node *FileNode) error {
			if node.Linter != nil {
				ctxCh <- node
			}
			return nil
		})
		defer close(ctxCh)
	}()
	var wg sync.WaitGroup
	for i := 0; i < concurrecy; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ctx := range ctxCh {
				// 检查
				ctx.Lint(ctx)
			}
		}()
	}
	wg.Wait()

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

func makeLintModels(output Outputer, models ...*Model) LintModels {

	//

	return func(ctx Context) {
		file := ctx.File()
		log.Debugf("handle file: %q", file)
		defer func() {
			output.Write(ctx)
		}()
		for index := range models {
			model := models[index]
			models[index].ModelFunc(model, ctx)
		}
	}
}
