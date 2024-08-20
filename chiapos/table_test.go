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
	ys := partial_ys(15, seed)
	y := ComputeF1(15, 5, ys, 5*15)
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
		21, 185, 27, 119, 214, 189, 172, 168,
		255, 193, 47, 112, 202, 51, 192, 31,
		33, 167, 102, 81, 207, 18, 55, 9,
		77, 234, 158, 72, 121, 171, 137, 229,
	}
	t1 := CreateTable1(k, seed)
	t1ys := t1.YS()
	fmt.Printf("t1ys %d %v\n", len(t1ys), t1ys[:8])
	t1xs := t1.XS()
	fmt.Printf("t1xs %d %v\n", len(t1xs), t1xs[:8])
	cache := &TablesCache{
		LeftTargets: calculate_left_targets(),
	}
	t2 := CreateTableN(k, 2, 1, t1, cache)
	t2ys := t2.YS()
	fmt.Printf("t2 %d %v\n", len(t2ys), t2ys[:8])
	t3 := CreateTableN(k, 3, 2, t2, cache)
	t3ys := t3.YS()
	fmt.Printf("t3 %d %v\n", len(t3ys), t3ys[:4])
	t4 := CreateTableN(k, 4, 3, t3, cache)
	t4ys := t4.YS()
	fmt.Printf("t4 %d %v\n", len(t4ys), t4ys[:4])
	t5 := CreateTableN(k, 5, 4, t4, cache)
	t5ys := t5.YS()
	fmt.Printf("t5 %d %v\n", len(t5ys), t5ys[:4])
	t6 := CreateTableN(k, 6, 5, t5, cache)
	t6ys := t6.YS()
	fmt.Printf("t6 %d %v\n", len(t6ys), t6ys[:4])
	t7 := CreateTableN(k, 7, 6, t6, cache)
	t7ys := t7.YS()
	fmt.Printf("t7 %d %v\n", len(t7ys), t7ys[:4])
	t.Fail()
}

func TestSortSearch(t *testing.T) {
	list := []uint32{0, 1, 1, 2, 3, 5, 5, 5, 6, 7, 8, 8, 8, 8, 9}
	idx := bsu32(list, 5)
	fmt.Println(idx)
	t.Fail()
}

func TestB4(t *testing.T) {
	k := byte(16)
	seed := []byte{
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1,
	}
	// ys := partial_ys(k, seed)
	// fmt.Println(ys[:], len(ys))
	// var x uint32
	// for ; x < 32; x++ {
	// 	y := ComputeF1(k, x, ys, uint(k)*uint(x))
	// 	fmt.Printf("x: %d, y: %d\n", x, y)
	// }
	cache := NewTablesCache()

	table := NewTableGeneric(k, seed, cache)
	// t1ys := table.TS[0].YS()
	// fmt.Println(t1ys[:4], len(t1ys))
	// t2ys := table.TS[1].YS()
	// fmt.Println(t2ys, len(t2ys))
	challengeIndex := uint32(800)
	// fmt.Printf("seed: %v\nchallengeIndex: %d\n", seed, challengeIndex)
	challenge := make([]byte, 32)
	binary.LittleEndian.PutUint32(challenge[:4], challengeIndex)
	// fmt.Printf("challenge: %v\n", challenge)

	proof, _ := table.FindProof(challenge)
	fmt.Printf("proof %v\n", proof)

	t.Fail()
}

// k16 404  3 proof
// k16 521 3 proof
// k16 568 6 proof
// k16 700 3 proof
// k16 800 3 proof
