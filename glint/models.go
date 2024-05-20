package glint

import (
	"strings"
	"sync"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

var manager = sync.Map{}

// InjectModels TODO
func InjectModels(lang utils.Language, models ...*Model) error {
	var lm *langManager
	if v, ok := manager.Load(lang); !ok {
		extends, err := utils.Extends(lang)
		if err != nil {
			return err
		}
		lm = &langManager{
			extends: make(map[string]struct{}, len(extends)),
			ms:      make(map[string]*Model, len(models)),
		}
		for index := range extends {
			lm.extends[extends[index]] = struct{}{}
		}
		if other, ok := manager.LoadOrStore(lang, lm); ok {
			lm = other.(*langManager)
		}
	} else {
		lm = v.(*langManager)
	}
	copy(models)
	return lm.adds(models...)
}

// copy 创建
func copy(ms []*Model) {
	for index := range ms {
		model := *ms[index]
		ms[index] = &model
	}
}

// ExportAllModels ...
func ExportAllModels() map[utils.Language][]*Model {
	langs := make(map[utils.Language][]*Model)
	manager.Range(func(key, value any) bool {
		lang, lm := key.(utils.Language), value.(*langManager)
		list := make([]*Model, 0, len(lm.ms))
		for index := range lm.ms {
			list = append(list, lm.ms[index])
		}
		langs[lang] = list
		return true
	})
	return langs
}

type Model struct {
	Name              string
	Tags              []string
	Options           map[string]any
	Inspect           string
	GenerateModelFunc func(model *Model) (CheckFuncType, error)
}

type CheckFuncType func(ctx Context) error

type langManager struct {
	sync.Mutex
	ms      map[string]*Model
	extends map[string]struct{}
}

// adds
func (l *langManager) adds(models ...*Model) error {
	if len(models) == 0 {
		return nil
	}
	l.Lock()
	defer l.Unlock()
	for _, model := range models {
		if _, ok := l.ms[model.Name]; ok {
			return errors.Newf("conflict model: %q", model.Name)
		}
		l.ms[model.Name] = model
	}
	return nil
}

func (l *langManager) ValidateExtends(extends []string) error {
	for _, ext := range extends {
		if _, ok := l.extends[strings.ToLower(ext)]; !ok {
			return errors.Newf("invalid extend suffix name %s", ext)
		}
	}
	return nil
}

func getLangManager(lang utils.Language) (*langManager, error) {
	value, ok := manager.Load(lang)
	if !ok {
		return nil, errors.Newf("unsupport language: %q", lang)
	}
	// validate language extends
	if langManager, ok := value.(*langManager); ok {
		return langManager, nil
	} else {
		return nil, errors.Newf("this is a bug, failed to assert %s is *langManager", value)
	}
}

// LoadModels ...
// 如果是非
func getModels(lang *config.Language) (utils.Language, []*Model, error) {

	language := utils.ToLanguage(lang.Name)
	if language == utils.Any {
		log.Infof("not found language: %q, associated language-independent models", lang.Name)
	}
	langManager, err := getLangManager(language)
	if err != nil {
		return utils.Unknown, nil, err
	}

	// validate language extends
	if err := langManager.ValidateExtends(lang.Extends); err != nil {
		return utils.Unknown, nil, errors.Newf("faield to validate language %q extends, err: %s", lang.Name, err)
	}

	modelList := make([]*Model, 0, len(lang.Models))
	for _, modelConf := range lang.Models {
		model, ok := langManager.ms[modelConf.Name]
		if !ok {
			return utils.Unknown, nil, errors.Newf("invalid %q language model: %q", language, modelConf.Name)
		} else {
			model.Tags = modelConf.Tags
			model.Options = modelConf.Options
			modelList = append(modelList, model)
		}
	}
	return language, modelList, nil
}
