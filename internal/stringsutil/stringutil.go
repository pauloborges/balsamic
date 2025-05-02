package stringsutil

import (
	"strings"
)

func StripMargin(s string) string {
	return StripMarginWithPrefix(s, "|")
}

func StripMarginWithPrefix(s, prefix string) string {
	var lines []string

	for line := range strings.SplitSeq(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, prefix)
		if idx >= 0 {
			line = line[idx+len(prefix):]
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
