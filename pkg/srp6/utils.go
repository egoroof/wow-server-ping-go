package srp6

import (
	"math/big"
)

// returns LE big integer from BE byte array
func bytesToInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(reverse(data))
}

// returns BE byte array from LE big integer
func intToBytes(padding int, bi *big.Int) []byte {
	return reverse(pad(padding, bi.Bytes()))
}

// returns a copy of data appended with zeroed bytes. Padding is added until len(data) == length.
// Initially, if len(data) >= length, no padding is added and the original data is returned.
func pad(resLength int, data []byte) []byte {
	dataLen := len(data)
	if dataLen >= resLength {
		return data
	}
	res := make([]byte, resLength)
	copy(res[resLength-dataLen:], data)
	return res
}

// returns a copy of data in reverse order
func reverse(data []byte) []byte {
	n := len(data)
	newData := make([]byte, n)
	for i := range n {
		newData[i] = data[n-i-1]
	}
	return newData
}
