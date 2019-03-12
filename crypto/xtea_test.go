package crypto

import (
	"github.com/sinoz/gokira/buffer"
	"testing"
)

var XTEAEncryptedBlock = []byte{210, 206, 60, 145, 145, 183, 102, 21, 210, 206, 60, 145, 145, 183, 102, 21, 210, 206, 60, 145, 145, 183, 102, 21, 210, 206, 60, 145, 145, 183, 102, 21}

func TestDecipherXTEA(t *testing.T) {
	encryptedBlock := make([]byte, len(XTEAEncryptedBlock))
	copy(encryptedBlock, XTEAEncryptedBlock)

	buf := buffer.HeapByteBufferWrap(encryptedBlock)
	_ = DecipherXTEA(buf, [4]int{1, 2, 3, 4})

	for buf.IsReadable() {
		result, _ := buf.ReadByte()
		if result != 1 && result != 23 {
			t.Error("failed to decrypt XTEA encrypted block")
		}
	}
}

func TestEncipherXTEA(t *testing.T) {
	unencryptedBlock := []byte{49, 47, 44, 31, 39, 45, 88, 28, 58, 28, 19, 48, 69, 99, 121, 27, 21, 33, 99, 98, 97, 94, 91, 12, 1, 56, 45, 88, 91, 57, 77, 71}

	// copy over values as enciphering mutates the buffer
	encryptableBlock := make([]byte, len(unencryptedBlock))
	copy(encryptableBlock, unencryptedBlock)

	buf := buffer.HeapByteBufferWrap(encryptableBlock)

	_ = EncipherXTEA(buf, [4]int{1, 2, 3, 4})

	for i := 0; i < len(encryptableBlock); i++ {
		value, _ := buf.ReadByteAt(i)
		if value == unencryptedBlock[i] {
			t.Error("unencrypted value and encrypted value match at same position after enciphering")
		}
	}

	_ = DecipherXTEA(buf, [4]int{1, 2, 3, 4})

	for i := 0; i < len(unencryptedBlock); i++ {
		value, _ := buf.ReadByteAt(i)
		if value != unencryptedBlock[i] {
			t.Error("unencrypted value and decrypted value mismatch at same position after deciphering")
		}
	}
}
