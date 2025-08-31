package jtt

import "fmt"

// T808_0x0301 事件报告
type T808_0x0301 struct {
	EventID byte
}

func (m *T808_0x0301) MsgID() MsgID { return MsgT808_0x0301 }

func (m *T808_0x0301) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.EventID)
	return w.Bytes(), nil
}

func (m *T808_0x0301) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data)
	var err error
	m.EventID, err = r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read event id: %w", err)
	}
	return len(data) - r.Len(), nil
}
