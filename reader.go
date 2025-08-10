package jtt

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Reader struct {
	d []byte
	r *bytes.Reader
}

func NewReader(data []byte) *Reader {
	return &Reader{d: data, r: bytes.NewReader(data)}
}

func (reader *Reader) Len() int {
	return reader.r.Len()
}

// Read 对应 JT808 类型 BYTE[n].
func (reader *Reader) Read(size ...int) ([]byte, error) {
	num := reader.r.Len()
	if len(size) > 0 {
		num = size[0]
	}

	if num > reader.r.Len() {
		return nil, io.ErrUnexpectedEOF
	}

	curr := len(reader.d) - reader.r.Len()
	buf := reader.d[curr : curr+num]
	_, _ = reader.r.Seek(int64(num), io.SeekCurrent)
	return buf, nil
}

// ReadByte 对应 JT808 类型 BYTE.
func (reader *Reader) ReadByte() (byte, error) {
	return reader.r.ReadByte()
}

// ReadUint16 对应 JT808 类型 WORD.
func (reader *Reader) ReadUint16() (uint16, error) {
	if reader.r.Len() < 2 {
		return 0, io.ErrUnexpectedEOF
	}

	var buf [2]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.BigEndian.Uint16(buf[:]), nil
}

// ReadWord 对应 JT808 类型 WORD.
func (reader *Reader) ReadWord() (uint16, error) {
	return reader.ReadUint16()
}

// ReadUint32 对应 JT808 类型 DWORD.
func (reader *Reader) ReadUint32() (uint32, error) {
	if reader.r.Len() < 4 {
		return 0, io.ErrUnexpectedEOF
	}

	var buf [4]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.BigEndian.Uint32(buf[:]), nil
}

// ReadDWord 对应 JT808 类型 DWORD.
func (reader *Reader) ReadDWord() (uint32, error) {
	return reader.ReadUint32()
}

// ReadUint64 对应 JT808 类型 BYTE[8].
func (reader *Reader) ReadUint64() (uint64, error) {
	if reader.r.Len() < 8 {
		return 0, io.ErrUnexpectedEOF
	}
	var buf [8]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if n != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return binary.BigEndian.Uint64(buf[:]), nil
}

// ReadBcdTime 对应 JT808 类型 BCD.
func (reader *Reader) ReadBcdTime() (time.Time, error) {
	if reader.r.Len() < 6 {
		return time.Time{}, io.ErrUnexpectedEOF
	}

	var buf [6]byte
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return time.Time{}, err
	}
	if n != len(buf) {
		return time.Time{}, io.ErrUnexpectedEOF
	}
	return FromBCDTime(buf[:])
}

// ReadBcd 对应 JT808 类型 BCD.
func (reader *Reader) ReadBcd(n int) (string, error) {
	if reader.r.Len() < n {
		return "", io.ErrUnexpectedEOF
	}

	var buf = make([]byte, n)
	n, err := reader.r.Read(buf[:])
	if err != nil {
		return "", err
	}
	if n != len(buf) {
		return "", io.ErrUnexpectedEOF
	}

	return BcdToString(buf), nil
}

// ReadString 对应 JT808 类型 STRING.
func (reader *Reader) ReadString(size ...int) (string, error) {
	data, err := reader.Read(size...)
	if err != nil {
		return "", err
	}

	text, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), simplifiedchinese.GB18030.NewDecoder()))
	if err != nil {
		return "", err
	}
	return BytesToString(text), nil
}
