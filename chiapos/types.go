package chiapos

func NewTablesCache() *TablesCache {
	return &TablesCache{
		Buckets:     make([]Bucket, 0, MAX_BUCKET_SIZE),
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
	LeftTargets []uint32
}

type BitSlice struct {
	HeadGap byte
	D       []byte
	TailGap byte
}

type table struct {
	n     byte
	Items []tItem
}

type tItem struct {
	x        uint32
	y        uint32
	position []uint32
	metadata []byte
}

func (t *table) YS() []uint32 {
	r := make([]uint32, len(t.Items))
	for i, item := range t.Items {
		r[i] = item.y
	}
	return r
}
