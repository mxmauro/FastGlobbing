package FastGlobbing

type GitWildcardMatcher struct {
	pattern        []rune
	patternLen     int
	patternHasPath bool
	directoryOnly  bool
}

// -----------------------------------------------------------------------------

// NewGitWildcardMatcher creates a new gitignore-style glob pattern matcher.
func NewGitWildcardMatcher(pattern string) GitWildcardMatcher {
	m := GitWildcardMatcher{
		pattern: []rune(pattern),
	}
	m.normalizePattern()
	if len(m.pattern) > 0 && m.pattern[len(m.pattern)-1] == '/' {
		m.directoryOnly = true
	}
	m.patternLen = len(m.pattern)
	m.patternHasPath = findRune(m.pattern, '/') >= 0
	return m
}

// Test returns true when text matches the gitignore-style glob pattern.
func (m *GitWildcardMatcher) Test(text string) bool {
	textRunes := []rune(text)

	if m.match(textRunes) {
		return true
	}

	if m.directoryOnly && (len(textRunes) == 0 || !isPathSeparator(textRunes[len(textRunes)-1])) {
		textRunes = append(textRunes, '/')
		return m.match(textRunes)
	}

	return false
}

func (m *GitWildcardMatcher) match(textRunes []rune) bool {
	textLen := len(textRunes)

	if m.patternLen == 1 && m.pattern[0] == '*' {
		return true // Speed-up
	}

	textOfs := 0
	patternOfs := 0

	// If the pattern starts with /, anchor it at the start of the path.
	if m.patternLen > 0 && m.pattern[0] == '/' {
		for textOfs+1 < textLen && textRunes[textOfs] == '.' && isPathSeparator(textRunes[textOfs+1]) {
			textOfs += 2
		}
		if textOfs < textLen && isPathSeparator(textRunes[textOfs]) {
			textOfs += 1
		}
		patternOfs = 1
	} else if !m.patternHasPath {
		sepIndex := findLastPathSeparator(textRunes)
		if sepIndex >= 0 {
			textOfs = sepIndex + 1
		}
	}

	// Main loop
	textBackup1 := -1
	patternBackup1 := -1
	textBackup2 := -1
	patternBackup2 := -1

	for textOfs < textLen {
		if patternOfs < m.patternLen {
			switch m.pattern[patternOfs] {
			case '*':
				patternOfs += 1
				if patternOfs < m.patternLen && m.pattern[patternOfs] == '*' {
					patternOfs += 1
					if patternOfs >= m.patternLen {
						return true
					}

					// ** followed by a / matches zero or more directories.
					if m.pattern[patternOfs] != '/' {
						return false
					}

					textBackup1 = -1
					patternBackup1 = -1
					textBackup2 = textOfs
					patternBackup2 = patternOfs
					if !isPathSeparator(textRunes[textOfs]) {
						patternOfs++
					}
					continue
				}

				// Trailing * matches everything except /
				textBackup1 = textOfs
				patternBackup1 = patternOfs
				continue

			case '?':
				if isPathSeparator(textRunes[textOfs]) {
					break
				}
				textOfs += 1
				patternOfs += 1
				continue

			case '[':
				if isPathSeparator(textRunes[textOfs]) {
					break
				}

				textRune := normalizeRune(textRunes[textOfs])
				reverse := patternOfs+1 < m.patternLen && (m.pattern[patternOfs+1] == '^' || m.pattern[patternOfs+1] == '!')
				if reverse {
					patternOfs += 1
				}

				matched := false
				lastChr := rune(-1)
				classHasChars := false
				for {
					patternOfs += 1
					if patternOfs >= m.patternLen {
						break
					}

					if m.pattern[patternOfs] == ']' && classHasChars {
						break
					}

					if lastChr >= 0 && m.pattern[patternOfs] == '-' && patternOfs+1 < m.patternLen && m.pattern[patternOfs+1] != ']' {
						patternOfs += 1
						if textRune <= m.pattern[patternOfs] && textRune >= lastChr {
							matched = true
						}
					} else if textRune == m.pattern[patternOfs] {
						matched = true
					}

					lastChr = m.pattern[patternOfs]
					classHasChars = true
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
				if m.pattern[patternOfs] != normalizeRune(textRunes[textOfs]) && !(m.pattern[patternOfs] == '/' && isPathSeparator(textRunes[textOfs])) {
					break
				}

				textOfs += 1
				patternOfs += 1
				continue
			}
		}

		if patternBackup1 >= 0 && !isPathSeparator(textRunes[textBackup1]) {
			textBackup1 += 1
			textOfs = textBackup1
			patternOfs = patternBackup1
			continue
		}
		if patternBackup2 >= 0 {
			textBackup2 += 1
			textOfs = textBackup2
			patternOfs = patternBackup2
			continue
		}
		return false
	}

	// Ignore trailing stars
	for patternOfs < m.patternLen && m.pattern[patternOfs] == '*' {
		patternOfs += 1
	}

	// At end of text means success if nothing else is left to match.
	return patternOfs >= m.patternLen
}
