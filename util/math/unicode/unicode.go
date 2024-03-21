package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// GOOS=darwin GOARCH=amd64 go build -o ./. ./util/math/unicode/unicode.go
func main() {
	if len(os.Args) != 2 {
		fmt.Println("提供一个字符")
		return
	}
	char, _ := utf8.DecodeRuneInString(os.Args[1])
	encoded := make([]byte, utf8.UTFMax)
	size := utf8.EncodeRune(encoded, char)
	result := fmt.Sprintf("%08b", encoded[:size])
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	unicode := fmt.Sprintf("u+%04x", char)
	fmt.Printf("Unicode: %s , %d\n", unicode, char)
	fmt.Println("Binary code:", result)
}
