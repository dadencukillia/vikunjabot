package texts

import (
	"errors"
	"fmt"
	"strings"
	"vikunjabot"
)

type LocalePack struct {
	code string
	lines map[string]string
}

func LoadLocalePack(localeCode string) (*LocalePack, error) {
	code := strings.ToLower(strings.TrimSpace(localeCode))

	langContent, err := vikunjabot.LocalizationPack.ReadFile("localizations/" + code + ".lang")
	if err != nil {
		return nil, errors.Join(ErrLocaleNotFound, err)
	}

	pack := LocalePack{
		code: code,
		lines: map[string]string{},
	}

	for line := range strings.SplitSeq(string(langContent), "\n") {
		pair := strings.SplitN(line, "=", 2)
		if len(pair) != 2 {
			continue
		}

		key := strings.ToLower(strings.TrimSpace(pair[0]))
		value := strings.TrimSpace(pair[1])

		if len(key) == 0 {
			continue
		}

		pack.lines[key] = value
	}

	return &pack, nil
}

func (a *LocalePack) Get(key string, variant string) (string, error) {
	lowKey := strings.ToLower(key)
	lowVariant := strings.ToLower(variant)

	if val, ok := a.lines[lowKey + ":" + lowVariant]; ok {
		return val, nil
	}

	if val, ok := a.lines[lowKey]; ok {
		return val, nil
	}

	return "", fmt.Errorf("no such key '%s:%s': %v", strings.ToUpper(key), lowVariant, ErrLocaleKeyNotFound)
}

func (a *LocalePack) GetOrDefault(key string, variant string, defaultVal string) string {
	v, err := a.Get(key, variant)
	if err != nil {
		return defaultVal
	}

	return v
}

func (a *LocalePack) GetWithDict(key string, variant string, processReplacement func(key string) (string, bool)) (string, error) {
	value, err := a.Get(key, variant)
	if err != nil {
		return "", err
	}

	return replaceRegex.ReplaceAllStringFunc(value, func(s string) string {
		replacement, ok := processReplacement(s[2:len(s) - 1])
		if !ok {
			return s
		}

		return replacement
	}), nil
}

func (a *LocalePack) GetWithDictOrDefault(key string, variant string, processReplacement func(key string) (string, bool), defaultVal string) string {
	v, err := a.GetWithDict(key, variant, processReplacement)
	if err != nil {
		return defaultVal
	}

	return v
}

func (a *LocalePack) GetWithRepl(key string, variant string, engine *ReplacementsEngine) (string, error) {
	return a.GetWithDict(key, variant, engine.ProcessReplacement)
}

func (a *LocalePack) GetWithReplOrDefault(key string, variant string, engine *ReplacementsEngine, defaultVal string) string {
	return a.GetWithDictOrDefault(key, variant, engine.ProcessReplacement, defaultVal)
}
