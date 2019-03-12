package main

import (
	"github.com/sinoz/gokira/pkg"
	"log"
)

func main() {
	fileBundle, err := cache.LoadFileBundle("cache/", 21)
	if err != nil {
		log.Fatal(err)
	}

	assetCache, err := cache.NewCache(fileBundle)
	if err != nil {
		log.Fatal(err)
	}

	println(assetCache.ArchiveCount()) // 21
}
