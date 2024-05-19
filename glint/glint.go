/*
A high-performance static analysis tool that is agnostic to the compiler or runtime environment
of the language.

1 初始化程序
2 定义上下文
3 执行调用
4 处理成功或失败的结果
*/

package glint

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

func init() {
	log.SetLevel(log.DEBUG)
	log.SetOutput(os.Stderr)
	errors.SetWarningPrefixf("%s warning", config.Program)
	errors.SetErrPrefixf("%s error", config.Program)
	errors.SetExitHandler(func(err error) {
		log.Error(err)
	})
}

// setup completion of pre-checks and initialization of the environment.
// * Checking version compatibility.
// * Filtering out excluded models.
// * Setting warnings to disable.
func setup(conf *config.Config) error {

	// Checking version compatibility and check
	if err := config.Validate(conf); err != nil {
		return err
	}

	// clean exclude model
	cleanModels(conf)

	// set warning
	if conf.WarningDisable {
		errors.DisableWarning()
	}
	return nil
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

// Lint Creates Glinter and check the given project path.
func Lint(conf *config.Config, project string) error {

	// init environment
	if err := setup(conf); err != nil {
		return err
	}
	log.Info("successfully setup glint run env")

	// create glinter instance
	glinter, err := NewGLinter(conf)
	if err != nil {
		return err
	}
	log.Infof("successfully created glinter: %s", glinter)

	// check project
	return glinter.run(project)
}

type SS struct {
	models        []*Model
	makeCtxFunc   MakeContextFunc
	preHandleFunc PreHandlerFunc
}

// GLinter ...
type GLinter struct {
	conf *config.Config

	// 分配器，为文件制定不同的检查器，每个文件只能有一个检查器。
	maker *ContextMaker

	// 输出对象，检查的缺陷将写入output
	output Outputer

	lints map[string]ModelFuncType

	defaultLint ModelFuncType

	// 需要执行的预先构建
	prehandlers []PreHandlerFunc

	// 上下文树的根节点
	ctx Context
}

// NewGLinter ...
func NewGLinter(conf *config.Config) (*GLinter, error) {

	// create outputer
	outputer, err := CreateOutput(conf.OutputFile, conf.OutputFormat)
	if err != nil {
		return nil, err
	}
	log.Infof("successfully created outputer: %s", outputer)

	// create ctx maker
	ctxMaker, err := NewContextMaker(conf.Languages)
	if err != nil {
		return nil, err
	}

	log.Infof("successfully created context maker: %s", ctxMaker)

	// create glinter
	glinter := &GLinter{
		conf:   conf,
		output: outputer,
		maker:  ctxMaker,
	}
	return glinter, err
}

// lint
func (g *GLinter) run(path string) error {

	defer func() {
		g.output.Close()
		log.Infof("successfully visited context tree")
	}()
	log.Infof("start lint project: %q", path)

	if err := g.loadProject(path); err != nil {
		return err
	}
	log.Infof("successfully parsed context tree: %s", g.ctx)

	g.PreHandle()
	log.Infof("successfully pre handle file tree")
	g.visitLint(tree, g.conf.Concurrecy)

	return nil
}

// loadModels
func (g *GLinter) loadModels() error {

	anonymityModels := []*Model{}
	for _, lang := range g.conf.Languages {
		meta, err := store.getLangMeta(lang)
		if err != nil {
			return err
		}

		// TODO 冲突检查
		if g.prehandlers != nil {
			g.prehandlers = append(g.prehandlers, meta.preHandleFunc)
		}

		// 未指定后缀名的
		if len(lang.Extends) == 0 {
			anonymityModels = append(anonymityModels, meta.models...)
			continue
		}

		lint, err := g.makeLint(meta.models)
		if err != nil {
			return err
		}
		for _, ext := range lang.Extends {
			if _, ok := g.lints[ext]; ok {
				return errors.Newf("conflict %s extend: %q", lang.Name, ext)
			}
			g.lints[ext] = lint
		}

		if len(anonymityModels) != 0 {
			g.defaultLint, err = g.makeLint(anonymityModels)
			if err != nil {
				return err
			}

		}
	}

	return nil
}

func (g *GLinter) makeLint(models []*Model) (ModelFuncType, error) {

	lints := make([]ModelFuncType, 0, len(models))

	for _, model := range models {
		lint, err := model.GenerateModelFunc(model)
		if err != nil {
			return nil, errors.Newf("failed to compile model: %q, err: %s", model.Name, err)
		} else {
			if lint != nil {
				lints = append(lints, lint)
			}
		}
	}

	f := func(ctx Context) {
		log.Debugf("lint: %q", ctx.File())
		defer func() {
			g.output.Write(ctx)
			log.Debugf("successfully lint %q", ctx.File())
		}()
		for index := range lints {
			lints[index](ctx)
		}
	}
	return f, nil

}

func (g *GLinter) createContext(path string) Context {
	ext := filepath.Ext(path)
	var f MakeContextFunc
	ctx := f(path)
	ctx.Check()

	return nil
}

// load ...
func (g *GLinter) loadProject(project string) error {
	isExclude, err := getExclude(g.conf.ExcludeFiles, g.conf.ExcludeDirs)
	if err != nil {
		return err
	}
	emptyCtx.acquire()
	defer emptyCtx.release()
	if err = g.buildCtxTree(emptyCtx, project, isExclude); err != nil {
		return err
	}
	if len(emptyCtx.children) != 1 {
		utils.Bugf("failed to build context tree")
		return errors.Newf("failed to build Context")
	} else {
		g.ctx = emptyCtx.children[0]
	}
	return nil
}

// buildCtxTree ...
func (g *GLinter) buildCtxTree(ctx Context, path string, exclude func(string, bool) bool) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	fname := info.Name()
	if info.IsDir() {
		// exclude directory
		if exclude(fname, false) {
			return nil
		} else {
			subCtx := g.maker.New(path)
			ctx.AddSubContext(subCtx)
			dirs, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for index := range dirs {
				subPath := filepath.Join(path, dirs[index].Name())
				if err = g.buildCtxTree(subCtx, subPath, exclude); err != nil {
					return err
				}
			}
		}
	} else {
		if exclude(fname, true) {
			return nil
		} else {
			subCtx := g.maker.New(fname)
			ctx.AddSubContext(subCtx)
		}
	}
	return nil
}

// PreHandle ...
func (g *GLinter) PreHandle() {
	// 依次调用语言的预处理程序
	g.PreHandle()

}

// VisitLint ...
func (g *GLinter) visitLint() {

	g.ctx.
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
	for i := 0; i < g.conf.Concurrecy; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ctx := range ctxCh {
				// 检查
				ctx.Check(g.output)
			}
		}()
	}
	wg.Wait()
}

func (g *GLinter) String() string {
	languags := make([]string, 0, len(g.conf.Languages))
	for index := range g.conf.Languages {
		languags = append(languags, g.conf.Languages[index].Name)
	}
	return fmt.Sprintf("<GLinter: %s>", strings.Join(languags, ","))
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
