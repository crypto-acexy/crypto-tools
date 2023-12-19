package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// GOOS=darwin GOARCH=amd64 go build  -o ./. ./tools/math/hex/h2d.go
func main() {
	args := os.Args
	if len(args) < 1 {
		fmt.Println("提供一个hex值")
		return
	}
	hex := args[1]

	decimal, flag := new(big.Int).SetString(strings.TrimPrefix(strings.TrimPrefix(hex, "0x"), "0X"), 16)
	if !flag {
		fmt.Println("无效的hex值", decimal)
		return
	}
	if len(args) >= 3 {
		scale, err := strconv.ParseUint(args[2], 10, 32)
		if err != nil {
			fmt.Println("无效的精度值")
			return
		}
		hex2decimal(hex, uint(scale))
	} else {
		hex2decimal(hex, 0)
	}
}

func hex2decimal(hex string, scale uint) {
	decimal, flag := new(big.Int).SetString(strings.TrimPrefix(strings.TrimPrefix(hex, "0x"), "0X"), 16)
	if !flag {
		fmt.Println("无效的hex值", decimal)
		return
	}
	result := decimal.String()
	if scale > 0 {
		for i := uint(len(result)); i < scale+1; i++ {
			result = "0" + result
		}
		result = result[:len(result)-int(scale)] + "." + result[len(result)-int(scale):]
	}
	fmt.Println(result)
}
