package chiapos

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestPY(t *testing.T) {
	seed := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	ys := partial_ys(15, seed)
	fmt.Println(ys[len(ys)-32:])
	t.Fail()
}

func TestF1(t *testing.T) {
	seed := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	ys := partial_ys(20, seed)
	y := ComputeF1(20, 3, ys, 3*20)
	fmt.Printf("Y: %d\n", y)
	t.Fail()
}

func TestB3(t *testing.T) {
	data := []byte{0, 1, 2, 3, 4, 5}
	h := blake3Hash(data)
	fmt.Println(h)
	t.Fail()
}

func TestLC(t *testing.T) {
	x := uint32(781432176)
	y := uint32(234356742)
	bytesx := make([]byte, 4)
	bytesy := make([]byte, 4)
	binary.BigEndian.PutUint32(bytesx, x)
	binary.BigEndian.PutUint32(bytesy, y)
	list := make([]BitSlice, 2)
	fmt.Println(bytesx)
	fmt.Println(bytesy)
	list[0] = BitSlice{
		HeadGap: 2,
		D:       bytesx[:],
	}
	list[1] = BitSlice{
		HeadGap: 4,
		D:       bytesy[:],
	}
	bs := bitsConcatLeft(list)

	fmt.Printf("headGap %d\n", bs.HeadGap)
	fmt.Printf("tailGap %d\n", bs.TailGap)
	for _, b := range bs.D {
		fmt.Printf("%08b\n", b)
	}
	fmt.Println(bs.D)
	x1 := uint64(x)
	y1 := uint64(y)
	z1 := (x1 << 34) | (y1 << 6)
	bytesz := make([]byte, 8)
	binary.BigEndian.PutUint64(bytesz, z1)
	fmt.Printf("%032b\n%032b\n%064b\n", x1, y1, z1)
	fmt.Println(bytesz)
	t.Fail()
}

func TestRC(t *testing.T) {
	x := uint32(781432176)
	y := uint32(234356742)
	bytesx := make([]byte, 4)
	bytesy := make([]byte, 4)
	binary.BigEndian.PutUint32(bytesx, x)
	binary.BigEndian.PutUint32(bytesy, y)
	list := make([]BitSlice, 2)
	fmt.Println(bytesx)
	fmt.Println(bytesy)
	list[0] = BitSlice{
		HeadGap: 4,
		D:       bytesx[:],
		TailGap: 0,
	}
	list[1] = BitSlice{
		HeadGap: 0,
		D:       bytesy[:],
		TailGap: 1,
	}
	bs := bitsConcatRight(list)

	fmt.Printf("headGap %d\n", bs.HeadGap)
	fmt.Printf("tailGap %d\n", bs.TailGap)
	for _, b := range bs.D {
		fmt.Printf("%08b\n", b)
	}
	fmt.Println(bs.D)
	x1 := uint64(x)
	y1 := uint64(y)
	z1 := ((x1 << 36) >> 5) | (y1 >> 1)
	bytesz := make([]byte, 8)
	binary.BigEndian.PutUint64(bytesz, z1)
	fmt.Printf("%032b\n%032b\n%064b\n", x1, y1, z1)
	fmt.Println(bytesz)
	t.Fail()
}

func TestTables(t *testing.T) {
	k := byte(20)
	seed := []byte{
		21, 185, 23, 119, 214, 189, 172, 168,
		255, 193, 47, 112, 202, 51, 192, 31,
		33, 167, 102, 81, 207, 18, 55, 9,
		77, 234, 158, 72, 106, 171, 137, 229,
	}
	t1 := CreateTable1(k, seed)
	t1ys := t1.YS()
	fmt.Printf("t1 %d %v\n", len(t1ys), t1ys[:4])
	cache := &TablesCache{
		LeftTargets: calculate_left_targets(),
	}
	t2 := CreateTablen(k, 2, 1, t1, cache)
	t2ys := t2.YS()
	fmt.Printf("t2 %d %v\n", len(t2ys), t2ys[:4])
	t3 := CreateTablen(k, 3, 2, t2, cache)
	t3ys := t3.YS()
	fmt.Printf("t3 %d %v\n", len(t3ys), t3ys[:4])
	t4 := CreateTablen(k, 4, 3, t3, cache)
	t4ys := t4.YS()
	fmt.Printf("t4 %d %v\n", len(t4ys), t4ys[:4])
	t5 := CreateTablen(k, 5, 4, t4, cache)
	t5ys := t5.YS()
	fmt.Printf("t5 %d %v\n", len(t5ys), t5ys[:4])
	t6 := CreateTablen(k, 6, 5, t5, cache)
	t6ys := t6.YS()
	fmt.Printf("t6 %d %v\n", len(t6ys), t6ys[:4])
	t7 := CreateTablen(k, 7, 6, t6, cache)
	t7ys := t7.YS()
	fmt.Printf("t7 %d %v\n", len(t7ys), t7ys[:4])
	t.Fail()
}

func TestSortSearch(t *testing.T) {
	list := []uint32{0, 1, 1, 2, 3, 5, 5, 5, 6, 7, 8, 8, 8, 8, 9}
	idx := bsu32(list, 8)
	fmt.Println(idx)
	t.Fail()
}

func TestB4(t *testing.T) {
	seed := []byte{
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
	}

	cache := NewTablesCache()
	table := NewTableGeneric(17, seed, cache)
	challengeIndex := uint32(212783992)
	// fmt.Printf("seed: %v\nchallengeIndex: %d\n", seed, challengeIndex)
	challenge := make([]byte, 32)
	binary.LittleEndian.PutUint32(challenge[:4], challengeIndex)
	// fmt.Printf("challenge: %v\n", challenge)
	proof, _ := table.FindProof(challenge)
	fmt.Printf("proof %v\n", proof)

	t.Fail()
}

// 955453486  1 proof
// 1378873433 3 proof
// 3501970825
// 3259531167
