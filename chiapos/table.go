package chiapos

import (
	"encoding/binary"
	"sort"

	"github.com/chocolatkey/chacha8"
)

const k = 20
const U8BITS = 8
const U32BITS = U8BITS * 4
const U64BITS = U8BITS * 8
const MAXU32 = 0xFFFFFFFF

const PARAM_EXT = 6
const PARAM_M = 1 << PARAM_EXT
const PARAM_B = 119
const PARAM_C = 127
const PARAM_BC = PARAM_B * PARAM_C
const Y_SIZE_BITS = k + PARAM_EXT
const METADATA_SIZE_BITS_1 = 1 * k
const METADATA_SIZE_BITS_2 = 2 * k
const METADATA_SIZE_BITS_3 = 4 * k
const METADATA_SIZE_BITS_4 = 4 * k
const METADATA_SIZE_BITS_5 = 3 * k
const METADATA_SIZE_BITS_6 = 2 * k
const METADATA_SIZE_BITS_7 = 0 * k
const T1_ITEM_SIZE = 1 << k
const MAX_BUCKET_SIZE = (T1_ITEM_SIZE/PARAM_BC + 1000)

func YSizeBits(k byte) uint {
	return uint(k + PARAM_EXT)
}

func MetadataSizeBytes(k, tableNubmer byte) uint {
	return divCeil(MetadataSizeBits(k, tableNubmer), U8BITS)
}

func MetadataSizeBits(k, tableNumber byte) (res uint) {
	switch tableNumber {
	case 1:
		res = 1
	case 2:
		res = 2
	case 3:
		res = 4
	case 4:
		res = 4
	case 5:
		res = 3
	case 6:
		res = 2
	default:
		res = 0
	}
	res = res * uint(k)
	return
}

func partial_ys(k byte, seed []byte) []byte {
	outputLenBits := uint(k) * (uint(1) << k)
	outputLen := divCeil(outputLenBits, U8BITS)
	output := make([]byte, outputLen)
	nonce := make([]byte, 12)
	cipher, err := chacha8.New(seed, nonce)
	if err != nil {
		panic(err)
	}
	cipher.KeyStream(output)
	return output
}

func calculate_left_targets() []uint32 {
	res := make([]uint32, 2*PARAM_BC*PARAM_M)
	var parity uint32
	for parity = 0; parity <= 1; parity++ {
		var r uint32
		for r = 0; r < PARAM_BC; r++ {
			c := r / PARAM_C
			var m uint32
			for m = 0; m < PARAM_M; m++ {
				cur := parity*PARAM_BC + r*PARAM_M + m
				res[cur] = ((c+m)%PARAM_B)*PARAM_C + (((2*m+parity)*(2*m+parity) + r) % PARAM_C)
			}
		}
	}
	return res
}

func calculate_left_target_on_demand(parity, r, m uint32) uint32 {
	c := r / PARAM_C
	return ((c+m)%PARAM_B)*PARAM_C + (((2*m+parity)*(2*m+parity) + r) % PARAM_C)
}

func ComputeF1(k byte, x uint32, partialY []byte, partialYOffset uint) (y uint32) {
	var preYMask, preExtMask uint32
	partialYLength := divCeil(partialYOffset%U8BITS+uint(k), U8BITS)
	preYBytes := make([]byte, 8)
	copy(preYBytes[:partialYLength], partialY[partialYOffset/U8BITS:partialYOffset/U8BITS+partialYLength])
	preY := uint32(binary.BigEndian.Uint64(preYBytes) >> (U64BITS - uint(k) - PARAM_EXT - partialYOffset%U8BITS))
	mx32 := uint64(MAXU32)
	preYMask = uint32(mx32<<PARAM_EXT) & (MAXU32 >> (U32BITS - k - PARAM_EXT))
	preExt := x >> (k - PARAM_EXT)
	preExtMask = MAXU32 >> (U32BITS - PARAM_EXT)
	return (preY & preYMask) | (preExt & preExtMask)
}

