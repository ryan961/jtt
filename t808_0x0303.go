package jtt

import "fmt"

// T808_0x0303 信息点播/取消
type T808_0x0303 struct {
	// 信息类型
	InfoType byte
	// 点播/取消标志 0:取消, 1:点播
	Flag byte
}

func (m *T808_0x0303) MsgID() MsgID { return MsgT808_0x0303 }

func (m *T808_0x0303) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.InfoType)
	w.WriteByte(m.Flag)
	return w.Bytes(), nil
}

func (m *T808_0x0303) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error
	if m.InfoType, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read InfoType: %w", err)
	}
	if m.Flag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read Flag: %w", err)
	}
	return len(data) - r.Len(), nil
}
