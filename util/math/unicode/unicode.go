package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

// GOOS=darwin GOARCH=amd64 go build -o ./. ./tools/math/hex/h2d.go
func main() {
	if len(os.Args) != 2 {
		fmt.Println("提供一个字符")
		return
	}

	char := []rune(os.Args[1])[0]
	encoded := make([]byte, utf8.UTFMax)
	size := utf8.EncodeRune(encoded, char)
	result := fmt.Sprintf("%08b", encoded[:size])
	result = strings.ReplaceAll(result, "[", "")
	result = strings.ReplaceAll(result, "]", "")
	fmt.Println(result)
}