func findMatches(last_table *table, left_bucket, right_bucket Bucket, left_targets []uint32) (ms []Match) {
	if left_bucket.Size == 0 || right_bucket.Size == 0 {
		return
	}
	rmap := make([]RmapItem, PARAM_BC)
	base := (last_table.Items[right_bucket.StartPosition].y / PARAM_BC) * PARAM_BC
	for i := uint32(0); i < right_bucket.Size; i++ {
		rightPosition := i + right_bucket.StartPosition
		r := last_table.Items[rightPosition].y - base
		if rmap[r].Count == 0 {
			rmap[r].StartPosition = rightPosition
		}
		rmap[r].Count += 1
	}
	base = base - PARAM_BC
	parity := (last_table.Items[left_bucket.StartPosition].y / PARAM_BC) % 2
	ms = make([]Match, 0)

	for i := uint32(0); i < left_bucket.Size; i++ {
		leftPosition := i + left_bucket.StartPosition
		y := last_table.Items[leftPosition].y
		r := y - base
		for m := 0; m < PARAM_M; m++ {
			rt := left_targets[parity*PARAM_BC+r*PARAM_M+uint32(m)]
			ritem := rmap[rt]
			if ritem.Count > 0 {
				for rightPosition := ritem.StartPosition; rightPosition < ritem.StartPosition+ritem.Count; rightPosition++ {
					ms = append(ms, Match{
						LeftPosition:  leftPosition,
						LeftY:         y,
						RightPosition: rightPosition,
					})
				}
			}
		}
	}
	return
}

func NumMatches(left_y, right_y uint32) (matches uint) {
	right_r := right_y % PARAM_BC
	parity := (left_y / PARAM_BC) % 2
	left_r := left_y % PARAM_BC

	for m := 0; m < PARAM_M; m++ {
		r_target := calculate_left_target_on_demand(parity, left_r, uint32(m))
		if r_target == right_r {
			matches++
		}
	}
	return
}

func ComputeFN(k, tn, pvtn byte, y uint32, left_metadata, right_metadata []byte) (cy uint32, cm []byte) {
	//parentMetadataBits := MetadataSizeBits(k, pvtn)
	// fmt.Printf("y: %d, left_m: %v, right_m: %v\n", y, left_metadata, right_metadata)
	ySizeBits := YSizeBits(k)
	mdSizeBits := MetadataSizeBits(k, pvtn)
	yBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(yBytes, y)
	ySizeBytes := divCeil(ySizeBits, U8BITS)
	headGap := ySizeBits % U8BITS
	if headGap > 0 {
		headGap = U8BITS - headGap
	}
	ySlice := BitSlice{
		HeadGap: byte(headGap),
		D:       yBytes[len(yBytes)-int(ySizeBytes):],
	}
	// fmt.Println("^^^^^^^^^")
	// fmt.Println(ySlice)
	headGap = mdSizeBits % U8BITS
	if headGap > 0 {
		headGap = U8BITS - headGap
	}
	lmSlice := BitSlice{
		HeadGap: byte(headGap),
		D:       left_metadata,
	}
	// fmt.Println(lmSlice)
	rmSlice := BitSlice{
		HeadGap: byte(headGap),
		D:       right_metadata,
	}
	// fmt.Println(rmSlice)
	// fmt.Println("--------")
	elements := []BitSlice{
		ySlice,
		lmSlice,
		rmSlice,
	}
	ot := bitsConcatLeft(elements)
	hash := blake3Hash(ot.D)
	// if tn == 7 {
	// 	fmt.Printf("y: %d, lm: %v, rm: %v, cat: %v\nhash: %v\n", y, left_metadata, right_metadata, ot.D, hash)
	// }

	cy = binary.BigEndian.Uint32(hash[:4])
	cy = cy >> (U32BITS - ySizeBits)
	// if tn == 7 {
	// 	fmt.Printf("%d from hash %v\n", cy, hash[:4])
	// }
	mdsb := MetadataSizeBits(k, tn)
	if tn < 4 {
		cm = bitsConcatRight(elements[1:]).D
	} else if mdsb > 0 {
		start := ySizeBits / U8BITS
		headGap := ySizeBits % U8BITS
		h := hash[start : start+divCeil(mdsb+headGap, U8BITS)]
		s := []BitSlice{
			{
				HeadGap: byte(headGap),
				D:       h,
				TailGap: byte(uint(len(h)*U8BITS) - mdsb - headGap),
			},
		}
		// if tn == 7 {
		// 	fmt.Printf("%v\n", s)
		// }
		cm = bitsConcatRight(s).D
		// if tn == 7 {
		// 	fmt.Printf("ySizeBits: %d, mdsb: %d,  h: %v\ncm: %v\n", ySizeBits, mdsb, h, cm)
		// }
	}

	return
}

