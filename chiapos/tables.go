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
	k  byte
	TS [7]*table
}

func NewTableGeneric(k byte, seed []byte, cache *TablesCache) *TableGeneric {
	t1 := CreateTable1(k, seed)
	t2 := CreateTableN(k, 2, 1, t1, cache)
	t3 := CreateTableN(k, 3, 2, t2, cache)
	t4 := CreateTableN(k, 4, 3, t3, cache)
	t5 := CreateTableN(k, 5, 4, t4, cache)
	t6 := CreateTableN(k, 6, 5, t5, cache)
	t7 := CreateTableN(k, 7, 6, t6, cache)
	return &TableGeneric{
		k:  k,
		TS: [7]*table{t1, t2, t3, t4, t5, t6, t7},
	}
}

func (tg *TableGeneric) FindProof(chalenge []byte) (proof []byte, found bool) {
	t7 := tg.TS[6]
	firstKChallengeBits := binary.BigEndian.Uint32(chalenge[:4]) >> (U32BITS - int(tg.k))
	t7Len := len(t7.Items)
	// fmt.Printf("ys size: %d\nfirstKChallengeBits: %d\n", len(ys), firstKChallengeBits)
	pos := uint32(bsearch(t7Len, func(i int) int {
		y := t7.Items[i].y >> PARAM_EXT
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
	for ; pos < uint32(t7Len); pos++ {
		y := t7.Items[pos].y
		y = y >> PARAM_EXT
		if y > firstKChallengeBits {
			break
		}
		if y == firstKChallengeBits {
			// fmt.Printf("match pos: %d, y-: %d, y: %d\n", pos, y, ys[pos])
			choosedxs := make([]uint32, 0, 64)
			for _, pos6 := range tg.TS[6].Position(pos) {
				for _, pos5 := range tg.TS[5].Position(pos6) {
					for _, pos4 := range tg.TS[4].Position(pos5) {
						for _, pos3 := range tg.TS[3].Position(pos4) {
							for _, pos2 := range tg.TS[2].Position(pos3) {
								for _, pos1 := range tg.TS[1].Position(pos2) {
									choosedxs = append(choosedxs, tg.TS[0].Items[pos1].x)
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
