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
}
