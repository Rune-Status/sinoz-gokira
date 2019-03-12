package cache

import (
	"errors"
	"log"
)

const (
	releaseManifestIdx = 255
)

// Cache is a file store that can serve information found within the contents of the FileBundle.
type Cache struct {
	bundle   *FileBundle
	mappings *indexTable
	archives map[int]*Archive
}

// LoadCache loads a FileBundle from the specified path and wraps it into an instance of a Cache.
func LoadCache(path string, indexCount int) (*Cache, error) {
	fileBundle, err := LoadFileBundle(path, indexCount)
	if err != nil {
		log.Fatal(err)
	}

	return NewCache(fileBundle)
}

// NewCache constructs a new file store for the given file bundle. May return an error.
func NewCache(bundle *FileBundle) (*Cache, error) {
	mappings, loadMappingsError := newIndexTable(bundle)
	if loadMappingsError != nil {
		return nil, loadMappingsError
	}

	archives := make(map[int]*Archive)
	archiveCount := len(mappings.entries)

	storage := &Cache{bundle: bundle, mappings: mappings, archives: archives}

	for archiveId := 0; archiveId < archiveCount; archiveId++ {
		archives[archiveId] = newArchive(archiveId, storage)
	}

	archives[releaseManifestIdx] = newArchive(releaseManifestIdx, storage)

	return storage, nil
}

func (cache *Cache) GetReleaseManifest() (*ReleaseManifest, error) {
	return newReleaseManifest(cache)
}

func (cache *Cache) GetArchive(id int) (*Archive, error) {
	archive, ok := cache.archives[id]
	if !ok {
		return nil, errors.New("specified archive does not exist")
	}

	return archive, nil
}

func (cache *Cache) GetUnencryptedFolder(archive, folderId int) (*Folder, error) {
	return cache.GetFolder(archive, folderId, [4]int{})
}

func (cache *Cache) GetFolder(archiveId, folderId int, keySet [4]int) (*Folder, error) {
	archive, archiveErr := cache.GetArchive(archiveId)
	if archiveErr != nil {
		return nil, archiveErr
	}

	return archive.GetFolder(folderId, keySet)
}

func (cache *Cache) GetFolderPages(archiveId, folderId int) ([]byte, error) {
	archive, archiveErr := cache.GetArchive(archiveId)
	if archiveErr != nil {
		return nil, archiveErr
	}

	return archive.GetFolderPages(folderId)
}

func (cache *Cache) GetArchiveManifest(archiveId int) (*ArchiveManifest, error) {
	folder, folderPageErr := cache.GetUnencryptedFolder(255, archiveId)
	if folderPageErr != nil {
		return nil, folderPageErr
	}

	return newArchiveManifest(archiveId, folder.Data)
}

func (cache *Cache) GetFolderManifest(archiveId, folderId int) (*FolderManifest, error) {
	archiveManifest, getManifestErr := cache.GetArchiveManifest(archiveId)
	if getManifestErr != nil {
		return nil, getManifestErr
	}

	if len(archiveManifest.FolderReferences) <= folderId {
		return nil, errors.New("specified folder id is out of bounds")
	}

	return archiveManifest.FolderReferences[folderId], nil
}

func (cache *Cache) GetFolderManifestByName(archiveId int, target string) (*FolderManifest, error) {
	archiveManifest, getManifestErr := cache.GetArchiveManifest(archiveId)
	if getManifestErr != nil {
		return nil, getManifestErr
	}

	for _, manifest := range archiveManifest.FolderReferences {
		currentNameHash := manifest.LabelHash
		targetNameHash := uint32(crypto.Djb2(target))

		if currentNameHash == targetNameHash {
			return manifest, nil
		}
	}

	return nil, errors.New("could not find an entry going by the specified name in the specified archive")
}

// ArchiveCount returns the amount of archives this storage has available. This does not
// include the release manifest file as an index, unlike its IndexCount() variant.
func (cache *Cache) ArchiveCount() int {
	return len(cache.bundle.indexResources)
}

// IndexCount returns the amount of index files this store has available. This may include
// the release manifest file as an index.
func (cache *Cache) IndexCount() int {
	return len(cache.mappings.entries)
}
