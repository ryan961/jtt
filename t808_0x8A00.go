package jtt

import (
	"fmt"
)

// T808_0x8A00 平台RSA公钥
type T808_0x8A00 struct {
	// 平台RSA公钥{e,n}中的e
	E uint32
	// RSA公钥{e,n}中的n
	N [128]byte
}

func (entity *T808_0x8A00) MsgID() MsgID { return MsgT808_0x8A00 }

func (entity *T808_0x8A00) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入e
	writer.WriteDWord(entity.E)

	// 写入n (128字节)
	writer.Write(entity.N[:])

	return writer.Bytes(), nil
}

func (entity *T808_0x8A00) Decode(data []byte) (int, error) {
	if len(data) < 132 { // 4 bytes for E + 128 bytes for N
		return 0, fmt.Errorf("invalid data length: %d, expected at least 132", len(data))
	}

	r := NewReader(data)
	var err error

	// 读取e
	if entity.E, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read E: %w", err)
	}

	// 读取n (128字节)
	var nBytes []byte
	if nBytes, err = r.Read(128); err != nil {
		return 0, fmt.Errorf("read N: %w", err)
	}
	copy(entity.N[:], nBytes)

	return len(data) - r.Len(), nil
}
