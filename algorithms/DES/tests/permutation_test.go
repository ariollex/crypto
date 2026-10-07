package tests

import (
	"bytes"
	"crypto/algorithms/DES"
	"errors"
	"strings"
	"testing"
)

func TestPermuteBits(t *testing.T) {
	input := []byte{0b_0001_0010, 0b_1011_0000}

	tests := []struct {
		name     string
		pBlock   []int
		order    DES.BitOrder
		firstBit int
		want     []byte
	}{
		{
			name:     "least significant first with numbering from zero",
			pBlock:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0},
			order:    DES.LeastSignificantFirst,
			firstBit: 0,
			want:     []byte{0b_0000_1001, 0b_0101_1000},
		},
		{
			name:     "most significant first with numbering from one",
			pBlock:   []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 1},
			order:    DES.MostSignificantFirst,
			firstBit: 1,
			want:     []byte{0b_0010_0101, 0b_0110_0000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DES.PermuteBits(input, tt.pBlock, tt.order, tt.firstBit)
			if err != nil {
				t.Fatalf("PermuteBits returned an error: %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Fatalf("PermuteBits() = % X, want % X", got, tt.want)
			}
		})
	}
}

func TestPermuteBitsDoesNotModifyInput(t *testing.T) {
	input := []byte{0b_0001_0010}
	want := make([]byte, len(input))
	copy(want, input)

	_, err := DES.PermuteBits(input, []int{7, 6, 5, 4, 3, 2, 1, 0}, DES.LeastSignificantFirst, 0)
	if err != nil {
		t.Fatalf("PermuteBits returned an error: %v", err)
	}
	if !bytes.Equal(input, want) {
		t.Fatalf("input was changed to % X, want % X", input, want)
	}
}

func TestPermuteBitsRejectsInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		pBlock   []int
		order    DES.BitOrder
		firstBit int
		wantErr  error
		contains string
	}{
		{
			name:     "wrong first bit number",
			pBlock:   []int{0, 1, 2, 3, 4, 5, 6, 7},
			order:    DES.LeastSignificantFirst,
			firstBit: 2,
			wantErr:  DES.ErrInvalidFirstBit,
		},
		{
			name:     "unknown bit order",
			pBlock:   []int{0, 1, 2, 3, 4, 5, 6, 7},
			order:    DES.BitOrder(67),
			firstBit: 0,
			wantErr:  DES.ErrInvalidBitOrder,
		},
		{
			name:     "wrong P block size",
			pBlock:   []int{0, 1},
			order:    DES.LeastSignificantFirst,
			firstBit: 0,
			contains: "length",
		},
		{
			name:     "repeated source bit",
			pBlock:   []int{0, 1, 2, 3, 4, 5, 6, 6},
			order:    DES.LeastSignificantFirst,
			firstBit: 0,
			contains: "repeated",
		},
		{
			name:     "source bit outside the value",
			pBlock:   []int{0, 1, 2, 3, 4, 5, 6, 8},
			order:    DES.LeastSignificantFirst,
			firstBit: 0,
			contains: "outside",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DES.PermuteBits([]byte{0}, tt.pBlock, tt.order, tt.firstBit)
			if err == nil {
				t.Fatal("PermuteBits returned nil, want error")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.contains != "" && !strings.Contains(err.Error(), tt.contains) {
				t.Fatalf("error = %q, want it to contain %q", err, tt.contains)
			}
		})
	}
}
