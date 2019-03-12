package main

import (
	"github.com/sinoz/gokira"
	"log"
)

func main() {
	assetCache, err := gokira.LoadCache("cache/", 21)
	if err != nil {
		log.Fatal(err)
	}

	println(assetCache.ArchiveCount()) // 21
}
