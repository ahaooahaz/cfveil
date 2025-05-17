package python

import (
	"regexp"
	"strings"
)

func parseModule(line string) (m []string, ok bool) {
	pattern := regexp.MustCompile(`^\s*(?:from\s+([\s\S]+?)\s+import\s+([\s\S]+?)(?:\s+as\b|$)|import\s+([\s\S]+?)(?:\s+as\b|$))`)

	matches := pattern.FindStringSubmatch(line)
	if len(matches) == 0 {
		ok = false
		return
	}
	ok = true

	behindFrom := matches[1]
	behindImport := matches[2]
	if behindFrom == "" {
		behindImport = matches[3]
	}
	if strings.TrimSpace(behindFrom) != "" {
		m = append(m, strings.TrimSpace(behindFrom))
	}
	behindImport = strings.TrimSpace(behindImport)
	behindImport = strings.Trim(behindImport, "()")
	for p := range strings.SplitSeq(behindImport, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			m = append(m, p)
		}
	}
	return
}
