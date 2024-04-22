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