func matchToResult(k, tn, pvtn byte, last_table *table, m Match) tItem {
	left_metadata := last_table.Metadata(m.LeftPosition)
	right_metadata := last_table.Metadata(m.RightPosition)
	y, metadata := ComputeFN(k, tn, pvtn, m.LeftY, left_metadata, right_metadata)
	return tItem{
		y:        y,
		position: []uint32{m.LeftPosition, m.RightPosition},
		metadata: metadata,
	}
}

func matchAndComputFn(k, tn, pvtn byte, last_table *table, left_bucket, right_bucket Bucket, left_targets []uint32) (results_table []tItem) {
	if left_bucket.Size == 0 || right_bucket.Size == 0 {
		return
	}
	if left_bucket.BucketIndex+1 != right_bucket.BucketIndex {
		return
	}
	results_table = make([]tItem, 0)

	// matches := findMatches(ys[left_bucket.StartPosition:left_bucket.StartPosition+left_bucket.Size], left_bucket.StartPosition, ys[right_bucket.StartPosition:right_bucket.StartPosition+right_bucket.Size], right_bucket.StartPosition, left_targets)
	matches := findMatches(last_table, left_bucket, right_bucket, left_targets)
	for _, m := range matches {
		// if tn == 7 {
		// 	fmt.Println(m)
		// }
		results_table = append(results_table, matchToResult(k, tn, pvtn, last_table, m))
	}
	return
}

func CreateTable1(k byte, seed []byte) *table {
	len := uint32(1 << k)
	partialYS := partial_ys(k, seed)
	t1 := &table{
		n:     1,
		Items: make([]tItem, len),
	}
	var x uint32
	for ; x < len; x++ {
		y := ComputeF1(k, x, partialYS, uint(k)*uint(x))
		t1.Items[x] = tItem{
			x: x,
			y: y,
		}
	}
	sort.Slice(t1.Items, func(i, j int) bool {
		return t1.Items[i].y < t1.Items[j].y
	})

	return t1
}

func (t *table) Position(pos uint32) []uint32 {
	if t.n == 1 {
		return nil
	}
	return t.Items[pos].position
}

func (t *table) Metadata(pos uint32) []byte {
	if t.n == 1 {
		x := t.Items[pos].x
		n := int(MetadataSizeBytes(k, 1))
		bs := make([]byte, 4)
		binary.BigEndian.PutUint32(bs, x)
		return bs[len(bs)-n:]
	}
	return t.Items[pos].metadata
}

func CreateTableN(k, tn, pvtn byte, last_table *table, cache *TablesCache) *table {
	buckets := cache.Buckets[:0]
	bucket := &Bucket{}

	for i, item := range last_table.Items {
		buketIndex := item.y / PARAM_BC
		if buketIndex == bucket.BucketIndex {
			bucket.Size += 1
			continue
		}
		buckets = append(buckets, *bucket)
		bucket = &Bucket{
			BucketIndex:   buketIndex,
			StartPosition: uint32(i),
			Size:          1,
		}
	}
	buckets = append(buckets, *bucket)
	num_values := 1 << k
	if len(buckets) < 2 {
		return &table{}
	}
	table := &table{
		n:     tn,
		Items: make([]tItem, 0, num_values),
	}
	for i := 0; i < len(buckets)-1; i++ {
		left_bucket := buckets[i]
		right_bucket := buckets[i+1]
		// fmt.Println(left_bucket)
		// fmt.Println(right_bucket)
		table.Items = append(table.Items, matchAndComputFn(k, tn, pvtn, last_table, left_bucket, right_bucket, cache.LeftTargets)...)
	}
	sort.Slice(table.Items, func(i, j int) bool {
		return table.Items[i].y < table.Items[j].y
	})
	return table
}
