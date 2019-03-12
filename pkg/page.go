package cache

import (
	"encoding/binary"
	"errors"
	"strconv"
)

const (
	pageHeaderSize  = 8
	pagePayloadSize = 512

	pageSize = pageHeaderSize + pagePayloadSize
)

// page is a fixed-sized block of data consisting of a header and a payload. A page is
// exactly 520 bytes with an 8-bytes header and a 512-bytes payload.
type page struct {
	id uint16

	// position is the position of this page in a series of linked pages
	// so if page A has page B as its tail, page B's position would be 1
	position uint16

	// tail is the next page that continues this page's data contents
	tail uint32

	// the contents of this page
	content []byte
}

// newPage produces a new page from the given data. May return an error.
func newPage(buf []byte) (*page, error) {
	if len(buf) < pageSize {
		return nil, errors.New("a page should consume at least " + strconv.Itoa(pageSize))
	}

	id := binary.BigEndian.Uint16(buf[0:])
	position := binary.BigEndian.Uint16(buf[2:])
	tail := uint32(buf[4])<<16 | uint32(buf[5])<<8 | uint32(buf[6])
	content := buf[8:pageSize]

	return &page{
		id:       id,
		position: position,
		tail:     tail,
		content:  content,
	}, nil
}
