package jtt

import "fmt"

// T808_0x8700 行驶记录数据采集命令
type T808_0x8700 struct {
	// 命令字，应符合GB/T 19056中相关要求
	Command byte
	// 数据块，内容格式应符合GB/T 19056要求的完整数据包，可为空
	DataBlock []byte
}

func (m *T808_0x8700) MsgID() MsgID { return MsgT808_0x8700 }

func (m *T808_0x8700) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.Command)
	if len(m.DataBlock) > 0 {
		w.Write(m.DataBlock)
	}
	return w.Bytes(), nil
}

func (m *T808_0x8700) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.Command, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read command: %w", err)
	}

	// 剩余所有字节作为数据块
	remaining := r.Len()
	if remaining > 0 {
		if m.DataBlock, err = r.Read(remaining); err != nil {
			return 0, fmt.Errorf("read data block: %w", err)
		}
	}

	return len(data) - r.Len(), nil
}
