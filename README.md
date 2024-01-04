# i18n




## import
```go
import "github.com/jingyuexing/i18n"
```


## usage

```go
package main

import (
	"fmt"

	i18n "github.com/jingyuexing/i18n"
)

func main(){
	messages := map[string]any{
        "en": i18n.Message{
            "greeting": i18n.Message{
                "welcome": "Welcome!",
            },
            "chat":i18n.Message{
                "button":"send",
            },
            "hint":i18n.Message{
            	"message":"the {name} is an user"
            }
        },
        "zh": i18n.Message{
            "greeting": i18n.Message{
                "welcome": "你好!",
            },
            "hint":i18n.Message{
            	"message":"the {name} is an user"
            }
        },
    }

	i18n := i18n.CreateI18n(&i18n.I18n{
	    Message:        messages,
	    // default language
	    Local:          "zh",
	    // When the key value specified in the Local language is not found,
	    // it will search again in the language specified by FallbackLocale.
	    // If it is also not found, the key passed in will be returned as is.
	    FallbackLocale: "en",
	})

	fmt.Printf("%s\n",i18n.T("greeting.welcome")) // will print "你好!"
	fmt.Printf("%s\n",i18n.T("hint.message"),map[string]any{
		"name":"Alan"
	}) // will print "the Alan is an user"
}
```

- struct tag

```go
package main

import (
	"fmt"

	i18n "github.com/jingyuexing/i18n"
)

type User struct {
	//
	Name string `i18n:"Error.Validate.Name"`
	Age string `i18n:"Error.Validate.Age"`
}

func main(){
	i18n_ := i18n.CreateI18n(&i18n.Options{
		Message: map[string]any{
			// the chinese translate information
			"zh":transZH,
			// the english translate information
			"en":transEN,
		},
		Local: "zh",
		FallbackLocale: "en",
	})

	user := &User{
		Name:"Altan",
		Age:20,
	}
	// validate the fields
	// ...
	i18n_.TS(user,"Name") // The internationalized prompt information for the 'Name' field will be returned
}

```
