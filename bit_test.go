package jtt

import "testing"

func TestGetBitByte(t *testing.T) {
	var v byte = 0b0101_0010 // bits set at 1, 4, 6
	cases := []struct {
		offset int
		expect bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, false},
		{4, true},
		{5, false},
		{6, true},
		{7, false},
		{8, false}, // out of range
	}

	for _, c := range cases {
		got := GetBitByte(v, c.offset)
		if got != c.expect {
			t.Fatalf("GetBitByte offset=%d: expect %v, got %v", c.offset, c.expect, got)
		}
	}
}

func TestSetBitByte_SetAndClear(t *testing.T) {
	var v byte = 0
	// set a few bits
	SetBitByte(&v, 0, true)
	SetBitByte(&v, 3, true)
	SetBitByte(&v, 7, true)
	if v != 0b1000_1001 {
		// 1<<0 + 1<<3 + 1<<7
		t.Fatalf("SetBitByte set failed: got %08b", v)
	}
	// clear them
	SetBitByte(&v, 3, false)
	SetBitByte(&v, 7, false)
	if v != 0b0000_0001 {
		t.Fatalf("SetBitByte clear failed: got %08b", v)
	}
}

func TestSetBitByte_OutOfRangeNoChange(t *testing.T) {
	var v byte = 0b0011_1100
	SetBitByte(&v, 8, true)  // out of range for byte
	SetBitByte(&v, 9, false) // out of range for byte
	if v != 0b0011_1100 {
		// should remain unchanged
		t.Fatalf("SetBitByte out-of-range should not change value, got %08b", v)
	}
}

func TestGetBitUint16(t *testing.T) {
	var v uint16 = 0
	v |= 1 << 0
	v |= 1 << 5
	v |= 1 << 12
	v |= 1 << 15

	cases := []struct {
		offset int
		expect bool
	}{
		{0, true},
		{1, false},
		{5, true},
		{12, true},
		{14, false},
		{15, true},
		{16, false}, // out of range
	}
	for _, c := range cases {
		got := GetBitUint16(v, c.offset)
		if got != c.expect {
			t.Fatalf("GetBitUint16 offset=%d: expect %v, got %v", c.offset, c.expect, got)
		}
	}
}

func TestSetBitUint16_SetAndClear(t *testing.T) {
	var v uint16 = 0
	SetBitUint16(&v, 0, true)
	SetBitUint16(&v, 10, true)
	SetBitUint16(&v, 15, true)
	if v != (uint16(1)<<0)|(uint16(1)<<10)|(uint16(1)<<15) {
		t.Fatalf("SetBitUint16 set failed: got %016b", v)
	}
	SetBitUint16(&v, 10, false)
	if v != (uint16(1)<<0)|(uint16(1)<<15) {
		t.Fatalf("SetBitUint16 clear failed: got %016b", v)
	}
}

func TestSetBitUint16_OutOfRangeNoChange(t *testing.T) {
	var v uint16 = (1 << 3) | (1 << 9)
	SetBitUint16(&v, 16, true) // out of range for uint16
	SetBitUint16(&v, 20, false)
	if v != (1<<3)|(1<<9) {
		t.Fatalf("SetBitUint16 out-of-range should not change value, got %016b", v)
	}
}

func TestGetBitUint32(t *testing.T) {
	var v uint32 = 0
	v |= 1 << 0
	v |= 1 << 16
	v |= 1 << 31

	cases := []struct {
		offset int
		expect bool
	}{
		{0, true},
		{1, false},
		{16, true},
		{30, false},
		{31, true},
		{32, false}, // out of range
	}
	for _, c := range cases {
		got := GetBitUint32(v, c.offset)
		if got != c.expect {
			t.Fatalf("GetBitUint32 offset=%d: expect %v, got %v", c.offset, c.expect, got)
		}
	}
}

func TestSetBitUint32_SetAndClear(t *testing.T) {
	var v uint32 = 0
	SetBitUint32(&v, 0, true)
	SetBitUint32(&v, 16, true)
	SetBitUint32(&v, 31, true)
	if v != (uint32(1)<<0)|(uint32(1)<<16)|(uint32(1)<<31) {
		t.Fatalf("SetBitUint32 set failed: got %032b", v)
	}
	SetBitUint32(&v, 16, false)
	if v != (uint32(1)<<0)|(uint32(1)<<31) {
		t.Fatalf("SetBitUint32 clear failed: got %032b", v)
	}
}
