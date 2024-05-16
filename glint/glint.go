package glint

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
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
	glinter, err := NewGLinter(conf)
	if err != nil {
		return err
	}
	log.Infof("successfully create glinter: %s", glinter)
	return glinter.Lint(project)
}

// Linter ...
type GLinter struct {
	conf       *config.Config
	dispatcher *Dispatcher
	output     Outputer
}

// NewGLinter ...
func NewGLinter(conf *config.Config) (*GLinter, error) {

	// create outputer
	outputer, err := CreateOutput(conf.ResultFile, conf.ResultFormat)
	if err != nil {
		return nil, err
	}
	log.Infof("successfully created outputer: %s", outputer)

	// create dispatcher
	dispatcher, err := NewDispatcher(conf.Languages)
	if err != nil {
		return nil, err
	}
	log.Infof("successfully created dispatcher: %s", dispatcher)

	// create glinter
	glinter := &GLinter{
		conf:       conf,
		output:     outputer,
		dispatcher: dispatcher,
	}

	return glinter, err
}

// Lint
func (g *GLinter) Lint(project string) error {

	defer func() {
		g.output.Flush()
	}()

	log.Infof("start lint project: %q", project)
	tree := NewFileTree(project)
	err := tree.Parse(g.conf.ExcludeFiles, g.conf.ExcludeDirs, g.dispatcher.Dispatch)
	if err != nil {
		return err
	}
	log.Infof("successfully parsed filetree: %s", tree)
	g.PreHandle(tree)
	log.Infof("successfully pre handle file tree")
	g.VisitLint(tree, g.conf.Concurrecy)
	log.Infof("visit file tree")
	return nil
}

func (g *GLinter) PreHandle(tree *FileTree) {
	log.Debugf("pre handle %s", tree)
}

// VisitLint ...
func (g *GLinter) VisitLint(tree *FileTree, concurrecy int) {

	ctxCh := make(chan Context, 1)
	// 遍历文件树
	go func() {
		tree.Walk(func(ctx Context) error {
			ctxCh <- ctx
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
				ctx.Lint(g.output, ctx)
			}
		}()
	}
	wg.Wait()
}

func (l *GLinter) String() string {
	return "<Linter>"
}

type Dispatcher struct {
	langModels      map[string]*Linter
	anonymityModels map[string]*Linter
	defaultLinter   *Linter
}

func NewDispatcher(languages []*config.Language) (*Dispatcher, error) {
	// build linter(s) mapping
	// 未指定语言的可有存在多个Any
	// 但是多个 Any 中 extend 是不能重复的
	// "" 会匹配任意文件

	dispatcher := Dispatcher{
		langModels:      make(map[string]*Linter),
		anonymityModels: make(map[string]*Linter),
	}
	for _, lang := range languages {
		linter, err := NewLinter(lang)
		if err != nil {
			return nil, err
		}
		err = dispatcher.RegisterLinter(linter, lang.Extends...)
		if err != nil {
			return nil, err
		}
	}
	return &dispatcher, nil
}

func (d *Dispatcher) RegisterLinter(linter *Linter, extends ...string) error {
	var models map[string]*Linter
	if linter.Lang == utils.Any {
		// 不指定扩展名时 默认的匹配所有的文件
		if len(extends) == 0 {
			d.defaultLinter = linter
			return nil
		}
		models = d.anonymityModels
	} else {
		if len(extends) == 0 {
			return errors.Error("必须指定扩展名")
		}
		models = d.langModels
	}
	for _, ext := range extends {
		if _, ok := models[ext]; !ok {
			models[ext] = linter
		} else {
			return errors.Newf("conflict %s extend: %q", linter.Lang, ext)
		}
	}
	return nil
}

func (d *Dispatcher) Dispatch(file string) *Linter {

	ext := filepath.Ext(file)
	if lint, ok := d.langModels[ext]; ok {
		return lint
	}

	if len(d.anonymityModels) != 0 {
		if lint, ok := d.anonymityModels[ext]; ok {
			return lint
		}
	}

	return d.defaultLinter
}

type LintFuncType func(Outputer, Context)

type Linter struct {
	Lang     utils.Language
	LintFunc LintFuncType
}

func NewLinter(lang *config.Language) (*Linter, error) {

	language, models, err := getModels(lang)
	if err != nil {
		return nil, err
	}

	modelFuncs := make([]ModelFuncType, 0, len(models)*2)
	for index := range models {
		if modelFunc, err := models[index].GenerateModelFunc(models[index]); err != nil {
			return nil, errors.Newf("failed to compile model: %q, err: %s", models[index].Name, err)
		} else {
			if modelFunc != nil {
				modelFuncs = append(modelFuncs, modelFunc)
			}
		}
	}

	lintFunc := func(outputer Outputer, ctx Context) {
		log.Debugf("lint: %q", ctx.File())
		defer func() {
			outputer.Write(ctx)
			log.Debugf("successfully lint %q", ctx.File())
		}()
		for index := range modelFuncs {
			modelFuncs[index](ctx)
		}
	}

	linter := &Linter{
		Lang:     language,
		LintFunc: lintFunc,
	}
	return linter, err
}

func (l *Linter) String() string {
	return fmt.Sprintf("<Linter: %s %p>", l.Lang, l.LintFunc)
}

// cleanModels 清除那些无用的规则
func cleanModels(conf *config.Config) {
	for _, lang := range conf.Languages {
		real := 0
		for index := range lang.Models {
			model := lang.Models[index]
			if !slices.Contains(conf.ExcludeNames, model.Name) &&
				!slices.ContainsFunc(model.Tags, func(tag string) bool { return slices.Contains(conf.ExcludeTags, tag) }) {
				if index != real {
					lang.Models[real] = model
				}
				real += 1
			}
		}
		lang.Models = lang.Models[:real]
	}
}
