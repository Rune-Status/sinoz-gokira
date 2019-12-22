package main

import (
	"encoding/json"
	"errors"

	"github.com/sinoz/gokira"
	"github.com/sinoz/gokira/bytes"

	"log"
	"strconv"
)

const EnumConfigId = 8

type EnumDescriptor struct {
	ID            uint32              `yaml:"id"`
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
	archiveManifest, err := cache.GetArchiveManifest(2)
	if err != nil {
		return nil, err
	}

	folder, err := cache.GetUnencryptedFolder(2, EnumConfigId)
	if err != nil {
		return nil, err
	}

	targetFolderManifest := archiveManifest.FolderReferences[EnumConfigId]
	packs, err := folder.GetPacks(targetFolderManifest)
	if err != nil {
		return nil, err
	}

	packCount := len(packs)
	descriptors := make([]*EnumDescriptor, packCount)

	for packId := 0; packId < packCount; packId++ {
		descriptors[packId] = &EnumDescriptor{ID: uint32(packId)}

		packData := bytes.StringWrap(packs[packId].Data)
		if err := decodeEnum(packData, descriptors[packId]); err != nil {
			return nil, err
		}
	}

	return descriptors, nil
}

func decodeEnum(bs *bytes.String, descriptor *EnumDescriptor) error {
	itr := bs.Iterator()
	for itr.IsReadable() {
		id, err := itr.ReadByte()
		if err != nil {
			return err
		}

		if id == 0 {
			break
		}

		switch id {
		case 1:
			keyTypeValue, err := itr.ReadByte()
			if err != nil {
				return err
			}

			descriptor.KeyType = string(keyTypeValue)

		case 2:
			valueTypeValue, err := itr.ReadByte()
			if err != nil {
				return err
			}

			descriptor.KeyType = string(valueTypeValue)

		case 3:
			if descriptor.DefaultString, err = itr.ReadCString(); err != nil {
				return err
			}

		case 4:
			descriptor.DefaultInt, err = itr.ReadUInt32()
			if err != nil {
				return err
			}

		case 5:
			paramCount, err := itr.ReadUInt16()
			if err != nil {
				return err
			}

			descriptor.Parameters = make(map[int]interface{}, paramCount)

			for i := 0; i < int(paramCount); i++ {
				key, err := itr.ReadUInt32()
				if err != nil {
					return err
				}

				value, err := itr.ReadCString()
				if err != nil {
					return err
				}

				descriptor.Parameters[int(key)] = value
			}

		case 6:
			paramCount, err := itr.ReadUInt16()
			if err != nil {
				return err
			}

			descriptor.Parameters = make(map[int]interface{}, paramCount)

			for i := 0; i < int(paramCount); i++ {
				key, err := itr.ReadUInt32()
				if err != nil {
					return err
				}

				value, err := itr.ReadUInt32()
				if err != nil {
					return err
				}

				descriptor.Parameters[int(key)] = int(value)
			}

		default:
			return errors.New("could not find a case for id " + strconv.Itoa(int(id)) + "")
		}
	}

	return nil
}
