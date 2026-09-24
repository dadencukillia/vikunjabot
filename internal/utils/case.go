package utils

import "strings"

func SnakeToPascalCase(name string) (string, error) {
	builder := strings.Builder{}

	for part := range strings.SplitSeq(strings.ToLower(name), "_") {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}

		if _, err := builder.WriteString(strings.ToUpper(string(runes[0]))); err != nil {
			return "", err
		}
		if _, err := builder.WriteString(string(runes[1:])); err != nil {
			return "", err
		}
	}

	return builder.String(), nil
}
