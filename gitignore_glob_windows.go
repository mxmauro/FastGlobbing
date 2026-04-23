package FastGlobbing

import (
	"unicode"
)

// -----------------------------------------------------------------------------

func (m *GitWildcardMatcher) normalizePattern() {
	for idx := range m.pattern {
		m.pattern[idx] = normalizeRune(m.pattern[idx])
	}
}

func isPathSeparator(r rune) bool {
	return r == '/' || r == '\\'
}

func normalizeRune(r rune) rune {
	if isPathSeparator(r) {
		return '/'
	}
	return unicode.ToUpper(r)
}
