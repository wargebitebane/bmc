package complement

import (
	"testing"
)

func TestTwos(t *testing.T) {
	tests := []struct {
		in   [2]byte
		bits uint8
		want int16
	}{
		// we test the following values for each of 4, 8, 10 and 16 bit numbers:
		//   lowest
		//   lowest + 1
		//   lowest + 2
		//   -2
		//   -1
		//   0
		//   1
		//   2
		//   highest - 2
		//   highest - 1
		//   highest

		// 4-bit
		{
			[...]byte{0, 0b00001000},
			4,
			-8,
		},
		{
			[...]byte{0, 0b00001001},
			4,
			-7,
		},
		{
			[...]byte{0, 0b00001010},
			4,
			-6,
		},
		{
			[...]byte{0, 0b00001110},
			4,
			-2,
		},
		{
			[...]byte{0, 0b00001111},
			4,
			-1,
		},
		{
			[...]byte{0, 0},
			4,
			0,
		},
		{
			[...]byte{0, 0b00000001},
			4,
			1,
		},
		{
			[...]byte{0, 0b00000010},
			4,
			2,
		},
		{
			[...]byte{0, 0b00000101},
			4,
			5,
		},
		{
			[...]byte{0, 0b00000110},
			4,
			6,
		},
		{
			[...]byte{0, 0b00000111},
			4,
			7,
		},

		// 8-bit
		{
			[...]byte{0, 0b10000000},
			8,
			-128,
		},
		{
			[...]byte{0, 0b10000001},
			8,
			-127,
		},
		{
			[...]byte{0, 0b10000010},
			8,
			-126,
		},
		{
			[...]byte{0, 0b11111110},
			8,
			-2,
		},
		{
			[...]byte{0, 0b11111111},
			8,
			-1,
		},
		{
			[...]byte{0, 0},
			8,
			0,
		},
		{
			[...]byte{0, 0b00000001},
			8,
			1,
		},
		{
			[...]byte{0, 0b00000010},
			8,
			2,
		},
		{
			[...]byte{0, 0b01111101},
			8,
			125,
		},
		{
			[...]byte{0, 0b01111110},
			8,
			126,
		},
		{
			[...]byte{0, 0b01111111},
			8,
			127,
		},

		// 10-bit
		{
			[...]byte{0b00000010, 0b00000000},
			10,
			-512,
		},
		{
			[...]byte{0b00000010, 0b00000001},
			10,
			-511,
		},
		{
			[...]byte{0b00000010, 0b00000010},
			10,
			-510,
		},
		{
			[...]byte{0b00000011, 0b11111110},
			10,
			-2,
		},
		{
			[...]byte{0b00000011, 0b11111111},
			10,
			-1,
		},
		{
			[...]byte{0, 0},
			10,
			0,
		},
		{
			[...]byte{0, 0b00000001},
			10,
			1,
		},
		{
			[...]byte{0, 0b00000010},
			10,
			2,
		},
		{
			[...]byte{0b00000001, 0b11111101},
			10,
			509,
		},
		{
			[...]byte{0b00000001, 0b11111110},
			10,
			510,
		},
		{
			[...]byte{0b00000001, 0b11111111},
			10,
			511,
		},

		// 16-bit
		{
			[...]byte{0b10000000, 0b00000000},
			16,
			-32768,
		},
		{
			[...]byte{0b10000000, 0b00000001},
			16,
			-32767,
		},
		{
			[...]byte{0b10000000, 0b00000010},
			16,
			-32766,
		},
		{
			[...]byte{0b11111111, 0b11111110},
			16,
			-2,
		},
		{
			[...]byte{0b11111111, 0b11111111},
			16,
			-1,
		},
		{
			[...]byte{0, 0},
			16,
			0,
		},
		{
			[...]byte{0, 0b00000001},
			16,
			1,
		},
		{
			[...]byte{0, 0b00000010},
			16,
			2,
		},
		{
			[...]byte{0b01111111, 0b11111101},
			16,
			32765,
		},
		{
			[...]byte{0b01111111, 0b11111110},
			16,
			32766,
		},
		{
			[...]byte{0b01111111, 0b11111111},
			16,
			32767,
		},
	}
	for _, test := range tests {
		got := Twos(test.in, test.bits)
		if got != test.want {
			t.Errorf("Twos(%v, %v) = %v, want %v", test.in, test.bits, got, test.want)
		}
	}
}