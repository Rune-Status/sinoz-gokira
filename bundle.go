package gokira

import (
	"io/ioutil"
	"strconv"
)

// FileBundle is the bundle of binary resource files.
type FileBundle struct {
	mainResource     []byte
	indexResources   [][]byte
	manifestResource []byte
}

// NewFileBundle constructs a new FileBundle using the given resources.
func NewFileBundle(data []byte, indices [][]byte, manifest []byte) *FileBundle {
	return &FileBundle{
		mainResource:     data,
		indexResources:   indices,
		manifestResource: manifest,
	}
}

// LoadFileBundle attempts to load a specific collection of resource files located in
// in the specified root path. May also return an error.
func LoadFileBundle(rootPath string, indexCount int) (*FileBundle, error) {
	bundle := new(FileBundle)

	mainFilePath := rootPath + "/main_file_cache.dat2"
	mainResource, err := ioutil.ReadFile(mainFilePath)
	if err != nil {
		return nil, err
	}

	bundle.mainResource = mainResource

	for idxId := 0; idxId < indexCount; idxId++ {
		indexFilePath := rootPath + "/main_file_cache.idx" + strconv.Itoa(idxId)
		idxResource, err := ioutil.ReadFile(indexFilePath)
		if err != nil {
			break
		}

		bundle.indexResources = append(bundle.indexResources, idxResource)
	}

	manifestFilePath := rootPath + "/main_file_cache.idx255"
	manifestResource, err := ioutil.ReadFile(manifestFilePath)
	if err != nil {
		return nil, err
	}

	bundle.manifestResource = manifestResource

	return bundle, nil
}
