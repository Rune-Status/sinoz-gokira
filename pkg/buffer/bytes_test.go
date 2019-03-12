package buffer

import (
	"strings"
	"testing"
)

func TestHeapBufferWrap(t *testing.T) {
	buffer := HeapByteBufferWrap([]byte{1, 2, 3, 4})
	if len(buffer.bytes) != 4 {
		t.Errorf("buffer internal array expected to have a length of 4 but was %v instead\n", len(buffer.bytes))
	}

	if buffer.writerIndex != 4 {
		t.Errorf("writer index expected to be 4 but was %v instead\n", buffer.writerIndex)
	}
}

func TestNewUnpooledHeapBuffer(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	if len(buffer.bytes) != 16 {
		t.Errorf("buffer internal array expected to have a length of 16 but was %v instead\n", len(buffer.bytes))
	}

	if buffer.writerIndex > 0 {
		t.Errorf("writer index expected to be 0 but was %v instead\n", buffer.writerIndex)
	}
}

func TestHeapBuffer_ReadByte(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = 32
	buffer.writerIndex++

	value, _ := buffer.ReadByte()
	if value != 32 {
		t.Errorf("value did not equal 32 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUInt16(t *testing.T) {
	myValue := 1024

	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte(myValue >> 8)
	buffer.bytes[1] = byte(myValue)
	buffer.writerIndex += 2

	value, _ := buffer.ReadUInt16()
	if value != uint16(myValue) {
		t.Errorf("value did not equal 1024 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUInt24(t *testing.T) {
	myValue := 84000

	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte(myValue >> 16)
	buffer.bytes[1] = byte(myValue >> 8)
	buffer.bytes[2] = byte(myValue)
	buffer.writerIndex += 3

	value, _ := buffer.ReadUInt24()
	if value != uint32(myValue) {
		t.Errorf("value did not equal 84000 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUInt32(t *testing.T) {
	myValue := 1 << 26

	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte(myValue >> 24)
	buffer.bytes[1] = byte(myValue >> 16)
	buffer.bytes[2] = byte(myValue >> 8)
	buffer.bytes[3] = byte(myValue)
	buffer.writerIndex += 4

	value, _ := buffer.ReadUInt32()
	if value != uint32(myValue) {
		t.Errorf("value did not equal 1 << 26 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUInt48(t *testing.T) {
	myValue := 1 << 41

	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte(myValue >> 40)
	buffer.bytes[1] = byte(myValue >> 32)
	buffer.bytes[2] = byte(myValue >> 24)
	buffer.bytes[3] = byte(myValue >> 16)
	buffer.bytes[4] = byte(myValue >> 8)
	buffer.bytes[5] = byte(myValue)
	buffer.writerIndex += 6

	value, _ := buffer.ReadUInt48()
	if value != uint64(myValue) {
		t.Errorf("value did not equal 1 << 41 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUInt64(t *testing.T) {
	myValue := 1 << 52

	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte(myValue >> 56)
	buffer.bytes[1] = byte(myValue >> 48)
	buffer.bytes[2] = byte(myValue >> 40)
	buffer.bytes[3] = byte(myValue >> 32)
	buffer.bytes[4] = byte(myValue >> 24)
	buffer.bytes[5] = byte(myValue >> 16)
	buffer.bytes[6] = byte(myValue >> 8)
	buffer.bytes[7] = byte(myValue)
	buffer.writerIndex += 8

	value, _ := buffer.ReadUInt64()
	if value != uint64(myValue) {
		t.Errorf("value did not equal 1 << 52 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUVarInt16(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte(42)
	buffer.writerIndex++

	value, _ := buffer.ReadUVarInt16()
	if value != 42 {
		t.Errorf("expected value to equal 42 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadUVarInt16_2(t *testing.T) {
	myValue := 1024

	buffer := NewHeapByteBuffer(16)
	buffer.bytes[0] = byte((myValue + 32768) >> 8)
	buffer.bytes[1] = byte(myValue + 32768)
	buffer.writerIndex += 2

	value, _ := buffer.ReadUVarInt16()
	if value != 1024 {
		t.Errorf("expected value to equal 1024 but was %v instead", value)
	}
}

func TestHeapBuffer_ReadCString(t *testing.T) {
	myValue := "Hello world"

	buffer := NewHeapByteBuffer(16)

	i := 0
	for _, character := range myValue {
		buffer.bytes[i] = byte(character)
		buffer.writerIndex++
		i++
	}

	buffer.bytes[i] = 0

	value, _ := buffer.ReadCString()
	if value != myValue {
		t.Errorf("value did not equal 'Hello world' but was %v instead", value)
	}
}

func TestHeapBuffer_WriteByte(t *testing.T) {
	expectedValue := 8

	buffer := NewHeapByteBuffer(16)
	buffer.WriteByte(byte(expectedValue))

	value := buffer.bytes[0]
	if value != byte(expectedValue) {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, value)
	}
}

func TestHeapBuffer_WriteInt16(t *testing.T) {
	expectedValue := 768

	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt16(int16(expectedValue))

	value := int16(buffer.bytes[0])<<8 | int16(buffer.bytes[1])
	if value != int16(expectedValue) {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, value)
	}
}

func TestHeapBuffer_WriteInt24(t *testing.T) {
	expectedValue := 1 << 18

	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt24(int32(expectedValue))

	value := int32(buffer.bytes[0])<<16 | int32(buffer.bytes[1])<<8 | int32(buffer.bytes[2])
	if value != int32(expectedValue) {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, value)
	}
}

func TestHeapBuffer_WriteInt32(t *testing.T) {
	expectedValue := 1 << 27

	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt32(int32(expectedValue))

	value := int32(buffer.bytes[0])<<24 | int32(buffer.bytes[1])<<16 | int32(buffer.bytes[2])<<8 | int32(buffer.bytes[3])
	if value != int32(expectedValue) {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, value)
	}
}

func TestHeapBuffer_WriteInt48(t *testing.T) {
	expectedValue := 1 << 40

	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt48(int64(expectedValue))

	value := int64(buffer.bytes[0])<<40 | int64(buffer.bytes[1])<<32 | int64(buffer.bytes[2])<<24 | int64(buffer.bytes[3])<<16 | int64(buffer.bytes[4])<<8 | int64(buffer.bytes[5])
	if value != int64(expectedValue) {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, value)
	}
}

func TestHeapBuffer_WriteInt64(t *testing.T) {
	expectedValue := 1 << 58

	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt64(int64(expectedValue))

	value := int64(buffer.bytes[0])<<56 | int64(buffer.bytes[1])<<48 | int64(buffer.bytes[2])<<40 | int64(buffer.bytes[3])<<32 | int64(buffer.bytes[4])<<24 | int64(buffer.bytes[5])<<16 | int64(buffer.bytes[6])<<8 | int64(buffer.bytes[7])
	if value != int64(expectedValue) {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, value)
	}
}

func TestHeapBuffer_WriteCString(t *testing.T) {
	expectedValue := "hello man"

	buffer := NewHeapByteBuffer(16)
	buffer.WriteCString(expectedValue)

	var bldr strings.Builder

	for buffer.IsReadable() {
		characterValue, readErr := buffer.ReadByte()
		if readErr != nil {
			t.Error(readErr)
		}

		if characterValue == 0 {
			break
		}

		bldr.WriteByte(characterValue)
	}

	result := bldr.String()
	if result != expectedValue {
		t.Errorf("written value did not equal %v but was %v instead", expectedValue, result)
	}
}

func TestHeapBuffer_VarSizeInt8(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.VarSizeInt8(func() {
		buffer.WriteByte(1)
		buffer.WriteByte(2)
		buffer.WriteByte(3)
	})

	sizeWritten := buffer.bytes[0]
	if sizeWritten != 3 {
		t.Errorf("expected size written to equal 3 but was %v instead", sizeWritten)
	}
}

func TestHeapBuffer_VarSizeInt8_2(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt32(0)

	buffer.VarSizeInt8(func() {
		buffer.WriteByte(1)
		buffer.WriteByte(2)
		buffer.WriteByte(3)
	})

	sizeWritten := buffer.bytes[4]
	if sizeWritten != 3 {
		t.Errorf("expected size written to equal 3 but was %v instead", sizeWritten)
	}
}

func TestHeapBuffer_VarSizeInt16(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.VarSizeInt16(func() {
		for i := 0; i < 512; i++ {
			buffer.WriteInt16(0)
		}
	})

	sizeWritten := int16(buffer.bytes[0])<<8 | int16(buffer.bytes[1])
	if sizeWritten != 1024 {
		t.Errorf("expected size written to equal 1024 but was %v instead", sizeWritten)
	}
}

func TestHeapBuffer_VarSizeInt16_2(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.WriteInt32(0)

	buffer.VarSizeInt16(func() {
		for i := 0; i < 512; i++ {
			buffer.WriteInt16(0)
		}
	})

	sizeWritten := int16(buffer.bytes[4])<<8 | int16(buffer.bytes[5])
	if sizeWritten != 1024 {
		t.Errorf("expected size written to equal 1024 but was %v instead", sizeWritten)
	}
}

func TestHeapBuffer_CanWrite(t *testing.T) {
	buffer := NewHeapByteBuffer(16)

	if !buffer.CanWrite(1) {
		t.Error("writer index is at 0 so should be able to write")
	}

	buffer.writerIndex = 12

	if buffer.CanWrite(6) {
		t.Error("writer index is at 12 so should be able to only write exactly 4 bytes")
	}
}

func TestHeapBuffer_IsWritable(t *testing.T) {
	buffer := NewHeapByteBuffer(16)

	if !buffer.IsWritable() {
		t.Error("writer index is at 0 so should be able to write")
	}

	buffer.writerIndex = 16

	if buffer.IsWritable() {
		t.Error("writer index is at 15 so should not be able to write")
	}
}

func TestHeapBuffer_IsReadable(t *testing.T) {
	buffer := NewHeapByteBuffer(16)

	if buffer.IsReadable() {
		t.Error("writer index is at 0 so should not be able to read")
	}

	buffer.writerIndex = 8

	if !buffer.IsReadable() {
		t.Error("writer index is at 8 so should be able to read exactly 8 bytes")
	}
}

func TestHeapBuffer_CanRead(t *testing.T) {
	buffer := NewHeapByteBuffer(16)

	if buffer.CanRead(1) {
		t.Error("writer index is at 0 so should not be able to read")
	}

	buffer.writerIndex = 8

	if !buffer.CanRead(8) {
		t.Error("writer index is at 8 so should be able to read exactly 8 bytes")
	}
}

func TestHeapBuffer_EnsureWritable(t *testing.T) {
	buffer := NewHeapByteBuffer(16)
	buffer.EnsureWritable(32)

	newCapacity := buffer.Capacity()
	if newCapacity != 32 {
		t.Errorf("expected capacity to equal 64 but was %v instead", newCapacity)
	}

	buffer.EnsureWritable(1024)

	newCapacity = buffer.Capacity()
	if newCapacity != 1024 {
		t.Errorf("expected capacity to equal 1024 but was %v instead", newCapacity)
	}

	// pretend we have written 512 bytes in the buffer
	buffer.writerIndex = 512

	// then ask for a potential growth of 1024 bytes
	buffer.EnsureWritable(1024)

	// buffer should double itself to 2048
	newCapacity = buffer.Capacity()
	if newCapacity != 2048 {
		t.Errorf("expected capacity to equal 2048 but was %v instead", newCapacity)
	}
}
