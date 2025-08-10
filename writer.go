package jtt

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Writer struct {
	b *bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{b: bytes.NewBuffer(nil)}
}

func (writer *Writer) Bytes() []byte {
	return writer.b.Bytes()
}

// Write 对应 JT808 类型 BYTE[n]，可以设定补充的 size 位数，超过 len(p) 时只会补充 p 的前 size 位，不足位时会在末尾补充 0.
func (writer *Writer) Write(p []byte, size ...int) *Writer {
	if len(size) == 0 {
		writer.b.Write(p)
		return writer
	}

	if len(p) >= size[0] {
		writer.b.Write(p[:size[0]])
	} else {
		writer.b.Write(p)
		end := size[0] - len(p)
		for i := 0; i < end; i++ {
			writer.b.WriteByte(0)
		}
	}
	return writer
}

// WriteByte 对应 JT808 类型 BYTE.
//
//nolint:govet,revive,stylecheck
func (writer *Writer) WriteByte(b byte) *Writer {
	_ = writer.b.WriteByte(b)
	return writer
}

// WriteUint16 对应 JT808 类型 WORD.
func (writer *Writer) WriteUint16(n uint16) *Writer {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], n)
	writer.b.Write(buf[:])
	return writer
}

// WriteWord 对应 JT808 类型 WORD.
func (writer *Writer) WriteWord(n uint16) *Writer {
	return writer.WriteUint16(n)
}

// WriteUint32 对应 JT808 类型 DWORD.
func (writer *Writer) WriteUint32(n uint32) *Writer {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], n)
	writer.b.Write(buf[:])
	return writer
}

// WriteDWord 对应 JT808 类型 DWORD.
func (writer *Writer) WriteDWord(n uint32) *Writer {
	return writer.WriteUint32(n)
}

// WriteUint64 对应 JT808 类型 BYTE[8].
func (writer *Writer) WriteUint64(n uint64) *Writer {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], n)
	writer.b.Write(buf[:])
	return writer
}

// WriteBcdTime 输入 time.Time, 转换为 JT808 协议定义的时间 format.
func (writer *Writer) WriteBcdTime(t time.Time) *Writer {
	writer.b.Write(ToBCDTime(t))
	return writer
}

// WriteBcd 对应 JT808 类型 BCD[n].
func (writer *Writer) WriteBcd(str string, n int) *Writer {
	writer.b.Write(StringToBCD(str, n))
	return writer
}

// WriteString 对应 JT808 类型 STRING.
func (writer *Writer) WriteString(str string, size ...int) error {
	reader := bytes.NewReader([]byte(str))
	data, err := io.ReadAll(transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder()))
	if err != nil {
		return err
	}
	writer.Write(data, size...)
	return nil
}
