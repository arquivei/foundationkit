package stringsutil

// Truncate cuts a given string if it's longer than the given size. Else it returns the string as is.
func Truncate(str string, size int) string {
	if len(str) <= size {
		return str
	}

	runes := []rune(str)
	return string(runes[:size])
}
