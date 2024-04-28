package chiapos

func PickPosition(pos []uint32, last_5_challenge_bits, table_number byte) uint32 {
	if ((last_5_challenge_bits >> (table_number - 2)) & 1) == 0 {
		return pos[0]
	}
	return pos[1]
}

type TableGeneric struct {
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
		Table1: t1,
		Table2: t2,
		Table3: t3,
		Table4: t4,
		Table5: t5,
		Table6: t6,
		Table7: t7,
	}
}

func (tg *TableGeneric) FindProof(chalenge []byte) (proof []byte) {
	return
}
