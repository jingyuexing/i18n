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
		Languages:      opts.Languages,
	}
	_i18n.AllLanguge()

	return _i18n
}

type Message = map[string]any

func (self *I18n) T(path string, templates ...map[string]any) string {
	local := self.Local
	if self.Local == "" {
		local = self.FallbackLocale
	}
	translate := ""

	var loadingFunction func(lang string, path []string) string
	_path := []string{path}
	if self.Delimiter != "" {
		_path = strings.Split(path, self.Delimiter)
	}
	switch self.Message[local].(type) {
	case map[string]any:
		loadingFunction = self.loadMapTranslate
	default:
		loadingFunction = self.loadStructTranslate
	}

	translate = loadingFunction(local, _path)

	if translate == "" {
		translate = loadingFunction(self.FallbackLocale, _path)
	}

	translate = utils.ToString(translate)
	if translate == "" {
		return path
	}

	if len(templates) != 0 {
		for _, temp := range templates {
			translate = utils.Template(translate, temp, "")
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

func (self *I18n) TS(val any, field string, templates ...map[string]any) string {
	v := reflect.ValueOf(val)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	path := strings.Split(field, self.Delimiter)

	for _, key := range path {
		if !v.IsValid() {
			return ""
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() == reflect.Struct {
			field, found := v.Type().FieldByName(key)
			if !found {
				return ""
			}
			tag := field.Tag.Get("i18n")
			if tag == "" {
				return ""
			}
			return self.T(tag, templates...)
		}
		return ""
	}
	return ""
}

func (self *I18n) loadStructTranslate(lang string, path []string) string {

	value := reflect.ValueOf(self.Message[lang])

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

func (self *I18n) LoadMessage(lang string, message any) {
	self.Message[lang] = message
	self.AllLanguge()
}

func (self *I18n) SetLocale(locale string) {
	if self.CheckLanguage(locale){
		self.Local = locale
	}
}

func (self *I18n) CheckLanguage(lang string) bool {
	_,ok := self.Message[lang];
	return ok
}

func (self *I18n) AllLanguge() []string {

	if len(self.Languages) != 0 {
		return self.Languages
	}

	lang := make([]string, 0)
	for _lang := range self.Message {
		lang = append(lang, _lang)
	}
	self.Languages = lang
	return self.Languages
}
