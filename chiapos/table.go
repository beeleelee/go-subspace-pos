package chiapos

import (
	"github.com/chocolatkey/chacha8"
)

const U8BITS = 8
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
