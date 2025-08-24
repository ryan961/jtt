package jtt

import "fmt"

// T808_0x8202 临时位置跟踪控制
type T808_0x8202 struct {
	// 时间间隔，单位：秒。0 表示停止跟踪
	IntervalSec uint16
	// 位置跟踪有效期，单位：秒
	DurationSec uint32
}

func (m *T808_0x8202) MsgID() MsgID { return MsgT808_0x8202 }

func (m *T808_0x8202) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteWord(m.IntervalSec)
	w.WriteDWord(m.DurationSec)
	return w.Bytes(), nil
}

func (m *T808_0x8202) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, fmt.Errorf("invalid body for T808_0x8202: %w (need >=6 bytes, got %d)", ErrInvalidBody, len(data))
	}
	r := NewReader(data)
	var err error
	if m.IntervalSec, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read IntervalSec: %w", err)
	}
	if m.DurationSec, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read DurationSec: %w", err)
	}
	return len(data) - r.Len(), nil
}
