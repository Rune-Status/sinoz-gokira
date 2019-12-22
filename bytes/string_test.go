package bytes

import (
	"testing"
)

func TestString_Concat(t *testing.T) {
	b1 := StringOf(1, 2, 3, 4, 5)
	b2 := StringOf(6, 7, 8, 9, 10)
	b3 := b1.Concat(b2)

	if b3.Length() != 10 {
		t.Error("expected byte string length after concatenation to be 10")
	}
}

func TestString_Drop(t *testing.T) {
	b1 := StringOf(1, 2, 3, 4, 5)
	b2 := b1.Drop(3)

	if b2.Length() != 2 {
		t.Error("expected byte string length after dropping to be 2")
	}
}

func TestString_Take(t *testing.T) {
	b1 := StringOf(1, 2, 3, 4, 5)
	b2 := b1.Take(3)

	if b2.Length() != 3 {
		t.Error("expected byte string length after dropping to be 3")
	}
}

func TestIterator_ReadByte(t *testing.T) {
	b := StringOf(1, 2, 3, 4)
	itr := b.Iterator()
	for i := 1; i <= 4; i++ {
		v, _ := itr.ReadByte()
		if int(v) != i {
			t.Errorf("expected byte value at index %v to equal %v", i-1, v)
		}
	}
}

func TestBuilder_WriteByte(t *testing.T) {
	bs := NewDefaultBuilder().
		WriteByte(8).
		Build()

	v, _ := bs.ByteAt(0)
	if v != 8 {
		t.Error("expected byte value at index 0 to equal 8")
	}
}

func TestByteBuilder_Growth(t *testing.T) {
	bb := NewDefaultBuilder()
	for i := 0; i < 128; i++ {
		bb.WriteByte(byte(i))
	}

	if bb.capacity() != 128 {
		t.Error("expected capacity to equal 128")
	}

	bb.WriteByte(byte(2))

	if bb.capacity() != 256 {
		t.Error("expected capacity to equal 256")
	}
}
