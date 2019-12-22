package gokira

import (
	"fmt"
	"hash/crc32"

	"github.com/sinoz/gokira/bytes"
)

const polynomial = 0xEDB88320

// ReleaseManifest contains metadata about every archive in a storage.
type ReleaseManifest struct {
	Versions  []uint32
	Checksums []uint32
}

// ArchiveManifest contains metadata about an archive.
type ArchiveManifest struct {
	Id               int
	Format           uint8
	Version          uint32
	Directive        uint8
	FolderReferences []*FolderManifest
}

// FolderManifest contains metadata about a folder in an archive.
type FolderManifest struct {
	Index          int
	Id             int
	LabelHash      uint32
	Version        uint32
	Checksum       uint32
	PackReferences []*PackManifest
}

// PackManifest contains metadata about a pack in a folder.
type PackManifest struct {
	// TODO
}

// newReleaseManifest constructs a new GetReleaseManifest that contains information about every archive in the given Cache.
func newReleaseManifest(cache *Cache) (*ReleaseManifest, error) {
	archiveCount := cache.ArchiveCount()

	release := &ReleaseManifest{
		Versions:  make([]uint32, archiveCount),
		Checksums: make([]uint32, archiveCount),
	}

	for archiveId := 0; archiveId < archiveCount; archiveId++ {
		archive, getManifestErr := cache.GetArchiveManifest(archiveId)
		if getManifestErr != nil {
			return nil, getManifestErr
		}

		pages, getPagesErr := cache.GetFolderPages(255, archiveId)
		if getPagesErr != nil {
			return nil, getPagesErr
		}

		crcTable := crc32.MakeTable(polynomial)
		crcValue := crc32.Checksum(pages, crcTable)

		release.Checksums[archiveId] = crcValue
		release.Versions[archiveId] = archive.Version
	}

	return release, nil
}

// newArchiveManifest constructs a new ArchiveManifest from the given data. May return an error.
func newArchiveManifest(id int, data []byte) (*ArchiveManifest, error) {
	itr := bytes.StringWrap(data).Iterator()

	manifest := new(ArchiveManifest)
	manifest.Id = id

	manifest.Format, _ = itr.ReadByte()
	if manifest.Format < 5 || manifest.Format > 7 {
		return nil, fmt.Errorf("format out of bounds (5-7) but is %v\n", manifest.Format)
	}

	if manifest.Format >= 6 {
		manifest.Version, _ = itr.ReadUInt32()
	}

	manifest.Directive, _ = itr.ReadByte()

	folderCount, _ := itr.ReadUInt16()
	folderIds := make([]int, folderCount)

	accumulator := 0
	lastId := -1

	// read the id of each and every referenced folder in the archive's manifest.
	// this is due to a sudden padding in between some folders where an id is skipped
	for i := 0; i < len(folderIds); i++ {
		idDelta, _ := itr.ReadUInt16()

		accumulator += int(idDelta)
		folderIds[i] = accumulator

		if folderIds[i] > lastId {
			lastId = folderIds[i]
		}
	}

	lastId++

	// and finally allocate folder manifests for each id we've read so we
	// can use this collection to easily read the rest of the data for each folder
	manifest.FolderReferences = make([]*FolderManifest, lastId)
	for index, folderId := range folderIds {
		manifest.FolderReferences[folderId] = &FolderManifest{
			Id:    folderId,
			Index: index,
		}
	}

	// check if the manifest has label hashes enlisted for each folder
	if manifest.containsLabels() {
		// and if so, we read each label hash
		for _, folder := range manifest.FolderReferences {
			if folder != nil {
				folder.LabelHash, _ = itr.ReadUInt32()
			}
		}
	}

	// read the crc checksum of each folder
	for _, folder := range manifest.FolderReferences {
		if folder != nil {
			folder.Checksum, _ = itr.ReadUInt32()
		}
	}

	// read the versions of each folder
	for _, folder := range manifest.FolderReferences {
		if folder != nil {
			folder.Version, _ = itr.ReadUInt32()
		}
	}

	// TODO
	for _, folder := range manifest.FolderReferences {
		if folder != nil {
			packCount, _ := itr.ReadUInt16()
			folder.PackReferences = make([]*PackManifest, packCount)
		}
	}

	return manifest, nil
}

// containsNames returns whether this manifest contains the DJB2 name hashes.
func (manifest *ArchiveManifest) containsLabels() bool {
	return manifest.Directive != 0
}
