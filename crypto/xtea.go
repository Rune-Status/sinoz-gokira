package crypto

import "github.com/sinoz/gokira/buffer"

const (
	GoldenRatio int = -1640531527
	Rounds      int = 32
)

func EncipherXTEA(buf *buffer.HeapByteBuffer, keySet [4]int) error {
	quadCount := buf.Capacity() / 8
	for quad := 0; quad < quadCount; quad++ {
		var sum = 0

		v0, readV0Err := buf.ReadUInt32At(quad * 8)
		if readV0Err != nil {
			return readV0Err
		}

		v1, readV1Err := buf.ReadUInt32At(quad*8 + 4)
		if readV1Err != nil {
			return readV1Err
		}

		for round := 0; round < Rounds; round++ {
			v0 += (((v1 << 4) ^ (v1 >> 5)) + v1) ^ uint32(sum+keySet[(sum&3)])
			sum += GoldenRatio
			v1 += (((v0 << 4) ^ (v0 >> 5)) + v0) ^ uint32(sum+keySet[(sum>>11)&3])
		}

		writeV0Err := buf.OverwriteUInt32(quad*8, v0)
		if writeV0Err != nil {
			return writeV0Err
		}

		writeV1Err := buf.OverwriteUInt32(quad*8+4, v1)
		if writeV1Err != nil {
			return writeV1Err
		}
	}

	return nil
}

func DecipherXTEA(buf *buffer.HeapByteBuffer, keySet [4]int) error {
	quadCount := buf.Capacity() / 8
	for quad := 0; quad < quadCount; quad++ {
		sum := GoldenRatio * Rounds

		v0, readV0Err := buf.ReadUInt32At(quad * 8)
		if readV0Err != nil {
			return readV0Err
		}

		v1, readV1Err := buf.ReadUInt32At(quad*8 + 4)
		if readV1Err != nil {
			return readV1Err
		}

		for round := 0; round < Rounds; round++ {
			v1 -= (((v0 << 4) ^ (v0 >> 5)) + v0) ^ uint32(sum+keySet[(sum>>11)&3])
			sum -= GoldenRatio
			v0 -= (((v1 << 4) ^ (v1 >> 5)) + v1) ^ uint32(sum+keySet[sum&3])
		}

		writeV0Err := buf.OverwriteUInt32(quad*8, v0)
		if writeV0Err != nil {
			return writeV0Err
		}

		writeV1Err := buf.OverwriteUInt32(quad*8+4, v1)
		if writeV1Err != nil {
			return writeV1Err
		}
	}

	return nil
}
