package utils

import (
	"fmt"
	"strings"
)

func UrlConcat(host string, path... any) string {
	return strings.TrimRight(host, "/") + "/" + strings.TrimLeft(fmt.Sprint(path...), "/")
}
