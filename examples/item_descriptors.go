package main

import (
	"errors"

	"github.com/sinoz/gokira"
	"github.com/sinoz/gokira/bytes"

	"log"
	"strconv"
)

const ItemConfigId = 10

type ItemDescriptor struct {
	ID                uint32    `yaml:"id"`
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
	assetCache, err := gokira.LoadCache("cache/", 21)
	if err != nil {
		log.Fatal(err)
	}

	descriptors, err := getItemDescriptors(assetCache)

	abyssalWhip := descriptors[4151] // 4151 = Abyssal Whip item
	println(abyssalWhip.Name)        // Abyssal whip
}

func getItemDescriptors(cache *gokira.Cache) ([]*ItemDescriptor, error) {
	archiveManifest, err := cache.GetArchiveManifest(2)
	if err != nil {
		return nil, err
	}

	folder, err := cache.GetUnencryptedFolder(2, ItemConfigId)
	if err != nil {
		return nil, err
	}

	targetFolderManifest := archiveManifest.FolderReferences[ItemConfigId]
	packs, err := folder.GetPacks(targetFolderManifest)
	if err != nil {
		return nil, err
	}

	packCount := len(packs)
	descriptors := make([]*ItemDescriptor, packCount)

	for id := 0; id < packCount; id++ {
		descriptors[id] = &ItemDescriptor{ID: uint32(id)}

		packData := bytes.StringWrap(packs[id].Data)
		if err := decodeItem(packData, descriptors[id]); err != nil {
			return nil, err
		}
	}

	return descriptors, nil
}

func decodeItem(bs *bytes.String, descriptor *ItemDescriptor) (err error) {
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
			if descriptor.InventoryModel, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 2:
			if descriptor.Name, err = itr.ReadCString(); err != nil {
				return err
			}

		case 4:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 5:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 6:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 7:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 8:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 11:
			descriptor.Stackable = true

		case 12:
			if descriptor.Cost, err = itr.ReadUInt32(); err != nil {
				return err
			}

		case 16:
			descriptor.Members = true

		case 23:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}
			if _, err = itr.ReadByte(); err != nil {
				return err
			}

		case 24:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 25:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

			if _, err = itr.ReadByte(); err != nil {
				return err
			}

		case 26:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 30, 31, 32, 33, 34:
			if descriptor.FloorOptions[id-30], err = itr.ReadCString(); err != nil {
				return err
			}

		case 35, 36, 37, 38, 39:
			if descriptor.BagOptions[id-35], err = itr.ReadCString(); err != nil {
				return err
			}

		case 40:
			count, _ := itr.ReadByte()
			for i := 0; i < int(count); i++ {
				if _, err = itr.ReadUInt16(); err != nil {
					return err
				}

				if _, err = itr.ReadUInt16(); err != nil {
					return err
				}
			}

		case 41:
			count, _ := itr.ReadByte()
			for i := 0; i < int(count); i++ {
				if _, err = itr.ReadUInt16(); err != nil {
					return err
				}

				if _, err = itr.ReadUInt16(); err != nil {
					return err
				}
			}

		case 42:
			if _, err = itr.ReadByte(); err != nil {
				return err
			}

		case 65:
			// STOCKMARKET

		case 78:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 79:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 90:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 91:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 92:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 93:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 95:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 97:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 98:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 100, 101, 102, 103, 104, 105, 106, 107, 108, 109:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 110:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 111:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 112:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 113:
			if _, err = itr.ReadByte(); err != nil {
				return err
			}

		case 114:
			if _, err = itr.ReadByte(); err != nil {
				return err
			}

		case 115:
			if _, err = itr.ReadByte(); err != nil {
				return err
			}

		case 139:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 140:
			if _, err := itr.ReadUInt16(); err != nil {
				return err
			}

		case 148:
			if descriptor.BankPlaceholderID, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 149:
			if _, err = itr.ReadUInt16(); err != nil {
				return err
			}

		case 249:
			count, err := itr.ReadByte()
			if err != nil {
				return err
			}

			for i := 0; i < int(count); i++ {
				flag, err := itr.ReadBool()
				if err != nil {
					return err
				}

				if _, err = itr.ReadUInt24(); err != nil {
					return err
				}

				if flag {
					if _, err = itr.ReadCString(); err != nil {
						return err
					}

				} else {
					if _, err = itr.ReadUInt32(); err != nil {
						return err
					}
				}
			}

		default:
			return errors.New("could not find a case for id " + strconv.Itoa(int(id)) + "")
		}
	}

	return nil
}
