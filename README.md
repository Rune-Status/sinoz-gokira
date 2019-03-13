# GoKira

## Install

```
go get github.com/sinoz/gokira
```

## Supported Revisions

GoKira has been tested on the following revisions of OldSchool RuneScape:

- 177
- 178

## How To Use

Loading the cache is as easy as:

```
assetCache, err := gokira.LoadCache("cache/", 21)
```

If you are interested in the raw file data of the underlying file bundle, you can also do:

```
fileBundle, err := gokira.LoadFileBundle("cache/", 21)
if err != nil {
    log.Fatal(err)
}

assetCache, err := gokira.NewCache(fileBundle)
if err != nil {
    log.Fatal(err)
}
```

Developers who aren't very familiar with the cache might only want to use this library for streaming purposes in their server application. The game client expects to receive a collection of pages that together make up a categorized folder. To fetch such a folder (or the release manifest / update keys):

```
if archiveId == 255 && folderId == 255 {
    releaseManifest, err := assetCache.getReleaseManifest()
    if err != nil {
        return nil, err
    }
    
    // encodes the release manifest with all the versions and checksums
    // of each archive, into a buffer
    bufLength := len(releaseManifest.Checksums) * 8
    buf := buffer.NewHeapByteBuffer(bufLength)
    
    for i := 0; i < len(releaseManifest.Checksums); i++ {
    	buf.WriteInt32(int32(releaseManifest.Checksums[i]))
    	buf.WriteInt32(int32(releaseManifest.Versions[i]))
    }
    
    // reads the written contents into a byte array
    byteData := buf.ReadSlice(buf.ReadableBytes())
    
    // and return it
    return byteData, nil
} else {
    // getFolderPages() returns ([]byte, error)
    return assetCache.getFolderPages(archiveId, folderId)
}
```

To learn more on how to use this library for your OldSchool RuneScape application, check out the examples directory.

## Extras

#### HeapByteBuffer

GoKira also comes with its own buffer implementation called `HeapByteBuffer` which operates similarily to netty's `ByteBuf`. It grows exponentially when the buffer has reached its limit during a write operation. 

#### Supported Cryptographic/Compression Utilities

- XTEA (deciphering, enciphering)
- RSA (decrypting, encrypting)
- DJB2 (One-way hashing)
- GZIP (Decompression)
- BZIP2 (Decompression)

## FAQ

#### Why the name?

I was watching the anime Death Note and couldn't really think of anything else. I needed something so here we go. If you have any suggestions, feel free to make an issue or comment on an existing issue.

#### Does this also support applying cache modifications?

No. The focus of this library is to give developers a cache library to build (server) applications with. Additionally, the Go standard library currently, at the time of this writing, does not support Bzip2 encoding. Although Gzip can be used, perhaps one day.

#### Will this library also support RuneScape 3?

No.

## Giving Credits

- Sini for some of the namings
- Authors of OpenRS for illustrating the cache encoding format in their work
