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
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
	"github.com/stkali/utility/tool"
	"golang.org/x/sync/errgroup"
)

func init() {

	level := log.ToLevelWithDefault(os.Getenv("GLINT_LOG_LEVEL"), log.INFO)
	log.SetLevel(level)
	log.SetOutput(os.Stderr)
	errors.SetWarningPrefixf("%s warning", config.Program)
	errors.SetErrPrefixf("%s error", config.Program)
	// errors.SetExitHandler(func(err error) {
	// 	log.Error(err)
	// })
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
	return glinter.run(tool.ToAbsPath(project))
}

// getExclude TODO
func getExclude(excFiles, excDirs []string) (func(path string, file bool) bool, error) {

	if len(excFiles) == 0 && len(excDirs) == 0 {
		return func(path string, file bool) bool { return false }, nil
	}

	veriryFile, err := makeExcludeFunc(excFiles...)
	if err != nil {
		return nil, err
	}

	verifyDir, err := makeExcludeFunc(excDirs...)
	if err != nil {
		return nil, err
	}
	if verifyDir == nil && veriryFile == nil {
		return func(path string, file bool) bool { return false }, nil
	}
	return func(path string, file bool) bool {
		if file && veriryFile != nil {
			return veriryFile(path)
		} else if verifyDir != nil {
			return verifyDir(path)
		}
		return false
	}, nil
}

type VerifyFileFunc func(path string) bool

func makeExcludeFunc(excludes ...string) (VerifyFileFunc, error) {
	var verify VerifyFileFunc
	switch len(excludes) {
	case 0:
	case 1:
		verify = func(path string) bool {
			if ok, err := filepath.Match(excludes[0], path); err != nil {
				panic(err)
			} else {
				return ok
			}
		}

	default:
		verify = func(path string) bool {
			for index := range excludes {
				if ok, err := filepath.Match(excludes[index], path); err != nil {
					panic(err)
				} else if ok {
					return true
				}
			}
			return false
		}
	}
	return verify, nil
}

// GLinter ...
type GLinter struct {
	conf *config.Config

	// // 分配器，为文件制定不同的检查器，每个文件只能有一个检查器。
	// maker *ContextMaker

	// 输出对象，检查的缺陷将写入output
	output Outputer

	ctxNewFuncMapping map[string]func(string) Context

	defaultCtxCreatFunc func(string) Context

	// 需要执行的预先构建
	prehandlers []PreHandlerType

	// 上下文树的根节点
	pkg Packager
}

