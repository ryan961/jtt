package jtt

import "fmt"

// T808_0x0701 电子运单上报
type T808_0x0701 struct {
	// 电子运单长度
	Length uint32
	// 电子运单内容，电子运单数据包
	Content []byte
}

func (m *T808_0x0701) MsgID() MsgID { return MsgT808_0x0701 }

func (m *T808_0x0701) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteDWord(m.Length)
	if len(m.Content) > 0 {
		w.Write(m.Content)
	}
	return w.Bytes(), nil
}

func (m *T808_0x0701) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.Length, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read length: %w", err)
	}

	// 读取电子运单内容
	if m.Length > 0 {
		if m.Content, err = r.Read(int(m.Length)); err != nil {
			return 0, fmt.Errorf("read content: %w", err)
		}
	}

	return len(data) - r.Len(), nil
}
