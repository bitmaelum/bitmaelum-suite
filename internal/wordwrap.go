package internal

import "strings"

// WordWrap wraps sentences to a maximum length of limit.
func WordWrap(s string, limit int) string {
	if strings.TrimSpace(s) == "" {
		return s
	}

	words := strings.Fields(strings.ToUpper(s))

	var result, line string
	for len(words) > 0 {
		if len(line)+len(words[0]) > limit {
			result += strings.TrimSpace(line) + "\n"
			line = ""
		}

		line = line + words[0] + " "
		words = words[1:]
	}
	if line != "" {
		result += strings.TrimSpace(line)
	}

	return result
}
