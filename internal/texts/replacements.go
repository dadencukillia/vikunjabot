package texts

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"vikunjabot/internal/utils"
)

var scopeRegex = regexp.MustCompile("^[A-Z_]+$")
var replaceRegex = regexp.MustCompile(`\^{[A-Z_:]+}`)

type ReplacementsEngine struct {
	sources map[string]any
	handlers map[string]func(key string) (string, bool)
	directs map[string]string
}

func NewReplacementsEngine() *ReplacementsEngine {
	return &ReplacementsEngine{
		sources: map[string]any{},
		handlers: map[string]func(key string) (string, bool){},
		directs: map[string]string{},
	}
}

func (a *ReplacementsEngine) RegisterSource(scope string, source any) error {
	if !scopeRegex.MatchString(scope) {
		return ErrWrongSourceScope
	}

	a.sources[scope] = source

	return nil
}

func (a *ReplacementsEngine) RegisterHandler(scope string, handler func(key string) (string, bool)) error {
	if !scopeRegex.MatchString(scope) {
		return ErrWrongSourceScope
	}

	a.handlers[scope] = handler

	return nil
}

func (a *ReplacementsEngine) RegisterDirect(key string, value string) {
	a.directs[key] = value
}

func (a *ReplacementsEngine) ProcessReplacement(key string) (string, bool) {
	if direct, ok := a.directs[key]; ok {
		return direct, true
	}

	delim := strings.SplitN(key, ":", 2)
	if len(delim) != 2 {
		return "", false
	}

	scope := delim[0]
	varPath := delim[1]

	if handler, ok := a.handlers[scope]; ok {
		return handler(varPath)
	}

	if source, ok := a.sources[scope]; ok {
		sourceRef := reflect.ValueOf(source)

		for field, value := range sourceRef.Fields() {
			if field.Tag.Get("replkey") == varPath {
				return fmt.Sprint(value.Interface()), true
			}
		}

		snakeCase, err := utils.SnakeToPascalCase(varPath)
		if err != nil {
			return "", false
		}
		expectedMethodName := "Repl" + snakeCase

		for method, value := range sourceRef.Methods() {
			if method.Name == expectedMethodName {
				returned := value.Call([]reflect.Value{})
				if len(returned) != 1 {
					return "", false
				}

				return fmt.Sprint(returned[0].Interface()), true
			}
		}
	}

	return "", false
}

func (a *ReplacementsEngine) ProcessString(text string) string {
	return replaceRegex.ReplaceAllStringFunc(text, func(s string) string {
		replacement, ok := a.ProcessReplacement(s[2:len(s) - 1])
		if !ok {
			return s
		}

		return replacement
	})
}
