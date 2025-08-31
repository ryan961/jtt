package jtt

import "fmt"

// T808_0x8400 电话回拨
type T808_0x8400 struct {
	// 标志 0:普通通话; 1:监听
	Flag byte
	// 电话号码，最长20字节
	Phone string
}

func (m *T808_0x8400) MsgID() MsgID { return MsgT808_0x8400 }

func (m *T808_0x8400) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.Flag)

	if len(m.Phone) > 20 {
		return nil, fmt.Errorf("phone too long: %d", len(m.Phone))
	}

	if len(m.Phone) > 0 {
		if err := w.WriteString(m.Phone); err != nil {
			return nil, fmt.Errorf("write phone: %w", err)
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x8400) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error
	if m.Flag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read Flag: %w", err)
	}
	if r.Len() > 0 {
		if m.Phone, err = r.ReadString(); err != nil {
			return 0, fmt.Errorf("read Phone: %w", err)
		}
	}
	return len(data) - r.Len(), nil
}
