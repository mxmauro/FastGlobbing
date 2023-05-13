//go:build !windows
// +build !windows

package FastGlobbing

//-----------------------------------------------------------------------------

func (m *GitWildcardMatcher) perOsInit() {
	m.caseSensitive = true
}
