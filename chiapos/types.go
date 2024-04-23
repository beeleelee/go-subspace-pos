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

type table1 struct {
	k  byte
	ys []uint32
	xs []uint32
}

type t1Item struct {
	x uint32
	y uint32
}

type tablen struct {
	k         byte
	n         byte
	ys        []uint32
	positions [][]uint32
	metadata  [][]byte
}

type tnItem struct {
	y        uint32
	position []uint32
	metadata []byte
}
