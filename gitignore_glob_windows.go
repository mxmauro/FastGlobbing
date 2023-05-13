package FastGlobbing

import (
	"os"
	"unicode"
)

// -----------------------------------------------------------------------------

func (m *GitWildcardMatcher) perOsInit() {
	for idx := range m.patternRunes {
		if m.patternRunes[idx] != '/' {
			m.patternRunes[idx] = unicode.ToUpper(m.patternRunes[idx])
		} else {
			m.patternRunes[idx] = os.PathSeparator
		}
	}
}
