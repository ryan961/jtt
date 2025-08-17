package jtt

import "fmt"

// T808_0x0001 终端通用应答
type T808_0x0001 struct {
	ReplyMsgSerialNo uint16 `json:"replyMsgSerialNo"` // 应答流水号，对应的终端消息的流水号
	ReplyMsgID       MsgID  `json:"replyMsgID"`       // 应答ID,对应的终端消息的ID
	Result           byte   `json:"result"`           // 结果，0;成功/确认;1:失败;2;消息有误;3:不支持
}

func (entity *T808_0x0001) MsgID() MsgID {
	return MsgT808_0x0001
}

func (entity *T808_0x0001) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入消息序列号
	writer.WriteUint16(entity.ReplyMsgSerialNo)

	// 写入响应消息ID
	writer.WriteUint16(uint16(entity.ReplyMsgID))

	// 写入响应结果
	writer.WriteByte(entity.Result)
	return writer.Bytes(), nil
}

func (entity *T808_0x0001) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, fmt.Errorf("invalid body for T808_0x0001: %w (need >=5 bytes, got %d)", ErrInvalidBody, len(data))
	}
	reader := NewReader(data)

	// 读取消息序列号
	var err error
	entity.ReplyMsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read ReplyMsgSerialNo: %w", err)
	}

	// 读取响应消息ID
	id, err := reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read ReplyMsgID: %w", err)
	}
	entity.ReplyMsgID = MsgID(id)

	// 读取响应结果
	result, err := reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read Result: %w", err)
	}
	entity.Result = result
	return len(data) - reader.Len(), nil
}
