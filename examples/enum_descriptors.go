package main

import (
	"encoding/json"
	"errors"
	"github.com/sinoz/gokira"
	"github.com/sinoz/gokira/buffer"
	"log"
	"strconv"
)

const EnumConfigId = 8

type EnumDescriptor struct {
	Id            uint32              `yaml:"id"`
	KeyType       string              `yaml:"key_type"`
	ValueType     string              `yaml:"value_type"`
	DefaultString string              `yaml:"default_str"`
	DefaultInt    uint32              `yaml:"default_int"`
	Parameters    map[int]interface{} `yaml:"params"`
}

func main() {
	assetCache, err := gokira.LoadCache("cache/", 21)
	if err != nil {
		log.Fatal(err)
	}

	descriptors, err := getEnumDescriptors(assetCache)
	descriptor := descriptors[1131]

	enumAsJson, _ := json.Marshal(descriptor.Parameters)
	println(string(enumAsJson)) // prints interface mappings
}

func getEnumDescriptors(cache *gokira.Cache) ([]*EnumDescriptor, error) {
	archiveManifest, getManifestErr := cache.GetArchiveManifest(2)
	if getManifestErr != nil {
		return nil, getManifestErr
	}

	folder, folderPageErr := cache.GetUnencryptedFolder(2, EnumConfigId)
	if folderPageErr != nil {
		return nil, folderPageErr
	}

	targetFolderManifest := archiveManifest.FolderReferences[EnumConfigId]
	packs, getPacksErr := folder.GetPacks(targetFolderManifest)
	if getPacksErr != nil {
		return nil, getPacksErr
	}

	packCount := len(packs)
	descriptors := make([]*EnumDescriptor, packCount)

	for packId := 0; packId < packCount; packId++ {
		descriptors[packId] = &EnumDescriptor{Id: uint32(packId)}

		packData := packs[packId].Data
		packBuffer := buffer.HeapByteBufferWrap(packData)

		decodeError := decodeEnum(packBuffer, descriptors[packId])
		if decodeError != nil {
			return nil, decodeError
		}
	}

	return descriptors, nil
}

func decodeEnum(buf *buffer.HeapByteBuffer, descriptor *EnumDescriptor) error {
	for buf.IsReadable() {
		id, readErr := buf.ReadByte()
		if readErr != nil {
			return readErr
		}

		if id == 0 {
			break
		}

		var err error

		switch id {
		case 1:
			keyTypeValue, err := buf.ReadByte()
			if err != nil {
				return err
			}

			descriptor.KeyType = string(keyTypeValue)

		case 2:
			valueTypeValue, err := buf.ReadByte()
			if err != nil {
				return err
			}

			descriptor.KeyType = string(valueTypeValue)

		case 3:
			descriptor.DefaultString, err = buf.ReadCString()
			if err != nil {
				return err
			}

		case 4:
			descriptor.DefaultInt, err = buf.ReadUInt32()
			if err != nil {
				return err
			}

		case 5:
			paramCount, err := buf.ReadUInt16()
			if err != nil {
				return err
			}

			descriptor.Parameters = make(map[int]interface{}, paramCount)

			for i := 0; i < int(paramCount); i++ {
				key, _ := buf.ReadUInt32()
				value, _ := buf.ReadCString()

				descriptor.Parameters[int(key)] = value
			}

		case 6:
			paramCount, err := buf.ReadUInt16()
			if err != nil {
				return err
			}

			descriptor.Parameters = make(map[int]interface{}, paramCount)

			for i := 0; i < int(paramCount); i++ {
				key, _ := buf.ReadUInt32()
				value, _ := buf.ReadUInt32()

				descriptor.Parameters[int(key)] = int(value)
			}

		default:
			return errors.New("could not find a case for id " + strconv.Itoa(int(id)) + "")
		}
	}

	return nil
}
