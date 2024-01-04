package i18n_test

import (
	"fmt"
	"testing"

	"github.com/jingyuexing/i18n"
)

func TestI18n(t *testing.T) {
	messages := map[string]any{
		"en": i18n.Message{
			"greeting": i18n.Message{
				"welcome": "Welcome!",
			},
			"chat": i18n.Message{
				"button": "send",
			},
			"hint": i18n.Message{
				"message": "the {name} is an user",
			},
		},
		"zh": i18n.Message{
			"greeting": i18n.Message{
				"welcome": "你好!",
			},
			"hint": i18n.Message{
				"message": "the {name} is an user",
			},
		},
	}

	i18n := i18n.CreateI18n(&i18n.Options{
		Message:        messages,
		Local:          "zh",
		FallbackLocale: "en",
	})

	if i18n.T("greeting.welcome") != "你好!" {
		t.Error("not pass")
	}
	fmt.Printf("%s\n", i18n.T("chat.button"))
	if i18n.T("chat.button") != "send" {
		t.Error("not pass")
	}
	if i18n.T("chat.submit") != "chat.submit" {
		t.Error("not pass")
	}
	if i18n.T("hint.message", map[string]any{
		"name": "Alan",
	}) != "the Alan is an user" {
		t.Error("not pass")
	}
}

func TestI18nStruct(t *testing.T) {
	type Translate struct {
		Greeting struct {
			Welcome string `json:"welcome"`
		} `json:"greeting"`

		Chat struct {
			Button string `json:"button"`
		} `json:"chat"`

		Hint struct {
			Message string `json:"message"`
		} `json:"hint"`
	}
	zh := &Translate{
		Greeting: struct {
			Welcome string "json:\"welcome\""
		}{
			Welcome: "你好",
		},
		Chat: struct {
			Button string "json:\"button\""
		}{
			Button: "发送",
		},
		Hint: struct {
			Message string "json:\"message\""
		}{
			Message: "提示信息",
		},
	}
	en := &Translate{
		Greeting: struct {
			Welcome string "json:\"welcome\""
		}{
			Welcome: "welcome",
		},
		Chat: struct {
			Button string "json:\"button\""
		}{
			Button: "send",
		},
		Hint: struct {
			Message string "json:\"message\""
		}{
			Message: "hint message",
		},
	}
	i18n_ := i18n.CreateI18n(&i18n.Options{
		Message: i18n.Message{
			"zh": zh,
			"en": en,
		},
		Local:          "zh",
		FallbackLocale: "en",
	})
	if i18n_.T("greeting.welcome") != "greeting.welcome" {
		t.Error("not pass")
	}
	i18n_.SetLocale("en")
	if i18n_.T("Greeting.Welcome") != "welcome" {
		t.Error("not pass")
	}
	i18n_.SetLocale("zh")
	if i18n_.T("Greeting.Welcome") != "你好" {
		t.Error("not pass")
	}

	type User struct {
		Name string `i18n:"Message.Name"`
		Age string `i18n:"Message.Age"`
		Address string `i18n:"Message.Address"`
	}
}

func TestI18nStructTag(t *testing.T) {

	transZH := map[string]any {
		"name":"姓名",
		"age":"年龄信息",
		"i18n":"国际化信息",
	}

	transEN := map[string]any {
		"name":"the name information",
		"age":"the age information",
		"i18n":"the i18n information",
	}

	type User  struct {
		Name string `i18n:"name"`
		Age string `i18n:"age"`
		I18n string `i18n:"i18n"`
	}

	i18nS := i18n.CreateI18n(&i18n.Options{
		Message: map[string]any{
			"zh":transZH,
			"en":transEN,
		},
		Local: "zh",
		FallbackLocale: "en",
	})
	user := &User{}
	result1 := i18nS.TS(user,"Name")
	if result1 != "姓名"{
		t.Error("Not Pass")
	}

	result2 := i18nS.TS(user,"Age")
	if result2 != "年龄信息"{
		t.Error("Not Pass")
	}
	i18nS.SetLocale("en")

	result3 := i18nS.TS(user,"Name")
	if result3 != "the name information" {
		t.Error("Not Pass")
	}

	result4 := i18nS.TS(user,"I18n")
	if result4 != "the i18n information" {
		t.Error("Not Pass")
	}
}
