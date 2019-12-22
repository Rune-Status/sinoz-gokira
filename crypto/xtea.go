package crypto

import (
	"encoding/binary"
)

const (
	GoldenRatio int = -1640531527
	Rounds      int = 32
)

func EncipherXTEA(bytes []byte, keySet [4]int) {
	quadCount := len(bytes) / 8
	for quad := 0; quad < quadCount; quad++ {
		var sum = 0

		v0 := binary.BigEndian.Uint32(bytes[quad*8:])
		v1 := binary.BigEndian.Uint32(bytes[quad*8+4:])
		for round := 0; round < Rounds; round++ {
			v0 += (((v1 << 4) ^ (v1 >> 5)) + v1) ^ uint32(sum+keySet[(sum&3)])
			sum += GoldenRatio
			v1 += (((v0 << 4) ^ (v0 >> 5)) + v0) ^ uint32(sum+keySet[(sum>>11)&3])
		}

		binary.BigEndian.PutUint32(bytes[quad*8:], v0)
		binary.BigEndian.PutUint32(bytes[quad*8+4:], v1)
	}
}

func DecipherXTEA(bytes []byte, keySet [4]int) {
	quadCount := len(bytes) / 8
	for quad := 0; quad < quadCount; quad++ {
		sum := GoldenRatio * Rounds

		v0 := binary.BigEndian.Uint32(bytes[quad*8:])
		v1 := binary.BigEndian.Uint32(bytes[quad*8+4:])

		for round := 0; round < Rounds; round++ {
			v1 -= (((v0 << 4) ^ (v0 >> 5)) + v0) ^ uint32(sum+keySet[(sum>>11)&3])
			sum -= GoldenRatio
			v0 -= (((v1 << 4) ^ (v1 >> 5)) + v1) ^ uint32(sum+keySet[sum&3])
		}

		binary.BigEndian.PutUint32(bytes[quad*8:], v0)
		binary.BigEndian.PutUint32(bytes[quad*8+4:], v1)
	}
}
