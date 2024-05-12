package models

import (
	"sync"

	"github.com/stkali/glint/config"
	"github.com/stkali/glint/parser"
	"github.com/stkali/glint/utils"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
	// "github.com/stkali/utility/log"
)

var manager = sync.Map{}

// InjectModels TODO
func InjectModels(lang utils.Language, models ...*Model) error {
	var lm *langManager
	if v, ok := manager.Load(lang); !ok {
		lm = &langManager{ms: make(map[string]*Model, len(models))}
		if other, ok := manager.LoadOrStore(lang, lm); ok {
			lm = other.(*langManager)
		}
	} else {
		lm = v.(*langManager)
	}
	return lm.adds(models...)
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
	Name      string
	Tags      []string
	Options   map[string]any
	ModelFunc ModelFuncType
}

type ModelFuncType func(model *Model, ctx parser.Context)

type langManager struct {
	sync.Mutex
	ms map[string]*Model
}

// adds
func (l *langManager) adds(models ...*Model) error {

	l.Lock()
	defer l.Unlock()
	for index := range models {
		if err := l.add(models[index]); err != nil {
			return err
		}
	}
	return nil
}

// add
func (l *langManager) add(model *Model) error {
	if _, ok := l.ms[model.Name]; ok {
		return errors.Newf("conflict model: %q", model.Name)
	}
	l.ms[model.Name] = model
	return nil
}

// LoadModels
func LoadModels(lang *config.Language) ([]*Model, error) {

	language := utils.ToLanguage(lang.Name)
	if language == utils.Unknown {
		return nil, errors.Newf("unsupport language: %q", lang)
	}
	value, ok := manager.Load(language)
	if !ok {
		return nil, errors.Newf("unsupport language: %q", lang)
	}
	langManager := value.(*langManager)
	modelList := make([]*Model, 0, len(lang.Models))
	for _, modelConf := range lang.Models {
		model, ok := langManager.ms[modelConf.Name]

		if !ok {
			return nil, errors.Newf("invalid %q language model: %q", language, modelConf.Name)
		} else {
			log.Infof("model: %s", model)
			log.Infof("modelConf: %s", modelConf.Options)
			model.Tags = modelConf.Tags
			model.Options = modelConf.Options
			modelList = append(modelList, model)
			log.Infof("model: %s", model)
		}
	}
	return modelList, nil
}
