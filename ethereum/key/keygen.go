package key

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/acexy/golang-toolkit/math/conversion"
)

func rsaGen(len int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, len)
	if err != nil {
		return nil, nil
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey
}

func rsaKeyToBit() {
	prik, _ := rsaGen(512)
	fmt.Println(prik.Size())
	fmt.Println(prik.N.BitLen())

	binaries := conversion.NewFormBytes(prik.D.Bytes())
	fmt.Println(binaries.To8Bits())
	// c726024de0c791d7cb8f1d951eff458b742dd25d264cb04b435af9d31824bd3f
	fmt.Println(binaries.To2HexString())
}
