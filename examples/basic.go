package main

import (
	"github.com/sinoz/gokira/pkg"
	"log"
)

func main() {
	assetCache, err := cache.LoadCache("cache/", 21)
	if err != nil {
		log.Fatal(err)
	}

	println(assetCache.ArchiveCount()) // 21
}
