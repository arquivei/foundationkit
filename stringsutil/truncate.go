package stringsutil

// Truncate cuts a given string if it's longer than the given size. Else it returns the string as is.
func Truncate(str string, size int) string {
	if size <= 0 {
		return ""
	}

	runes := []rune(str)
	if len(runes) <= size {
		return str
	}

	return string(runes[:size])
}
