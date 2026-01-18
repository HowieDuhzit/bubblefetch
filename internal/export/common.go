package export

import (
	"strings"
	"unicode"
)

func splitRendered(rendered, separator string) (string, string, bool) {
	clean := stripANSI(rendered)
	if separator == "" {
		separator = ""
	}

	if strings.TrimSpace(separator) == "" {
		separator = ""
	}

	candidates := []string{
		separator,
		": ",
		":",
		" | ",
		" → ",
		" - ",
	}

	for _, sep := range candidates {
		if sep == "" {
			continue
		}
		if !strings.Contains(clean, sep) {
			continue
		}

		parts := strings.SplitN(clean, sep, 2)
		if len(parts) != 2 {
			continue
		}

		label := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		label = strings.TrimRight(label, ":")

		if label == "" && value == "" {
			continue
		}

		return label, value, true
	}

	return "", "", false
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
