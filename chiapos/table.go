package chiapos

import (
	"encoding/binary"

	"github.com/chocolatkey/chacha8"
)

const U8BITS = 8
const U32BITS = U8BITS * 4
const U64BITS = U8BITS * 8
const MAXU32 = 0xFFFFFFFF
const PARAM_EXT = 6
const PARAM_M = 1 << PARAM_EXT
const PARAM_B = 119
const PARAM_C = 127
const PARAM_BC = PARAM_B * PARAM_C

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

func calculate_left_targets() [][][]uint32 {
	res := make([][][]uint32, 2)
	for i := range res {
		parity := uint32(i)
		l1 := make([][]uint32, PARAM_BC)
		res[i] = l1
		var r uint32
		for r = 0; r < PARAM_BC; r++ {
			c := r / PARAM_C
			l2 := make([]uint32, PARAM_M)
			var m uint32
			for m = 0; m < PARAM_M; m++ {
				l2[m] = ((c+m)%PARAM_B)*PARAM_C + (((2*m+parity)*(2*m+parity) + r) % PARAM_C)
			}
			l1[r] = l2
		}

	}
	return res
}

func calculate_left_target_on_demand(parity, r, m uint) uint {
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
