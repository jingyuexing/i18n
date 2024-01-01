package i18n

import (
	"reflect"
	"strings"

	"github.com/jingyuexing/go-utils"
)

type I18n struct {
	Message        map[string]any
	Local          string
	FallbackLocale string
	Delimiter      string
	Languages      []string
	isMap          bool
}

type Options struct {
	Message        map[string]any
	Local          string
	FallbackLocale string
	Delimiter      string
	Languages      []string
}

func CreateI18n(opts *Options) *I18n {
	d := "."
	if opts.Delimiter != "" {
		d = opts.Delimiter
	}
	_i18n := &I18n{
		Message:        opts.Message,
		Local:          opts.Local,
		FallbackLocale: opts.FallbackLocale,
		Delimiter:      d,
		Languages: opts.Languages,
	}
	_i18n.AllLanguge()

	return _i18n
}

type Message = map[string]any

func (i *I18n) T(path string, templates ...map[string]any) string {
	local := i.Local
	if i.Local == "" {
		local = i.FallbackLocale
	}
	translate := ""

	var loadingFunction func(lang string, path []string) string
	_path := []string{path}
	if i.Delimiter != "" {
		_path = strings.Split(path, i.Delimiter)
	}
	switch i.Message[local].(type) {
	case map[string]any:
		loadingFunction = i.loadMapTranslate
	default:
		loadingFunction = i.loadStructTranslate
	}

	translate = loadingFunction(local, _path)

	if translate == "" {
		translate = loadingFunction(i.FallbackLocale, _path)
	}

	translate = utils.ToString(translate)
	if translate == "" {
		return path
	}

	if len(templates) != 0 {
		for _, temp := range templates {
			translate = utils.Template(translate, temp)
		}
	}
	return translate
}

func (i *I18n) loadMapTranslate(lang string, path []string) string {

	value := i.Message[lang]
	for _, key := range path {
		if value == nil {
			value = ""
			break
		}
		valMap, ok := value.(map[string]any)
		if !ok {
			value = ""
			break
		}
		if val, exists := valMap[key]; exists {
			value = val
		} else {
			value = ""
		}
	}
	return utils.ToString(value)
}

func (i *I18n) loadStructTranslate(lang string, path []string) string {

	value := reflect.ValueOf(i.Message[lang])

	for _, key := range path {
		if !value.IsValid() {
			return ""
		}
		if value.Kind() == reflect.Pointer {
			value = reflect.Indirect(value)
		}
		if value.Kind() == reflect.Struct {
			field := value.FieldByName(key)
			if !field.IsValid() {
				return ""
			}
			value = field
		} else {
			return ""
		}
	}

	if value.IsValid() && value.CanInterface() {
		return utils.ToString(value.Interface())
	}

	return ""
}

func (i *I18n) LoadMessage(lang string, message any) {
	i.Message[lang] = message
	i.AllLanguge()
}

func (i *I18n) SetLocale(locale string) {
	i.Local = locale
}

func (i *I18n) CheckLanguage(lang string) bool {
	for _lang := range i.Message {
		if _lang == lang {
			return true
		}
	}
	return false
}

func (i *I18n) AllLanguge() []string{

	if len(i.Languages) != 0 {
		return i.Languages
	}

	lang := make([]string,0)
	for _lang := range i.Message {
		lang = append(lang, _lang)
	}
	i.Languages = lang
	return i.Languages
}
