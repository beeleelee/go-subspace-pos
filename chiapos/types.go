package chiapos

func NewTablesCache() *TablesCache {
	return &TablesCache{
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
	LeftTargets []uint32
}

type BitSlice struct {
	HeadGap byte
	D       []byte
	TailGap byte
}

type table struct {
	k     byte
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

func (t *table) XS() []uint32 {
	if t.n != 1 {
		return nil
	}
	r := make([]uint32, len(t.Items))
	for i, item := range t.Items {
		r[i] = item.x
	}
	return r
}
