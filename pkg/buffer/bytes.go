package buffer

import (
	"errors"
	"io"
	"strings"
)

type HeapByteBuffer struct {
	bytes       []byte
	readerIndex int
	writerIndex int
}

func HeapByteBufferWrap(bytes []byte) *HeapByteBuffer {
	return &HeapByteBuffer{
		bytes:       bytes,
		writerIndex: len(bytes),
	}
}

func NewHeapByteBuffer(capacity int) *HeapByteBuffer {
	return &HeapByteBuffer{bytes: make([]byte, capacity)}
}

func (buffer *HeapByteBuffer) Fill(reader io.Reader, amtRequiredBytes int) error {
	buffer.EnsureWritable(amtRequiredBytes)

	for amtRequiredBytes > 0 {
		amtBytesRead, readError := reader.Read(buffer.bytes[buffer.writerIndex : buffer.writerIndex+amtRequiredBytes])
		if readError != nil {
			return readError
		}

		amtRequiredBytes -= amtBytesRead
		buffer.writerIndex += amtBytesRead
	}

	return nil
}

func (buffer *HeapByteBuffer) ReadByteAt(index int) (byte, error) {
	if index >= buffer.Capacity() {
		return 0, errors.New("index out of bounds")
	}

	return buffer.bytes[index], nil
}

func (buffer *HeapByteBuffer) ReadUInt16At(index int) (uint16, error) {
	if (index + 1) >= buffer.Capacity() {
		return 0, errors.New("index out of bounds")
	}

	return uint16(buffer.bytes[index])<<8 | uint16(buffer.bytes[index+1]), nil
}

func (buffer *HeapByteBuffer) ReadUInt24At(index int) (uint32, error) {
	if (index + 2) >= buffer.Capacity() {
		return 0, errors.New("index out of bounds")
	}

	return uint32(buffer.bytes[index])<<16 | uint32(buffer.bytes[index+1])<<8 | uint32(buffer.bytes[index+2]), nil
}

func (buffer *HeapByteBuffer) ReadUInt32At(index int) (uint32, error) {
	if (index + 3) >= buffer.Capacity() {
		return 0, errors.New("index out of bounds")
	}

	return uint32(buffer.bytes[index])<<24 | uint32(buffer.bytes[index+1])<<16 | uint32(buffer.bytes[index+2])<<8 | uint32(buffer.bytes[index+3]), nil
}

func (buffer *HeapByteBuffer) ReadUInt48At(index int) (uint64, error) {
	if (index + 5) >= buffer.Capacity() {
		return 0, errors.New("index out of bounds")
	}

	return uint64(buffer.bytes[index])<<40 | uint64(buffer.bytes[index+1])<<32 | uint64(buffer.bytes[index+2])<<24 | uint64(buffer.bytes[index+3])<<16 | uint64(buffer.bytes[index+4])<<8 | uint64(buffer.bytes[index+5]), nil
}

func (buffer *HeapByteBuffer) ReadUInt64At(index int) (uint64, error) {
	if (index + 7) >= buffer.Capacity() {
		return 0, errors.New("index out of bounds")
	}

	return uint64(buffer.bytes[index])<<56 | uint64(buffer.bytes[index+1])<<48 | uint64(buffer.bytes[index+2])<<40 | uint64(buffer.bytes[index+3])<<32 | uint64(buffer.bytes[index+4])<<24 | uint64(buffer.bytes[index+5])<<16 | uint64(buffer.bytes[index+6])<<8 | uint64(buffer.bytes[index+7]), nil
}

func (buffer *HeapByteBuffer) ReadByte() (byte, error) {
	if !buffer.CanRead(1) {
		return 0, errors.New("reader index + amount to read out of bounds")
	}

	value := buffer.bytes[buffer.readerIndex]
	buffer.readerIndex++
	return value, nil
}

func (buffer *HeapByteBuffer) ReadBool() (bool, error) {
	result, err := buffer.ReadByte()

	return result == 1, err
}

func (buffer *HeapByteBuffer) ReadUInt16() (uint16, error) {
	if !buffer.CanRead(2) {
		return 0, errors.New("reader index + amount to read out of bounds")
	}

	value := uint16(buffer.bytes[buffer.readerIndex])<<8 | uint16(buffer.bytes[buffer.readerIndex+1])
	buffer.readerIndex += 2
	return value, nil
}

