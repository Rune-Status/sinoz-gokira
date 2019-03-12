package main

import (
	"errors"
	"github.com/sinoz/gokira/pkg"
	"github.com/sinoz/gokira/pkg/buffer"
	"log"
	"strconv"
)

const ItemConfigId = 10

type Descriptor struct {
	Id                uint32    `yaml:"id"`
	Name              string    `yaml:"name"`
	Examine           string    `yaml:"examine"`
	InventoryModel    uint16    `yaml:"inv_model"`
	Stackable         bool      `yaml:"can_stack"`
	Members           bool      `yaml:"members"`
	Cost              uint32    `yaml:"cost"`
	NotedID           int       `yaml:"noted"`
	BankPlaceholderID uint16    `yaml:"bank_placeholder_id"`
	FloorOptions      [5]string `yaml:"floor_opt"`
	BagOptions        [5]string `yaml:"bag_opt"`
}

func main() {
	fileBundle, err := cache.LoadFileBundle("cache/", 21)
	if err != nil {
		log.Fatal(err)
	}

	assetCache, err := cache.NewCache(fileBundle)
	if err != nil {
		log.Fatal(err)
	}

	descriptors, err := getDescriptors(assetCache)

	abyssalWhip := descriptors[4151] // 4151 = Abyssal Whip item
	println(abyssalWhip.Name)        // Abyssal whip
}

func getDescriptors(cache *cache.Cache) ([]*Descriptor, error) {
	archiveManifest, getManifestErr := cache.GetArchiveManifest(2)
	if getManifestErr != nil {
		return nil, getManifestErr
	}

	folder, folderPageErr := cache.GetUnencryptedFolder(2, ItemConfigId)
	if folderPageErr != nil {
		return nil, folderPageErr
	}

	targetFolderManifest := archiveManifest.FolderReferences[ItemConfigId]
	packs, getPacksErr := folder.GetPacks(targetFolderManifest)
	if getPacksErr != nil {
		return nil, getPacksErr
	}

	packCount := len(packs)
	descriptors := make([]*Descriptor, packCount)

	for id := 0; id < packCount; id++ {
		descriptors[id] = &Descriptor{Id: uint32(id)}

		packData := packs[id].Data
		packBuffer := buffer.HeapByteBufferWrap(packData)

		decodeError := decode(packBuffer, descriptors[id])
		if decodeError != nil {
			return nil, decodeError
		}
	}

	return descriptors, nil
}

func decode(buf *buffer.HeapByteBuffer, descriptor *Descriptor) error {
	for buf.IsReadable() {
		id, readErr := buf.ReadByte()
		if readErr != nil {
			return readErr
		}

		if id == 0 {
			break
		}

		switch id {
		case 1:
			descriptor.InventoryModel, _ = buf.ReadUInt16()

		case 2:
			descriptor.Name, _ = buf.ReadCString()

		case 4:
			buf.ReadUInt16()

		case 5:
			buf.ReadUInt16()

		case 6:
			buf.ReadUInt16()

		case 7:
			buf.ReadUInt16()

		case 8:
			buf.ReadUInt16()

		case 11:
			descriptor.Stackable = true

		case 12:
			descriptor.Cost, _ = buf.ReadUInt32()

		case 16:
			descriptor.Members = true

		case 23:
			buf.ReadUInt16()
			buf.ReadByte()

		case 24:
			buf.ReadUInt16()

		case 25:
			buf.ReadUInt16()
			buf.ReadByte()

		case 26:
			buf.ReadUInt16()

		case 30, 31, 32, 33, 34:
			descriptor.FloorOptions[id-30], _ = buf.ReadCString()

		case 35, 36, 37, 38, 39:
			descriptor.BagOptions[id-35], _ = buf.ReadCString()

		case 40:
			count, _ := buf.ReadByte()
			for i := 0; i < int(count); i++ {
				buf.ReadUInt16()
				buf.ReadUInt16()
			}

		case 41:
			count, _ := buf.ReadByte()
			for i := 0; i < int(count); i++ {
				buf.ReadUInt16()
				buf.ReadUInt16()
			}

		case 42:
			buf.ReadByte()

		case 65:
			// STOCKMARKET

		case 78:
			buf.ReadUInt16()

		case 79:
			buf.ReadUInt16()

		case 90:
			buf.ReadUInt16()

		case 91:
			buf.ReadUInt16()

		case 92:
			buf.ReadUInt16()

		case 93:
			buf.ReadUInt16()

		case 95:
			buf.ReadUInt16()

		case 97:
			buf.ReadUInt16()

		case 98:
			buf.ReadUInt16() // TODO noted template??

		case 100, 101, 102, 103, 104, 105, 106, 107, 108, 109:
			buf.ReadUInt16()
			buf.ReadUInt16()

		case 110:
			buf.ReadUInt16()

		case 111:
			buf.ReadUInt16()

		case 112:
			buf.ReadUInt16()

		case 113:
			buf.ReadByte()

		case 114:
			buf.ReadByte()

		case 115:
			buf.ReadByte()

		case 139:
			buf.ReadUInt16()

		case 140:
			buf.ReadUInt16()

		case 148:
			descriptor.BankPlaceholderID, _ = buf.ReadUInt16()

		case 149:
			buf.ReadUInt16()

		case 249:
			count, _ := buf.ReadByte()

			for i := 0; i < int(count); i++ {
				flag, _ := buf.ReadBool()
				buf.ReadUInt24()
				if flag {
					buf.ReadCString()
				} else {
					buf.ReadUInt32()
				}
			}

		default:
			return errors.New("could not find a case for id " + strconv.Itoa(int(id)) + "")
		}
	}

	return nil
}
