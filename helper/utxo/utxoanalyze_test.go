package utxo

import (
	"testing"
)

func TestBlockcrypher(t *testing.T) {
	err := Analyze(NewBlockcrypherPlatform(Bitcoin, "bc1qnfesnxke9ekf0zs8p75gd2qsj6kqvvh0s6j73v", "localhost:7890"))
	if err != nil {
		return
	}
}

func TestMempool(t *testing.T) {
	err := Analyze(NewMempoolPlatform(Bitcoin, "bc1qur5ym67kljnwrqkw75t0qafe3mq9fxfxn07dwf", "localhost:7890"))
	if err != nil {
		return
	}
}
