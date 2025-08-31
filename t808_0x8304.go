package jtt

import "fmt"

// T808_0x8304 信息服务
type T808_0x8304 struct {
	// 信息类型
	InfoType byte
	// 信息内容，经 GBK 编码
	Content string
}

func (m *T808_0x8304) MsgID() MsgID { return MsgT808_0x8304 }

func (m *T808_0x8304) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.InfoType)
	ln, err := GB18030Length(m.Content)
	if err != nil {
		return nil, fmt.Errorf("get content length: %w", err)
	}
	w.WriteWord(uint16(ln))
	if ln > 0 {
		if err := w.WriteString(m.Content); err != nil {
			return nil, fmt.Errorf("write content: %w", err)
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x8304) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error
	if m.InfoType, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read InfoType: %w", err)
	}
	var l uint16
	if l, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read length: %w", err)
	}
	if l > 0 {
		if m.Content, err = r.ReadString(int(l)); err != nil {
			return 0, fmt.Errorf("read content: %w", err)
		}
	}
	return len(data) - r.Len(), nil
}
