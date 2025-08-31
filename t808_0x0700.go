package jtt

import "fmt"

// T808_0x0700 行驶记录数据上传
type T808_0x0700 struct {
	// 应答流水号，对应的行驶记录数据采集命令消息的流水号
	ReplyMsgSerialNo uint16
	// 命令字，对应平台发出的命令字
	Command byte
	// 数据块，内容格式见GB/T 19056中相关内容
	DataBlock []byte
}

func (m *T808_0x0700) MsgID() MsgID { return MsgT808_0x0700 }

func (m *T808_0x0700) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteWord(m.ReplyMsgSerialNo)
	w.WriteByte(m.Command)
	if len(m.DataBlock) > 0 {
		w.Write(m.DataBlock)
	}
	return w.Bytes(), nil
}

func (m *T808_0x0700) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.ReplyMsgSerialNo, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read reply msg serial number: %w", err)
	}

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
