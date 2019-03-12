package gokira

import (
	"errors"
	"strconv"
)

const (
	indexSize = 6
)

// index contains the start address and the size of a single folder.
type index struct {
	address uint32
	size    uint32
}

func newIndex(data []byte) (*index, error) {
	if len(data) < indexSize {
		return nil, errors.New("folder index requires at least " + strconv.Itoa(indexSize) + " bytes")
	}

	size := uint32(data[0])<<16 | uint32(data[1])<<8 | uint32(data[2])
	blockId := uint32(data[3])<<16 | uint32(data[4])<<8 | uint32(data[5])

	address := blockId * pageSize

	return &index{address: address, size: size}, nil
}
