package terminal

import "fmt"

const (
	Block     string = "█"
	SemiBlock string = "▒"
)

func Abbr(s string) string {
	return fmt.Sprintf("%s...%s", s[:5], s[len(s)-5:])
}

func FirstN(s string, n int) string {
	return fmt.Sprintf("%s...", s[:n])
}

func LastN(s string, n int) string {
	return fmt.Sprintf("%s...", s[len(s)-n:])
}
