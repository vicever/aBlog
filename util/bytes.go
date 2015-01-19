package util

import (
	"encoding/binary"
	"math"
)

func Int642Bytes(value int64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	binary.PutVarint(b, value)
	return b
}

func Bytes2Int64(b []byte) int64 {
	i, _ := binary.Varint(b)
	return i
}

func Float2Bytes(value float64) []byte {
	b := make([]byte, binary.MaxVarintLen64)
	i := math.Float64bits(value)
	binary.PutUvarint(b, i)
	return b
}

func Bytes2Float(b []byte) float64 {
	i, _ := binary.Uvarint(b)
	return math.Float64frombits(i)
}
