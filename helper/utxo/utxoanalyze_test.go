package utxo

import (
	"testing"
)

func TestBlockcrypherBitcoin(t *testing.T) {
	err := Analyze(NewBlockcrypherPlatform(Bitcoin, "bc1qlr2enam7fk6gff3ynqp5jp4gpdegsgd2v0ptx6", "localhost:7890"))
	if err != nil {
		return
	}
}

func TestBlockcrypherLitecoin(t *testing.T) {
	err := Analyze(NewBlockcrypherPlatform(Litecoin, "ltc1q8v3usfl6tht0r975kl2hex6pepxwu68ncv3twd", "localhost:7890"))
	if err != nil {
		return
	}
}

func TestMempoolBitcoin(t *testing.T) {
	err := Analyze(NewMempoolPlatform(Bitcoin, "bc1qlr2enam7fk6gff3ynqp5jp4gpdegsgd2v0ptx6", "localhost:7890"))
	if err != nil {
		return
	}
}
