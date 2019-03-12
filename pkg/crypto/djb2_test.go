package crypto

import "testing"

func TestDjb2(t *testing.T) {
	const valueA = "hello_world"
	const valueB = "hello world"

	hashA := Djb2(valueA)
	hashB := Djb2(valueB)

	if hashA == hashB {
		t.Error("hash A and B should not be equal")
	}

	hashC := Djb2(valueA)
	if hashC != hashA {
		t.Error("hash A and C should be equal")
	}

	if hashB == hashC {
		t.Error("hash B and C should not be equal")
	}
}
