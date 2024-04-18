package chiapos

import (
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
	y := ComputeF1(15, 3, ys, 3*15)
	fmt.Printf("Y: %d\n", y)
	t.Fail()
}
