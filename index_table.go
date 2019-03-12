package gokira

import (
	"errors"
	"strconv"
)

// indexTable contains the memory address and size of every single folder packed into the main data file.
type indexTable struct {
	entries map[int][]*index
}

// GetIndex looks up a folder index by its archive and folder id. May return an error.
func (indexTable *indexTable) GetIndex(archiveId, folderId int) (*index, error) {
	if archiveId < 0 {
		return nil, errors.New("index out of bounds")
	}

	archive, ok := indexTable.entries[archiveId]
	if !ok {
		return nil, errors.New("index " + strconv.Itoa(folderId) + " does not exist")
	}

	return archive[folderId], nil
}

// newIndexTable constructs a new table of indices where each index contains the memory address
// and the size of each and every folder that is packed into the main file and ordened
// in a collection of archives. May throw an error.
func newIndexTable(bundle *FileBundle) (*indexTable, error) {
	indices := make(map[int][]*index)

	for indexFileId := 0; indexFileId < len(bundle.indexResources); indexFileId++ {
		resource := bundle.indexResources[indexFileId]
		if len(resource) == 0 {
			continue
		}

		archive, readErr := newIndexList(resource)
		if readErr != nil {
			return nil, readErr
		}

		indices[indexFileId] = archive
	}

	manifestArchive, readManifestArchiveErr := newIndexList(bundle.manifestResource)
	if readManifestArchiveErr != nil {
		return nil, readManifestArchiveErr
	}

	indices[releaseManifestIdx] = manifestArchive

	return &indexTable{entries: indices}, nil
}

// newIndexList combines a collection of mappings of folders that are stored
// in the given data block that implicitly represents the archive itself. May throw an error.
func newIndexList(data []byte) ([]*index, error) {
	folderCount := len(data) / indexSize

	var indices []*index

	for folderId := 0; folderId < folderCount; folderId++ {
		entryStartAddr := folderId * indexSize
		entryEndAddr := entryStartAddr + indexSize

		indexBlock := data[entryStartAddr:entryEndAddr]
		index, entryErr := newIndex(indexBlock)
		if entryErr != nil {
			return nil, entryErr
		}

		indices = append(indices, index)
	}

	return indices, nil
}
