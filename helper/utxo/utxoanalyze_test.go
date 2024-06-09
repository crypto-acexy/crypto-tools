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
	err := Analyze(NewBlockcrypherPlatform(Litecoin, "ltc1q7tz5cldmjwcrdgyu2ld38n3qkvc6cy8f4zv7hc", "localhost:7890"))
	if err != nil {
		return
	}
}

func TestBlockcrypherDogecoin(t *testing.T) {
	err := Analyze(NewBlockcrypherPlatform(Dogecoin, "ltc1qtcf8j6rc50hu3arfhz23mj64e3eue0z8ph99va", "localhost:7890"))
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
