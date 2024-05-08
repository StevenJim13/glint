package models

import (
	"sync"

	"github.com/stkali/glint/util"
	"github.com/stkali/utility/errors"
)

var manager = sync.Map{}

// InjectModels TODO
func InjectModels(lang util.Language, models ... *Model) error {
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
func ExportAllModels() map[util.Language][]*Model {
	langs := make(map[util.Language][]*Model)
	manager.Range(func(key, value any) bool {
		lang, lm := key.(util.Language), value.(*langManager)
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
	Name string
	Tags []string
	Options map[string]any
	ModelFunc ModelFuncType
}

type Context interface {
}

type ModelFuncType func(model *Model, ctx Context)

type langManager struct {
	sync.Mutex
	ms map[string]*Model
}

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

func (l *langManager) add(model *Model) error {
	if _, ok := l.ms[model.Name]; ok {
		return errors.Newf("conflict model: %q", model.Name)
	}
	l.ms[model.Name] = model
	return nil
}