// NewGLinter ...
func NewGLinter(conf *config.Config) (*GLinter, error) {

	// create outputer
	outputer, err := CreateOutput(conf.OutputFile, conf.OutputFormat)
	if err != nil {
		return nil, err
	}
	log.Infof("successfully created outputer: %s", outputer)

	// create glinter
	glinter := &GLinter{
		conf:              conf,
		output:            outputer,
		ctxNewFuncMapping: make(map[string]func(string) Context),
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

	if err := g.loadLanguages(); err != nil {
		return err
	}

	if err := g.loadProject(path); err != nil {
		return err
	}
	log.Infof("successfully parsed context tree: %s", g.pkg)

	g.PreHandle()
	log.Infof("successfully pre handle file tree")
	return g.visitLint()
}

// loadModels
func (g *GLinter) loadLanguages() error {

	anonymityModels := []*Model{}

	for _, lang := range g.conf.Languages {

		language, err := store.queryLanguage(lang)
		// TODO check models double
		if err != nil {
			return err
		}
		log.Infof("get language: %s", language)

		if language.Prehandler != nil {
			g.prehandlers = append(g.prehandlers, language.Prehandler)
		}

		models := make([]*Model, 0, len(lang.Models))
		for _, modelConf := range lang.Models {
			model := language.getModel(modelConf.Name)
			if model == nil {
				return errors.Newf("unregister language: %s", modelConf)
			}
			model.Tags = modelConf.Tags
			model.Options = modelConf.Options
			models = append(models, model)
		}

		if len(language.Extends) == 0 {
			anonymityModels = append(anonymityModels, models...)
			continue
		}

		check, err := g.makeCheckFunc(models)
		if err != nil {
			return err
		}

		newCxt := func(path string) Context {
			return language.NewConext(path, check)
		}
		log.Infof("register extends: %s ", language.Extends)
		for _, ext := range language.Extends {
			if _, ok := g.ctxNewFuncMapping[ext]; ok {
				return errors.Newf("conflict %s extend: %q", language.Name, ext)
			} else {
				g.ctxNewFuncMapping[ext] = newCxt
			}
		}

		if len(anonymityModels) != 0 {
			newCxt, err := g.makeCheckFunc(anonymityModels)
			if err != nil {
				return err
			}
			g.defaultCtxCreatFunc = func(file string) Context {
				return NewFileContext(file, newCxt)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *GLinter) makeCheckFunc(models []*Model) (CheckFuncType, error) {

	lints := make([]CheckFuncType, 0, len(models))

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

	f := func(ctx Context) (err error) {
		defer func() {
			g.output.Write(ctx)
		}()
		for index := range lints {
			errors.Join(err, lints[index](ctx))
		}
		return err
	}
	return f, nil

}

func (g *GLinter) makeContext(path string) Context {
	ext := filepath.Ext(path)
	if creatFunc, ok := g.ctxNewFuncMapping[ext]; ok {
		return creatFunc(path)
	} else if g.defaultCtxCreatFunc != nil {
		return g.defaultCtxCreatFunc(path)
	}
	return nil
}

// load ...
func (g *GLinter) loadProject(path string) error {
	isExclude, err := getExclude(g.conf.ExcludeFiles, g.conf.ExcludeDirs)
	if err != nil {
		return err
	}
	g.pkg = NewPackage("root")
	if err = g.buildCtxTree(g.pkg, path, isExclude); err != nil {
		return err
	}
	return nil
}

// buildCtxTree ...
func (g *GLinter) buildCtxTree(pkg Packager, path string, exclude func(string, bool) bool) error {
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
			subPackage := NewPackage(path)
			pkg.AddPackage(subPackage)
			dirs, err := os.ReadDir(path)
			if err != nil {
				return err
			}
			for index := range dirs {
				subPath := filepath.Join(path, dirs[index].Name())
				if err = g.buildCtxTree(subPackage, subPath, exclude); err != nil {
					return err
				}
			}
		}
	} else {
		if exclude(fname, true) {
			return nil
		} else {
			subCtx := g.makeContext(path)
			if subCtx != nil {
				pkg.AddContext(subCtx)
			}
		}
	}
	return nil
}

// PreHandle ...
func (g *GLinter) PreHandle() error {
	switch len(g.prehandlers) {
	case 0:
		return nil
	case 1:
		return g.prehandlers[0](g.conf, g.pkg)
	default:
		var ge *errgroup.Group
		for index := range g.prehandlers {
			ge.Go(func() error {
				return g.prehandlers[index](g.conf, g.pkg)
			})

		}
		return ge.Wait()
	}
}

// VisitLint ...
func (g *GLinter) visitLint() error {
	ctxCh := make(chan Context, 1)
	var wg sync.WaitGroup
	var err error
	for i := 0; i < g.conf.Concurrecy; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ctx := range ctxCh {
				err = errors.Join(err, ctx.Check())
			}
		}()
	}
	g.pkg.Walk(func(ctx Context) {
		ctxCh <- ctx
	})
	close(ctxCh)
	wg.Wait()
	return err

}

func (g *GLinter) String() string {
	languags := make([]string, 0, len(g.conf.Languages))
	for index := range g.conf.Languages {
		languags = append(languags, g.conf.Languages[index].Name)
	}
	return fmt.Sprintf("<GLinter: %s>", strings.Join(languags, ","))
}
