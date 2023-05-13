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
