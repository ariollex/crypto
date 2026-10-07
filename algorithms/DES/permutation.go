package DES

import (
	"errors"
	"fmt"
)

type BitOrder uint8

const (
	LeastSignificantFirst BitOrder = iota
	MostSignificantFirst
)

var (
	ErrInvalidBitOrder = errors.New("invalid bit order")
	ErrInvalidFirstBit = errors.New("first bit number must be 0 or 1")
)

func PermuteBits(value []byte, pBlock []int, order BitOrder, firstBit int) ([]byte, error) {
	if firstBit != 0 && firstBit != 1 {
		return nil, ErrInvalidFirstBit
	}
	if order != LeastSignificantFirst && order != MostSignificantFirst {
		return nil, ErrInvalidBitOrder
	}

	bitCount := len(value) * 8
	if len(pBlock) != bitCount {
		return nil, fmt.Errorf("the P-block length is %d, want %d", len(pBlock), bitCount)
	}

	seen := make([]bool, bitCount)
	for i, srcNumber := range pBlock {
		idx := srcNumber - firstBit
		if idx < 0 || idx >= bitCount {
			return nil, fmt.Errorf("the P-block[%d] = %d: bit number is outside [%d, %d]", i, srcNumber, firstBit, firstBit+bitCount-1)
		}
		if seen[idx] {
			return nil, fmt.Errorf("the P-block[%d] = %d: source bit is repeated", i, srcNumber)
		}
		seen[idx] = true
	}

	result := make([]byte, len(value))
	for idx, srcNumber := range pBlock {
		source := srcNumber - firstBit
		if bitAt(value, source, order) == 1 {
			setBit(result, idx, order)
		}
	}

	return result, nil
}

func bitAt(value []byte, position int, order BitOrder) byte {
	byteIndex, bitIndex := bitAddress(position, order)
	return (value[byteIndex] >> bitIndex) & 1
}

func setBit(value []byte, position int, order BitOrder) {
	byteIndex, bitIndex := bitAddress(position, order)
	value[byteIndex] |= 1 << bitIndex
}

func bitAddress(position int, order BitOrder) (byteIndex int, bitIndex uint) {
	byteIndex = position / 8
	bitIndex = uint(position % 8)
	if order == MostSignificantFirst {
		bitIndex = 7 - bitIndex
	}
	return byteIndex, bitIndex
}
