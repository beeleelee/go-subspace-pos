package chiapos

import (
	"lukechampine.com/blake3"
)

func divCeil(a, b uint) uint {
	res := a / b
	if a%b > 0 {
		res++
	}
	return res
}

func blake3Hash(data []byte) []byte {
	h := blake3.Sum256(data)
	return h[:]
}

func bitsConcatLeft(bs []BitSlice) (ot BitSlice) {
	if len(bs) == 0 {
		return
	}
	ot.D = make([]byte, 0)
	var lastByteTailGap byte
	otDLastIdx := -1
	for _, bsItem := range bs {
		if len(bsItem.D) == 0 {
			continue
		}
		firstByte := bsItem.D[0]
		innerBytes := bsItem.D[1:]
		if otDLastIdx == -1 {
			lastByteTailGap = bsItem.HeadGap
			// fmt.Printf("firstByte %08b\n", firstByte)
			ot.D = append(ot.D, firstByte<<lastByteTailGap)
			otDLastIdx++
		} else {
			prevbyte := ot.D[otDLastIdx]
			// fmt.Printf("prevByte %08b\nlastByteTailGap %d\n cur %08b\n", prevbyte, lastByteTailGap, firstByte)

			ot.D[otDLastIdx] = prevbyte | ((firstByte << bsItem.HeadGap) >> (8 - lastByteTailGap))
			if 8-bsItem.HeadGap > lastByteTailGap {
				// fmt.Printf("8-bsItem.HeadGap > lastByteTailGap  \n")
				ot.D = append(ot.D, firstByte<<(bsItem.HeadGap+lastByteTailGap))
				otDLastIdx++
				lastByteTailGap += bsItem.HeadGap
			} else {
				// fmt.Println("8-bsItem.HeadGap <= lastByteTailGap")
				lastByteTailGap = lastByteTailGap + bsItem.HeadGap - 8
			}
		}
		endIdx := len(innerBytes) - 1
		if endIdx < 0 {
			ot.D[otDLastIdx] = (ot.D[otDLastIdx] >> bsItem.TailGap) << bsItem.TailGap
			lastByteTailGap += bsItem.TailGap
		}
		for j, d := range innerBytes {
			tailGap := byte(0)
			if j == endIdx {
				d = (d >> bsItem.TailGap) << bsItem.TailGap
				tailGap = bsItem.TailGap
			}
			prevbyte := ot.D[otDLastIdx]
			// fmt.Printf("prevByte %08b\nlastByteTailGap %d\ncur %08b\n", prevbyte, lastByteTailGap, d)

			ot.D[otDLastIdx] = prevbyte | (d >> (8 - lastByteTailGap))
			if 8-tailGap > lastByteTailGap {
				ot.D = append(ot.D, d<<lastByteTailGap)
				otDLastIdx++
				lastByteTailGap += tailGap
			} else {
				lastByteTailGap = lastByteTailGap + tailGap - 8
			}
		}
	}
	ot.TailGap = lastByteTailGap
	return
}

func bitsConcatRight(bs []BitSlice) (ot BitSlice) {
	if len(bs) == 0 {
		return
	}
	ot.D = make([]byte, 0)
	var prevByteHeadGap byte
	otDLastIdx := -1
	for ri := len(bs) - 1; ri > -1; ri-- {
		bsItem := bs[ri]
		if len(bsItem.D) == 0 {
			continue
		}
		endIdx := len(bsItem.D) - 1
		lastByte := bsItem.D[endIdx]
		innerBytes := bsItem.D[0:endIdx]
		if otDLastIdx == -1 {
			prevByteHeadGap = bsItem.TailGap
			// fmt.Printf("lastByte %08b\n", lastByte)
			ot.D = append(ot.D, lastByte>>prevByteHeadGap)
			otDLastIdx++
		} else {
			prevbyte := ot.D[otDLastIdx]
			// fmt.Printf("prevByte %08b\nprevByteHeadGap %d\n cur %08b\n", prevbyte, prevByteHeadGap, lastByte)

			ot.D[otDLastIdx] = prevbyte | ((lastByte >> bsItem.TailGap) << (8 - prevByteHeadGap))
			if 8-bsItem.TailGap > prevByteHeadGap {
				// fmt.Printf("8-bsItem.TailGap > prevByteHeadGap  \n")
				ot.D = append(ot.D, lastByte>>(bsItem.TailGap+prevByteHeadGap))
				otDLastIdx++
				prevByteHeadGap += bsItem.TailGap
			} else {
				// fmt.Println("8-bsItem.TailGap <= prevByteHeadGap")
				prevByteHeadGap = prevByteHeadGap + bsItem.TailGap - 8
			}
		}
		endIdx--
		if endIdx < 0 {
			ot.D[otDLastIdx] = (ot.D[otDLastIdx] << bsItem.HeadGap) >> bsItem.HeadGap
			prevByteHeadGap += bsItem.HeadGap
		}
		for rj := endIdx; rj > -1; rj-- {
			d := innerBytes[rj]
			headGap := byte(0)
			if rj == 0 {
				d = (d << bsItem.HeadGap) >> bsItem.HeadGap
				headGap = bsItem.HeadGap
			}
			prevbyte := ot.D[otDLastIdx]
			// fmt.Printf("prevByte %08b\nprevByteHeadGap %d\ncur %08b\n", prevbyte, prevByteHeadGap, d)

			ot.D[otDLastIdx] = prevbyte | (d << (8 - prevByteHeadGap))
			if 8-headGap > prevByteHeadGap {
				ot.D = append(ot.D, d>>prevByteHeadGap)
				otDLastIdx++
				prevByteHeadGap += headGap
			} else {
				prevByteHeadGap += headGap - 8
			}
		}

	}
	ot.HeadGap = prevByteHeadGap
	//reverse ot.D
	for i, j := 0, len(ot.D)-1; i < j; i, j = i+1, j-1 {
		ot.D[i], ot.D[j] = ot.D[j], ot.D[i]
	}
	return
}

func bsearch(size int, f func(int) int) int {
	left := 0
	right := size
	for left < right {
		mid := left + size/2
		cmp := f(mid)
		if cmp == 0 {
			return mid
		} else if cmp < 0 {
			left = mid + 1
		} else {
			right = mid
		}
		size = right - left
	}
	return left
}

func bsu32(list []uint32, x uint32) int {
	return bsearch(len(list), func(i int) int {
		y := list[i]
		if y == x {
			return 0
		} else if y < x {
			return -1
		} else {
			return 1
		}
	})
}
