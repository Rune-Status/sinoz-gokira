package crypto

// Djb2 computes a hash value of the given string input using
// Dan Bernstein's DJB2 algorithm.
func Djb2(value string) int {
	var hash int

	for _, character := range value {
		hash = int(character) + ((hash << 5) - hash)
	}

	return hash
}
