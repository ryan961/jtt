package jtt

import "fmt"

// T808_0x8805 单条存储多媒体数据检索上传命令
type T808_0x8805 struct {
	// 多媒体ID，值大于0
	MultimediaID uint32
	// 删除标志
	// 0: 保留
	// 1: 删除
	DeleteFlag byte
}

func (m *T808_0x8805) MsgID() MsgID { return MsgT808_0x8805 }

func (m *T808_0x8805) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteDWord(m.MultimediaID)
	w.WriteByte(m.DeleteFlag)
	return w.Bytes(), nil
}

func (m *T808_0x8805) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.MultimediaID, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read multimedia id: %w", err)
	}

	if m.DeleteFlag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read delete flag: %w", err)
	}

	return len(data) - r.Len(), nil
}
