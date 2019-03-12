package cache

import (
	"encoding/binary"
	"errors"
	"github.com/sinoz/gokira/pkg/buffer"
	"github.com/sinoz/gokira/pkg/compression"
	"github.com/sinoz/gokira/pkg/crypto"
)

const (
	noCompression    = 0
	bzip2Compression = 1
	gzipCompression  = 2
)

type Folder struct {
	CompressionType byte
	Data            []byte
}

func newFolder(data []byte, keySet [4]int) (*Folder, error) {
	if data == nil || len(data) == 0 {
		return nil, errors.New("given folder data is nil")
	}

	if len(data) < 5 {
		return nil, errors.New("reading folder contents requires at least 5 bytes")
	}

	compressionType := data[0]

	folderSize := binary.BigEndian.Uint32(data[1:])
	folderPayload := data[5:]

	isCompressed := compressionType != noCompression

	if keySet[0] != 0 || keySet[1] != 0 || keySet[2] != 0 || keySet[3] != 0 {
		sizeEncryptedBlock := folderSize
		if isCompressed {
			// + 4 because of the decompressed length in the compressed block's header
			sizeEncryptedBlock += 4
		}

		// deciphering will modify the contents.. let's copy it first
		encryptedBlock := make([]byte, sizeEncryptedBlock)
		copy(encryptedBlock, folderPayload[:sizeEncryptedBlock])

		// and now we can safely decipher the block
		decipherErr := crypto.DecipherXTEA(buffer.HeapByteBufferWrap(encryptedBlock), keySet)
		if decipherErr != nil {
			return nil, decipherErr
		}

		// and assign the folder payload with the decrypted block
		folderPayload = encryptedBlock
	}

	if isCompressed {
		decompressedLength := binary.BigEndian.Uint32(folderPayload)

		if decompressedLength < 0 {
			return nil, errors.New("negative decompressed size")
		}

		if decompressedLength >= 20000000 {
			return nil, errors.New("decompressed size larger than allowed")
		}

		decompressedData, decompressErr := decompressFolder(folderPayload[4:folderSize+4], decompressedLength, compressionType)
		if decompressErr != nil {
			return nil, decompressErr
		}

		// TODO version

		return &Folder{CompressionType: compressionType, Data: decompressedData}, nil
	} else {
		// TODO version

		return &Folder{CompressionType: compressionType, Data: folderPayload[:folderSize]}, nil
	}
}

func decompressFolder(compressedData []byte, decompressedLength uint32, compressionType byte) ([]byte, error) {
	switch compressionType {
	case bzip2Compression:
		decompressedData, decompressErr := compression.DecompressBzip2(compressedData)
		if decompressErr != nil {
			return nil, decompressErr
		}

		if len(decompressedData) != int(decompressedLength) {
			return nil, errors.New("mismatch in bzip2 decompression size")
		}

		return decompressedData, nil

	case gzipCompression:
		decompressedData := make([]byte, decompressedLength)

		decompressedData, decompressErr := compression.DecompressGzip(compressedData)
		if decompressErr != nil {
			return nil, decompressErr
		}

		if len(decompressedData) != int(decompressedLength) {
			return nil, errors.New("mismatch in gzip decompression size")
		}

		return decompressedData, nil

	default:
		return nil, errors.New("unsupported compression type")
	}
}

func (folder *Folder) GetPacks(manifest *FolderManifest) ([]*Pack, error) {
	folderSizeInBytes := len(folder.Data)

	amtPacks := len(manifest.PackReferences)
	amtChunks := int(folder.Data[folderSizeInBytes-1])

	packs := make([]*Pack, amtPacks)

	controlInfoOffset := folderSizeInBytes - 1 - amtChunks*amtPacks*4
	controlInfoBlock := folder.Data[controlInfoOffset:]

	controlInfo := buffer.HeapByteBufferWrap(controlInfoBlock)

	chunkSizes := make([][]int, amtChunks)
	for chunk := range chunkSizes {
		chunkSizes[chunk] = make([]int, amtPacks)
	}

	fileSizes := make([]int, amtPacks)

	for chunk := 0; chunk < amtChunks; chunk++ {
		var chunkSize int

		for pack := 0; pack < amtPacks; pack++ {
			delta, _ := controlInfo.ReadUInt32()
			actualDelta := int32(delta)

			chunkSize += int(actualDelta)
			chunkSizes[chunk][pack] = chunkSize

			fileSizes[pack] += chunkSize
		}
	}

	for pack := 0; pack < amtPacks; pack++ {
		packs[pack] = &Pack{}
	}

	var address int

	for chunk := 0; chunk < amtChunks; chunk++ {
		for pack := 0; pack < amtPacks; pack++ {
			chunkSize := chunkSizes[chunk][pack]

			chunkStart := address
			chunkEnd := address + chunkSize
			chunkData := folder.Data[chunkStart:chunkEnd]

			targetPack := packs[pack]
			targetPack.Data = append(targetPack.Data, chunkData...)

			address += chunkSize
		}
	}

	return packs, nil
}
