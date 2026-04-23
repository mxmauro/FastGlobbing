package FastGlobbing

// -----------------------------------------------------------------------------

func findRune(runes []rune, toFind rune) int {
	for idx, r := range runes {
		if r == toFind {
			return idx
		}
	}
	return -1
}

func findLastRune(runes []rune, toFind rune) int {
	for idx := len(runes) - 1; idx >= 0; idx-- {
		if runes[idx] == toFind {
			return idx
		}
	}
	return -1
}

func findLastPathSeparator(runes []rune) int {
	for idx := len(runes) - 1; idx >= 0; idx-- {
		if isPathSeparator(runes[idx]) {
			return idx
		}
	}
	return -1
}
