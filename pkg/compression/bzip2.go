package compression

import (
	"bytes"
	"compress/bzip2"
	"io/ioutil"
)

var bzip2Header = []byte{
	'B',
	'Z',
	'h',
	'9',
}

func DecompressBzip2(compressed []byte) ([]byte, error) {
	includingHeader := append(bzip2Header, compressed...)

	compressedBuf := bytes.NewBuffer(includingHeader)
	bzipReader := bzip2.NewReader(compressedBuf)

	return ioutil.ReadAll(bzipReader)
}