func (buffer *HeapByteBuffer) ReadUInt24() (uint32, error) {
	if !buffer.CanRead(3) {
		return 0, errors.New("reader index + amount to read out of bounds")
	}

	value := uint32(buffer.bytes[buffer.readerIndex])<<16 |
		uint32(buffer.bytes[buffer.readerIndex+1])<<8 |
		uint32(buffer.bytes[buffer.readerIndex+2])
	buffer.readerIndex += 3
	return value, nil
}

func (buffer *HeapByteBuffer) ReadUInt32() (uint32, error) {
	if !buffer.CanRead(4) {
		return 0, errors.New("reader index + amount to read out of bounds")
	}

	value := uint32(buffer.bytes[buffer.readerIndex])<<24 | uint32(buffer.bytes[buffer.readerIndex+1])<<16 | uint32(buffer.bytes[buffer.readerIndex+2])<<8 | uint32(buffer.bytes[buffer.readerIndex+3])
	buffer.readerIndex += 4
	return value, nil
}

func (buffer *HeapByteBuffer) ReadUInt48() (uint64, error) {
	if !buffer.CanRead(6) {
		return 0, errors.New("reader index + amount to read out of bounds")
	}

	value := uint64(buffer.bytes[buffer.readerIndex])<<40 | uint64(buffer.bytes[buffer.readerIndex+1])<<32 | uint64(buffer.bytes[buffer.readerIndex+2])<<24 | uint64(buffer.bytes[buffer.readerIndex+3])<<16 | uint64(buffer.bytes[buffer.readerIndex+4])<<8 | uint64(buffer.bytes[buffer.readerIndex+5])
	buffer.readerIndex += 6
	return value, nil
}

func (buffer *HeapByteBuffer) ReadUInt64() (uint64, error) {
	if !buffer.CanRead(8) {
		return 0, errors.New("reader index + amount to read out of bounds")
	}

	value := uint64(buffer.bytes[buffer.readerIndex])<<56 | uint64(buffer.bytes[buffer.readerIndex+1])<<48 | uint64(buffer.bytes[buffer.readerIndex+2])<<40 | uint64(buffer.bytes[buffer.readerIndex+3])<<32 | uint64(buffer.bytes[buffer.readerIndex+4])<<24 | uint64(buffer.bytes[buffer.readerIndex+5])<<16 | uint64(buffer.bytes[buffer.readerIndex+6])<<8 | uint64(buffer.bytes[buffer.readerIndex+7])
	buffer.readerIndex += 8
	return value, nil
}

func (buffer *HeapByteBuffer) ReadUVarInt32() (uint, error) {
	value, readErr := buffer.ReadByteAt(buffer.readerIndex)
	if readErr != nil {
		return 0, readErr
	}

	if value < 0 {
		value, readErr := buffer.ReadUInt32()

		return uint(value), readErr
	} else {
		value, readErr := buffer.ReadUInt16()

		return uint(value), readErr
	}
}

func (buffer *HeapByteBuffer) ReadBigUVarInt16() (uint, error) {
	var result uint

	var accumulator uint
	var readErr error

	for accumulator, readErr = buffer.ReadUVarInt16(); accumulator == 32767; accumulator, readErr = buffer.ReadUVarInt16() {
		if readErr != nil {
			return 0, readErr
		}

		result += 32767
	}

	result += accumulator
	return result, nil
}

func (buffer *HeapByteBuffer) ReadUVarInt16() (uint, error) {
	value, readErr := buffer.ReadByteAt(buffer.readerIndex)
	if readErr != nil {
		return 0, readErr
	}

	if value < 128 {
		value, readErr := buffer.ReadByte()

		return uint(value), readErr
	} else {
		value, readErr := buffer.ReadUInt16()

		return uint(value) - 32768, readErr
	}
}

func (buffer *HeapByteBuffer) ReadCString() (string, error) {
	var bldr strings.Builder

	for buffer.IsReadable() {
		characterValue, readErr := buffer.ReadByte()
		if readErr != nil {
			return "", readErr
		}

		if characterValue == 0 {
			break
		}

		bldr.WriteByte(characterValue)
	}

	return bldr.String(), nil
}

func (buffer *HeapByteBuffer) ReadDoubleEndedCString() (string, error) {
	terminatorValue, err := buffer.ReadByte()
	if err != nil {
		return "", err
	}

	if terminatorValue != 0 {
		return "", errors.New("not a double ended c-string")
	}

	return buffer.ReadCString()
}

func (buffer *HeapByteBuffer) ReadSlice(size int) []byte {
	slice := buffer.Slice(buffer.readerIndex, buffer.readerIndex+size)
	buffer.readerIndex += size
	return slice
}

