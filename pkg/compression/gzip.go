package compression

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func DecompressGzip(compressed []byte) ([]byte, error) {
	compressedBuf := bytes.NewBuffer(compressed)
	gzipReader, gzipErr := gzip.NewReader(compressedBuf)
	if gzipErr != nil {
		return nil, gzipErr
	}

	return ioutil.ReadAll(gzipReader)
}
