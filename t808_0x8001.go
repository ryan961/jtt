package jtt

import "fmt"

// T808_0x8001 平台通用应答
type T808_0x8001 struct {
	ReplyMsgSerialNo uint16 `json:"replyMsgSerialNo"` // 应答流水号，对应的终端消息的流水号
	ReplyMsgID       MsgID  `json:"replyMsgID"`       // 应答ID,对应的终端消息的ID
	Result           byte   `json:"result"`           // 结果，0;成功/确认;1:失败;2;消息有误;3:不支持;4:报警处理确认(2013、2019新增)
}

func (entity *T808_0x8001) MsgID() MsgID { return MsgT808_0x8001 }

func (entity *T808_0x8001) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteWord(entity.ReplyMsgSerialNo)

	// 写入应答ID
	writer.WriteWord(uint16(entity.ReplyMsgID))

	// 写入结果
	writer.WriteByte(entity.Result)

	return writer.Bytes(), nil
}

func (entity *T808_0x8001) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, fmt.Errorf("invalid body for T808_0x8001: %w (need >=5 bytes, got %d)", ErrInvalidBody, len(data))
	}
	reader := NewReader(data)

	// 读取应答流水号
	var err error
	entity.ReplyMsgSerialNo, err = reader.ReadWord()
	if err != nil {
		return 0, fmt.Errorf("read ReplyMsgSerialNo: %w", err)
	}

	// 读取应答ID
	msgID, err := reader.ReadWord()
	if err != nil {
		return 0, fmt.Errorf("read ReplyMsgID: %w", err)
	}
	entity.ReplyMsgID = MsgID(msgID)

	// 读取结果
	entity.Result, err = reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read Result: %w", err)
	}

	return len(data) - reader.Len(), nil
}
