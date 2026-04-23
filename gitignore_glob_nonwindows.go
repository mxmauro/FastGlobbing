//go:build !windows
// +build !windows

package FastGlobbing

import "os"

//-----------------------------------------------------------------------------

func (m *GitWildcardMatcher) normalizePattern() {}

func isPathSeparator(r rune) bool {
	return r == os.PathSeparator
}

func normalizeRune(r rune) rune {
	return r
}
