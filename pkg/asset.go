package cache

// Asset is a simple asset file inside of a Pack.
type Asset struct {
	Data []byte
}

// newAsset constructs a new Asset file from the given data. May return an error.
func newAsset(data []byte) (*Asset, error) {
	return &Asset{Data: data}, nil
}
