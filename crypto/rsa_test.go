package crypto

import (
	"testing"
)

func TestCryptRSA(t *testing.T) {
	originalValue := "Hello World"
	originalValueBytes := []byte(originalValue)

	keyPair, generateErr := GenerateRSAKeyPair(1024)
	if generateErr != nil {
		t.Fatal(generateErr)
	}

	encryptedBytes := CryptRSA(originalValueBytes, keyPair.GetPublicModulus(), keyPair.GetPublicExponent())
	decryptedBytes := CryptRSA(encryptedBytes, keyPair.GetPrivateModulus(), keyPair.GetPrivateExponent())

	decryptedValue := string(decryptedBytes)
	if decryptedValue != originalValue {
		t.Errorf("decrypted value did not match %v\n", originalValue)
	}
}
