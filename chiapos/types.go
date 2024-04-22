package chiapos

func NewTablesCache() *TablesCache {
	return &TablesCache{
		Buckets:     make([]Bucket, 0),
		RmapItem:    make([]RmapItem, 0),
		LeftTargets: calculate_left_targets(),
	}
}

type Bucket struct {
	BucketIndex   uint32
	StartPosition uint32
	Size          uint32
}

type Match struct {
	LeftPosition  uint32
	LeftY         uint32
	RightPosition uint32
}

type RmapItem struct {
	Count         uint32
	StartPosition uint32
}

type TablesCache struct {
	Buckets     []Bucket
	RmapItem    []RmapItem
	LeftTargets [][][]uint32
}

type BitSlice struct {
	HeadGap byte
	D       []byte
	TailGap byte
}
