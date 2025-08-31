package jtt

import (
	"fmt"
)

// T808_0x8900 数据下行透传
type T808_0x8900 struct {
	// 透传消息类型
	// 0x00: GNSS模块详细定位数据
	// 0x0B: 道路运输证IC卡信息
	// 0x41: 串口1透传
	// 0x42: 串口2透传
	// 0xF0-0xFF: 用户自定义透传
	TransparentMsgType uint8
	// 透传消息内容
	TransparentMsgContent []byte
}

func (entity *T808_0x8900) MsgID() MsgID { return MsgT808_0x8900 }

func (entity *T808_0x8900) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入透传消息类型
	writer.WriteByte(entity.TransparentMsgType)

	// 写入透传消息内容
	if len(entity.TransparentMsgContent) > 0 {
		writer.Write(entity.TransparentMsgContent)
	}

	return writer.Bytes(), nil
}

func (entity *T808_0x8900) Decode(data []byte) (int, error) {
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
