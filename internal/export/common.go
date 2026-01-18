package export

import (
	"strings"
	"unicode"
)

func splitRendered(rendered, separator string) (string, string, bool) {
	clean := stripANSI(rendered)
	if separator == "" {
		return "", "", false
	}

	parts := strings.SplitN(clean, separator, 2)
	if len(parts) != 2 {
		return "", "", false
	}

	label := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if label == "" && value == "" {
		return "", "", false
	}

	return label, value, true
}

func sanitizeText(input string) string {
	var b strings.Builder
	b.Grow(len(input))
	for _, r := range input {
		if r == '\n' || r == '\t' || r == '\r' {
			b.WriteRune(r)
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
