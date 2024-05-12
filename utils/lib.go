package utils

func MatchNewChar(char string) int {
	switch char {
	case "\\r":
		return 0
	case "\\n":
		return 1
	case "\\r\\n":
		return 2
	default:
		return 100
	}
}
