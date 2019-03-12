package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"math/big"
)

type RSAKeySet struct {
	privateKey *rsa.PrivateKey
	publicKey  rsa.PublicKey
}

func GenerateRSAKeyPair(bitSize int) (*RSAKeySet, error) {
	reader := rand.Reader

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return nil, err
	}

	return &RSAKeySet{privateKey: key, publicKey: key.PublicKey}, nil
}

func (keyPair *RSAKeySet) GetPrivateExponent() *big.Int {
	return keyPair.privateKey.D
}

func (keyPair *RSAKeySet) GetPrivateModulus() *big.Int {
	return keyPair.privateKey.N
}

func (keyPair *RSAKeySet) GetPublicModulus() *big.Int {
	return keyPair.publicKey.N
}

func (keyPair *RSAKeySet) GetPublicExponent() *big.Int {
	return big.NewInt(int64(keyPair.publicKey.E))
}

func CryptRSA(data []byte, modulus, exponent *big.Int) []byte {
	dataInt := &big.Int{}
	dataInt.SetBytes(data)

	resultInt := &big.Int{}
	resultInt.Exp(dataInt, exponent, modulus)

	return resultInt.Bytes()
}