func (buffer *HeapByteBuffer) Slice(start, end int) []byte {
	return buffer.bytes[start:end]
}

func (buffer *HeapByteBuffer) WriteByte(value byte) {
	buffer.EnsureWritable(1)

	buffer.bytes[buffer.writerIndex] = value
	buffer.writerIndex++
}

func (buffer *HeapByteBuffer) WriteBool(value bool) {
	if value {
		buffer.WriteByte(1)
	} else {
		buffer.WriteByte(0)
	}
}

func (buffer *HeapByteBuffer) WriteInt16(value int16) {
	buffer.EnsureWritable(2)

	buffer.bytes[buffer.writerIndex] = byte(value >> 8)
	buffer.bytes[buffer.writerIndex+1] = byte(value)

	buffer.writerIndex += 2
}

func (buffer *HeapByteBuffer) WriteInt24(value int32) {
	buffer.EnsureWritable(3)

	buffer.bytes[buffer.writerIndex] = byte(value >> 16)
	buffer.bytes[buffer.writerIndex+1] = byte(value >> 8)
	buffer.bytes[buffer.writerIndex+2] = byte(value)

	buffer.writerIndex += 3
}

func (buffer *HeapByteBuffer) WriteInt32(value int32) {
	buffer.EnsureWritable(4)

	buffer.bytes[buffer.writerIndex] = byte(value >> 24)
	buffer.bytes[buffer.writerIndex+1] = byte(value >> 16)
	buffer.bytes[buffer.writerIndex+2] = byte(value >> 8)
	buffer.bytes[buffer.writerIndex+3] = byte(value)

	buffer.writerIndex += 4
}

func (buffer *HeapByteBuffer) WriteInt48(value int64) {
	buffer.EnsureWritable(6)

	buffer.bytes[buffer.writerIndex] = byte(value >> 40)
	buffer.bytes[buffer.writerIndex+1] = byte(value >> 32)
	buffer.bytes[buffer.writerIndex+2] = byte(value >> 24)
	buffer.bytes[buffer.writerIndex+3] = byte(value >> 16)
	buffer.bytes[buffer.writerIndex+4] = byte(value >> 8)
	buffer.bytes[buffer.writerIndex+5] = byte(value)

	buffer.writerIndex += 6
}

func (buffer *HeapByteBuffer) WriteInt64(value int64) {
	buffer.EnsureWritable(8)

	buffer.bytes[buffer.writerIndex] = byte(value >> 56)
	buffer.bytes[buffer.writerIndex+1] = byte(value >> 48)
	buffer.bytes[buffer.writerIndex+2] = byte(value >> 40)
	buffer.bytes[buffer.writerIndex+3] = byte(value >> 32)
	buffer.bytes[buffer.writerIndex+4] = byte(value >> 24)
	buffer.bytes[buffer.writerIndex+5] = byte(value >> 16)
	buffer.bytes[buffer.writerIndex+6] = byte(value >> 8)
	buffer.bytes[buffer.writerIndex+7] = byte(value)

	buffer.writerIndex += 8
}

func (buffer *HeapByteBuffer) WriteBytes(data []byte) {
	buffer.EnsureWritable(len(data))

	for i := 0; i < len(data); i++ {
		buffer.WriteByte(data[i])
	}
}

func (buffer *HeapByteBuffer) OverwriteByte(index int, value byte) error {
	if index >= buffer.Capacity() {
		return errors.New("index out of bounds")
	}

	buffer.bytes[index] = byte(value)

	return nil
}

func (buffer *HeapByteBuffer) OverwriteUInt16(index int, value uint16) error {
	if index+1 >= buffer.Capacity() {
		return errors.New("index out of bounds")
	}

	buffer.bytes[index] = byte(value >> 8)
	buffer.bytes[index+1] = byte(value)

	return nil
}

func (buffer *HeapByteBuffer) OverwriteUInt24(index int, value uint32) error {
	if index+2 >= buffer.Capacity() {
		return errors.New("index out of bounds")
	}

	buffer.bytes[index] = byte(value >> 16)
	buffer.bytes[index+1] = byte(value >> 8)
	buffer.bytes[index+2] = byte(value)

	return nil
}

func (buffer *HeapByteBuffer) OverwriteUInt32(index int, value uint32) error {
	if index+3 >= buffer.Capacity() {
		return errors.New("index out of bounds")
	}

	buffer.bytes[index] = byte(value >> 24)
	buffer.bytes[index+1] = byte(value >> 16)
	buffer.bytes[index+2] = byte(value >> 8)
	buffer.bytes[index+3] = byte(value)

	return nil
}

