package main

import (
	"fmt"
	"strings"
)

func countWords(str string) map[string]int {
	count := make(map[string]int)
	for _, char := range strings.Fields(str) {
		char = strings.ToLower(char)
		lst := []string{}
		for _, c := range char {
			if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
				lst = append(lst, string(c))
			}
		}
		char = strings.Join(lst, "")
		count[strings.ToLower(char)]++
	}
	return count
}

func isPalindrome(str string) bool {
	l, r := 0, len(str)-1
	for l < r {
		if str[l] != str[r] {
			return false
		}
		l++
		r--
	}
	return true
}

func main() {
	fmt.Println(countWords("a A! A1 a1"))
	fmt.Println(isPalindrome("abba"))
}
