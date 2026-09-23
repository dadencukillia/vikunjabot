package bot

import (
	"fmt"
	"strings"
)

func wrapApiUrl(token string, endpoint string) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/%s", token, strings.TrimLeft(endpoint, "/"))
}
