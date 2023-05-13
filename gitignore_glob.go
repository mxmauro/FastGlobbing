package FastGlobbing

import (
	"os"
	"unicode"
)

type GitWildcardMatcher struct {
	patternRunes   []rune
	patternLen     int
	patternHasPath bool
	caseSensitive  bool
}

// -----------------------------------------------------------------------------

// NewGitWildcardMatcher creates a new gitignore-style glob pattern matcher.
func NewGitWildcardMatcher(pattern string) GitWildcardMatcher {
	m := GitWildcardMatcher{
		patternRunes: []rune(pattern),
	}
	m.patternLen = len(m.patternRunes)
	m.patternHasPath = findRune(m.patternRunes, os.PathSeparator) >= 0
	m.perOsInit()
	return m
}

// Test returns a patter if text string matches gitignore-style glob pattern.
func (m *GitWildcardMatcher) Test(text string) bool {
	textRunes := []rune(text)
	textLen := len(textRunes)

	if m.patternLen == 1 && m.patternRunes[0] == '*' {
		return true // speed-up
	}

	textOfs := 0

	// if pattern does not contain a path, skip it
	if !m.patternHasPath {
		sepIndex := findRune(textRunes, os.PathSeparator)
		if sepIndex >= 0 {
			textOfs = sepIndex + 1
		}
	}

	//main loop
	textBackup1 := -1
	patternBackup1 := -1
	textBackup2 := -1
	patternBackup2 := -1
	patternOfs := 0

	for textOfs < textLen {
		if patternOfs < m.patternLen {
			switch m.patternRunes[patternOfs] {
			case '*':
				// match anything except . after /
				patternOfs += 1
				if patternOfs < m.patternLen && m.patternRunes[patternOfs] == '*' {
					// trailing ** match everything after /
					patternOfs += 1
					if patternOfs >= m.patternLen {
						return true
					}

					// ** followed by a / match zero or more directories
					if m.patternRunes[patternOfs] != os.PathSeparator {
						return false
					}

					// new **-loop, discard *-loop
					textBackup1 = -1
					patternBackup1 = -1
					textBackup2 = textOfs
					patternBackup2 = patternOfs
					if textRunes[textOfs] != os.PathSeparator {
						patternOfs++
					}
					continue
				}

				// trailing * matches everything except /
				textBackup1 = textOfs
				patternBackup1 = patternOfs
				continue

			case '?':
				// match any character except /
				if textRunes[textOfs] == os.PathSeparator {
					break
				}
				textOfs += 1
				patternOfs += 1
				continue

			case '[':
				// match any character in [...] except /
				if textRunes[textOfs] == os.PathSeparator {
					break
				}

				// inverted character class
				reverse := patternOfs+1 < m.patternLen && (m.patternRunes[patternOfs+1] == '^' || m.patternRunes[patternOfs+1] == '!')
				if reverse {
					patternOfs += 1
				}

				// match character class
				matched := false
				lastChr := rune(-1)
				for {
					patternOfs += 1
					if patternOfs >= m.patternLen || m.patternRunes[patternOfs] == ']' {
						break
					}

					r := textRunes[textOfs]
					if !m.caseSensitive {
						r = unicode.ToUpper(r)
					}
					if lastChr >= 0 && m.patternRunes[patternOfs] == '-' && patternOfs+1 < m.patternLen && m.patternRunes[patternOfs+1] != ']' {
						patternOfs += 1
						if r <= m.patternRunes[patternOfs] && r >= lastChr {
							matched = true
						}
					} else {
						if r == m.patternRunes[patternOfs] {
							matched = true
						}
					}

					lastChr = m.patternRunes[patternOfs]
				}
				if matched == reverse {
					break
				}
				textOfs += 1
				if patternOfs < m.patternLen {
					patternOfs += 1
				}
				continue

			default:
				// match the current non-NUL character
				r := textRunes[textOfs]
				if !m.caseSensitive {
					r = unicode.ToUpper(r)
				}
				if m.patternRunes[patternOfs] != r && (!(m.patternRunes[patternOfs] == os.PathSeparator && textRunes[textOfs] == os.PathSeparator)) {
					break
				}

				// do not match a . with *, ? [] after /
				textOfs += 1
				patternOfs += 1
				continue
			}
		}

		if patternBackup1 >= 0 && textRunes[textBackup1] != os.PathSeparator {
			// *-loop: backtrack to the last * but do not jump over /
			textBackup1 += 1
			textOfs = textBackup1
			patternOfs = patternBackup1
			continue
		}
		if patternBackup2 >= 0 {
			// **-loop: backtrack to the last **
			textBackup2 += 1
			textOfs = textBackup2
			patternOfs = patternBackup2
			continue
		}
		return false
	}

	//ignore trailing stars
	for patternOfs < m.patternLen && m.patternRunes[patternOfs] == '*' {
		patternOfs += 1
	}

	//at end of text means success if nothing else is left to match
	return patternOfs >= m.patternLen
}