func (buffer *HeapByteBuffer) OverwriteUInt48(index int, value uint64) error {
	if index+5 >= buffer.Capacity() {
		return errors.New("index out of bounds")
	}

	buffer.bytes[index] = byte(value >> 40)
	buffer.bytes[index+1] = byte(value >> 32)
	buffer.bytes[index+2] = byte(value >> 24)
	buffer.bytes[index+3] = byte(value >> 16)
	buffer.bytes[index+4] = byte(value >> 8)
	buffer.bytes[index+5] = byte(value)

	return nil
}

func (buffer *HeapByteBuffer) OverwriteUInt64(index int, value uint32) error {
	if index+7 >= buffer.Capacity() {
		return errors.New("index out of bounds")
	}

	buffer.bytes[index] = byte(value >> 56)
	buffer.bytes[index+1] = byte(value >> 48)
	buffer.bytes[index+2] = byte(value >> 40)
	buffer.bytes[index+3] = byte(value >> 32)
	buffer.bytes[index+4] = byte(value >> 24)
	buffer.bytes[index+5] = byte(value >> 16)
	buffer.bytes[index+6] = byte(value >> 8)
	buffer.bytes[index+7] = byte(value)

	return nil
}

func (buffer *HeapByteBuffer) WriteCString(value string) {
	for _, character := range value {
		buffer.WriteByte(byte(character))
	}

	buffer.WriteByte(0)
}

func (buffer *HeapByteBuffer) VarSizeInt8(block func()) {
	offset := buffer.writerIndex

	buffer.WriteByte(0)

	defer func() {
		amtBytesWritten := buffer.writerIndex - offset - 1

		buffer.bytes[offset] = byte(amtBytesWritten)
	}()

	block()
}

func (buffer *HeapByteBuffer) VarSizeInt16(block func()) {
	offset := buffer.writerIndex

	buffer.WriteInt16(0)

	defer func() {
		amtBytesWritten := buffer.writerIndex - offset - 2

		buffer.bytes[offset] = byte(amtBytesWritten >> 8)
		buffer.bytes[offset+1] = byte(amtBytesWritten)
	}()

	block()
}

func (buffer *HeapByteBuffer) EnsureWritable(amount int) {
	for !buffer.CanWrite(amount) {
		trailingBuffer := NewHeapByteBuffer(int(buffer.Capacity()))
		buffer.bytes = append(buffer.bytes, trailingBuffer.bytes...)
	}
}

func (buffer *HeapByteBuffer) ResetReaderIndex() {
	buffer.readerIndex = 0
}

func (buffer *HeapByteBuffer) SetReaderIndex(index int) {
	buffer.readerIndex = index
}

func (buffer *HeapByteBuffer) ResetWriterIndex() {
	buffer.writerIndex = 0
}

func (buffer *HeapByteBuffer) ReaderIndex() int {
	return buffer.readerIndex
}

func (buffer *HeapByteBuffer) SetWriterIndex(index int) {
	buffer.writerIndex = index
}

func (buffer *HeapByteBuffer) WriterIndex() int {
	return buffer.writerIndex
}

func (buffer *HeapByteBuffer) ReadableBytes() int {
	return buffer.writerIndex - buffer.readerIndex
}

func (buffer *HeapByteBuffer) Skip(amount int) {
	buffer.SetReaderIndex(buffer.readerIndex + amount)
}

func (buffer *HeapByteBuffer) Capacity() int {
	return len(buffer.bytes)
}

func (buffer *HeapByteBuffer) IsWritable() bool {
	return buffer.CanWrite(1)
}

func (buffer *HeapByteBuffer) CanWrite(amount int) bool {
	return (buffer.Capacity() - buffer.writerIndex) >= amount
}

func (buffer *HeapByteBuffer) IsReadable() bool {
	return buffer.CanRead(1)
}

func (buffer *HeapByteBuffer) CanRead(amount int) bool {
	return (buffer.writerIndex - buffer.readerIndex) >= amount
}

func (buffer *HeapByteBuffer) BytesWrittenSoFar() []byte {
	return buffer.bytes[:buffer.writerIndex]
}

func (buffer *HeapByteBuffer) Duplicate() *HeapByteBuffer {
	return HeapByteBufferWrap(buffer.bytes)
}

func (buffer *HeapByteBuffer) Clear() {
	buffer.readerIndex = 0
	buffer.writerIndex = 0
}

func (buffer *HeapByteBuffer) String() string {
	return string(buffer.bytes)
}
