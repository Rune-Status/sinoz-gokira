package cache

import "errors"

// Archive is an aggregate of folders of packs of assets.
type Archive struct {
	Id      int
	storage *Cache
}

// newArchive constructs a new Archive.
func newArchive(id int, storage *Cache) *Archive {
	return &Archive{Id: id, storage: storage}
}

// GetFolder produces a Folder of pages. May return an error.
func (archive *Archive) GetFolder(id int, keySet [4]int) (*Folder, error) {
	pages, folderPagesErr := archive.GetFolderPages(id)
	if folderPagesErr != nil {
		return nil, folderPagesErr
	}

	return newFolder(pages, keySet)
}

// GetFolderPages collects a set of raw pages that together make up the requested folder.
// May throw an error.
func (archive *Archive) GetFolderPages(folderId int) ([]byte, error) {
	folderMapping, getError := archive.storage.mappings.GetIndex(archive.Id, folderId)
	if getError != nil {
		return nil, getError
	}

	offset := folderMapping.address
	remaining := int(folderMapping.size)

	var pageId int
	var pageContents []byte

	for remaining > 0 {
		pageData := archive.storage.bundle.mainResource[offset:]
		page, pageError := newPage(pageData)
		if pageError != nil {
			return nil, pageError
		}

		if int(page.position) != pageId {
			return nil, errors.New("page index mismatch")
		}

		pageAddition := page.content
		if remaining <= pagePayloadSize {
			pageAddition = pageAddition[:remaining]
		}

		pageContents = append(pageContents, pageAddition...)

		offset = page.tail * pageSize
		remaining -= pagePayloadSize

		pageId++
	}

	return pageContents, nil
}
