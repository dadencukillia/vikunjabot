package utils

import (
	"strings"

	"golang.org/x/net/html"
)

func ExtractHTMLText(htmlPart string) string {
	domDoc := html.NewTokenizer(strings.NewReader(htmlPart))
	content := strings.Builder{}

loopDom:
	for {
		tt := domDoc.Next()
		switch tt {
		case html.ErrorToken:
			break loopDom
		case html.TextToken:
			TxtContent := strings.TrimSpace(html.UnescapeString(string(domDoc.Text())))
			if len(TxtContent) > 0 {
				content.WriteString(TxtContent)
				content.WriteByte('\n')
			}
		}
	}

	return content.String()
}
