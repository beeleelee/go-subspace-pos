package chiapos

import (
	"encoding/binary"
)

func PickPosition(pos []uint32, last_5_challenge_bits, table_number byte) uint32 {
	if ((last_5_challenge_bits >> (table_number - 2)) & 1) == 0 {
		return pos[0]
	}
	return pos[1]
}

type TableGeneric struct {
	k      byte
	Table1 Table
	Table2 Table
	Table3 Table
	Table4 Table
	Table5 Table
	Table6 Table
	Table7 Table
}

func NewTableGeneric(k byte, seed []byte, cache *TablesCache) *TableGeneric {
	t1 := CreateTable1(k, seed)
	t2 := CreateTablen(k, 2, 1, t1, cache)
	t3 := CreateTablen(k, 3, 2, t2, cache)
	t4 := CreateTablen(k, 4, 3, t3, cache)
	t5 := CreateTablen(k, 5, 4, t4, cache)
	t6 := CreateTablen(k, 6, 5, t5, cache)
	t7 := CreateTablen(k, 7, 6, t6, cache)
	return &TableGeneric{
		k:      k,
		Table1: t1,
		Table2: t2,
		Table3: t3,
		Table4: t4,
		Table5: t5,
		Table6: t6,
		Table7: t7,
	}
}

func (tg *TableGeneric) FindProof(chalenge []byte) (proof []byte, found bool) {
	ys := tg.Table7.YS()
	firstKChallengeBits := binary.BigEndian.Uint32(chalenge[:4]) >> (U32BITS - int(tg.k))

	// fmt.Printf("ys size: %d\nfirstKChallengeBits: %d\n", len(ys), firstKChallengeBits)
	pos := uint32(bsearch(len(ys), func(i int) int {
		y := ys[i] >> PARAM_EXT
		if y == firstKChallengeBits {
			return 0
		} else if y < firstKChallengeBits {
			return -1
		} else {
			return 1
		}
	}))
	// fmt.Printf("%v\n", ys[pos-3:pos+3])
	// fmt.Printf("search pos: %d\n", pos)
	for ; pos < uint32(len(ys)); pos++ {
		y := ys[pos]
		y = y >> PARAM_EXT
		if y > firstKChallengeBits {
			break
		}
		if y == firstKChallengeBits {
			// fmt.Printf("match pos: %d, y-: %d, y: %d\n", pos, y, ys[pos])
			xs := tg.Table1.XS()
			choosedxs := make([]uint32, 0, 64)
			for _, pos6 := range tg.Table7.Position(pos) {
				for _, pos5 := range tg.Table6.Position(pos6) {
					for _, pos4 := range tg.Table5.Position(pos5) {
						for _, pos3 := range tg.Table4.Position(pos4) {
							for _, pos2 := range tg.Table3.Position(pos3) {
								for _, pos1 := range tg.Table2.Position(pos2) {
									choosedxs = append(choosedxs, xs[pos1])
								}
							}
						}
					}
				}
			}
			bslices := make([]BitSlice, len(choosedxs))
			keep := divCeil(uint(tg.k), U8BITS)
			headgap := tg.k % U8BITS
			if headgap > 0 {
				headgap = U8BITS - headgap
			}
			for i, x := range choosedxs {
				// fmt.Printf("offset: %d, x: %d\n", i, x)
				bs := make([]byte, 4)
				binary.BigEndian.PutUint32(bs, x)
				bslices[i] = BitSlice{
					HeadGap: headgap,
					D:       bs[len(bs)-int(keep):],
				}
			}
			proof = bitsConcatLeft(bslices).D
			found = true
			break
		}
	}
	return
}
