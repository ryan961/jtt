package jtt

import (
	"fmt"
)

// T808_0x0901 数据压缩上报
type T808_0x0901 struct {
	// 压缩消息长度
	CompressedMsgLength uint32
	// 压缩消息体（需要压缩的消息经过GZIP压缩算法后的消息）
	CompressedMsgBody []byte
}

func (entity *T808_0x0901) MsgID() MsgID { return MsgT808_0x0901 }

func (entity *T808_0x0901) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入压缩消息长度
	writer.WriteDWord(entity.CompressedMsgLength)

	// 写入压缩消息体
	if len(entity.CompressedMsgBody) > 0 {
		writer.Write(entity.CompressedMsgBody)
	}

	return writer.Bytes(), nil
}

func (entity *T808_0x0901) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}

	r := NewReader(data)
	var err error

	// 读取压缩消息长度
	if entity.CompressedMsgLength, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read compressed msg length: %w", err)
	}

	// 读取压缩消息体
	if r.Len() > 0 {
		if entity.CompressedMsgBody, err = r.Read(); err != nil {
			return 0, fmt.Errorf("read compressed msg body: %w", err)
		}
	}

	return len(data) - r.Len(), nil
}
