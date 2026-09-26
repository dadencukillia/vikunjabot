package utils

func CutString(text string, maxLength int, trailingDots bool) string {
	if len(text) <= maxLength {
		return text
	}

	cutPos := min(len(text), maxLength)
	if trailingDots {
		cutPos = max(0, cutPos - 3)
	}

	tail := ""
	if trailingDots {
		tail = "..."
	}

	return string([]rune(text)[:cutPos]) + tail
}
