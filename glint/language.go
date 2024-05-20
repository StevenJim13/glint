package glint

import (
	"fmt"
	"slices"
	"strings"

	"github.com/stkali/glint/config"
	"github.com/stkali/utility/errors"
	"github.com/stkali/utility/log"
)

type Langauge struct {
	// Name 语言名
	Names []string
	// Extends 语言的扩展名
	Extends []string
	// NewContext 语言上下文的创建函数
	NewConext NewContextType
	// PreHandler 语言的预处理函数
	Prehandler PreHandlerType
	// Models 语言所拥有的检查规则
	Models []*Model
	//
	modelsMapping map[string]*Model
}

func (l *Langauge) String() string {
	return fmt.Sprintf("Language(%s)", l.Name())
}

func (l *Langauge) getModel(name string) *Model {
	if model, ok := l.modelsMapping[name]; ok {
		return model
	} else {
		return nil
	}
}

func (l *Langauge) Name() string {
	if len(l.Names) != 0 {
		return l.Names[0]
	}
	return "??language"
}

func (l *Langauge) ValidateExtends(extends []string) error {
	for _, ext := range extends {
		if !slices.Contains(l.Extends, ext) {
			return errors.Newf("invalid extend suffix name %s", ext)
		}
	}
	return nil
}

var anonymityLang = &Langauge{}

type LangStore struct {
	lm map[string]*Langauge
}

// register 注册的可能是匿名语言
// 通过语言名称获取
//
//	验证扩展名
//
// 通过扩展名获取
//
//	提示用户这是什么语言
func (l *LangStore) register(lang *Langauge) error {

	for _, name := range lang.Names {
		name := strings.ToLower(name)
		log.Info(l, l.lm)
		l.lm[name] = lang
	}
	lang.modelsMapping = make(map[string]*Model, len(lang.Models))
	for _, model := range lang.Models {
		if _, ok := lang.modelsMapping[model.Name]; ok {
			return errors.Newf("conflict %q model: %q", lang.Name(), model.Name)
		} else {
			lang.modelsMapping[model.Name] = model
		}
	}

	return nil
}

func (l *LangStore) queryLanguage(lang *config.Language) (*Langauge, error) {

	name := strings.ToLower(lang.Name)
	language, ok := l.lm[name]
	if ok {
		if err := language.ValidateExtends(lang.Extends); err != nil {
			return nil, errors.Newf("faield to validate language %q extends, err: %s", lang.Name, err)
		}
		return language, nil
	}

	if len(lang.Extends) == 0 {
		return anonymityLang, nil
	}

	for _, l := range l.lm {
		if ok := slices.ContainsFunc(l.Extends, func(s string) bool {
			return slices.Contains(lang.Extends, s)
		}); ok {
			errors.Warningf("语言名为: %s", l.Names[0])
			return l, nil
		}
	}
	return anonymityLang, nil
}

// func (l *LangStore) getLangMeta(lang *config.Language) (*Langauge, []*Model, error) {
// 	language, err := l.queryLanguage(lang)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	models := make([]*Model, 0, len(lang.Models))
// 	for _, conf := range lang.Models {
// 		model, ok := language.modelsMapping[conf.Name]
// 		if !ok {
// 			return nil, nil, errors.Newf("invalid %q language model: %q", language.Name(), conf.Name)
// 		} else {
// 			model.Tags = conf.Tags
// 			model.Options = conf.Options
// 			models = append(models, model)
// 		}
// 	}
// 	return language,
// }

var store = &LangStore{lm: make(map[string]*Langauge)}

func RegisterLangauge(lang *Langauge) error {
	return store.register(lang)
}
