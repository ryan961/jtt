package jtt

import (
	"fmt"
)

// T808_0x0900 数据上行透传
type T808_0x0900 struct {
	// 透传消息类型
	TransparentMsgType uint8
	// 透传消息内容
	TransparentMsgContent []byte
}

func (entity *T808_0x0900) MsgID() MsgID { return MsgT808_0x0900 }

func (entity *T808_0x0900) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入透传消息类型
	writer.WriteByte(entity.TransparentMsgType)

	// 写入透传消息内容
	if len(entity.TransparentMsgContent) > 0 {
		writer.Write(entity.TransparentMsgContent)
	}

	return writer.Bytes(), nil
}

func (entity *T808_0x0900) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data)
	var err error

	// 读取透传消息类型
	if entity.TransparentMsgType, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read transparent msg type: %w", err)
	}

	// 读取透传消息内容
	if r.Len() > 0 {
		if entity.TransparentMsgContent, err = r.Read(); err != nil {
			return 0, fmt.Errorf("read transparent msg content: %w", err)
		}
	}

	return len(data) - r.Len(), nil
}
