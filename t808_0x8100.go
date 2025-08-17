package jtt

import "fmt"

// T808_0x8100 终端注册应答
type T808_0x8100 struct {
	// 应答流水号，对应的终端注册消息的流水号
	ReplyMsgSerialNo uint16 `json:"replyMsgSerialNo"`
	// 结果，0：成功，1：车辆已被注册，2：数据库中无该车辆，3：终端已被注册，4：数据库中无该终端
	Result byte `json:"result"`
	// 鉴权码，注册成功时才有该字段
	AuthCode string `json:"authCode"`
}

func (entity *T808_0x8100) MsgID() MsgID { return MsgT808_0x8100 }

func (entity *T808_0x8100) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteWord(entity.ReplyMsgSerialNo)
	writer.WriteByte(entity.Result)
	if entity.Result == 0 {
		if err := writer.WriteString(entity.AuthCode); err != nil {
			return nil, err
		}
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8100) Decode(data []byte) (int, error) {
	if len(data) < 3 { // WORD + BYTE
		return 0, fmt.Errorf("invalid body for T808_0x8100: %w (need >=3 bytes, got %d)", ErrInvalidBody, len(data))
	}
	reader := NewReader(data)
	var err error
	if entity.ReplyMsgSerialNo, err = reader.ReadWord(); err != nil {
		return 0, fmt.Errorf("read ReplyMsgSerialNo: %w", err)
	}
	if entity.Result, err = reader.ReadByte(); err != nil {
		return 0, fmt.Errorf("read Result: %w", err)
	}
	if entity.Result == 0 {
		if entity.AuthCode, err = reader.ReadString(); err != nil {
			return 0, fmt.Errorf("read AuthCode: %w", err)
		}
	}
	return len(data) - reader.Len(), nil
}
