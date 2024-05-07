package utxo

import (
	"fmt"
	"github.com/acexy/golang-toolkit/util/json"
	"testing"
)

func TestBlockcrypher(t *testing.T) {
	result, err := Analyze(NewBlockcrypherPlatform(Bitcoin, "bc1qur5ym67kljnwrqkw75t0qafe3mq9fxfxn07dwf", "localhost:7890"))
	if err != nil {
		return
	}
	fmt.Println(json.ToJson(result))
}
