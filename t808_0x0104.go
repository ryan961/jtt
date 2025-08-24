package jtt

import "fmt"

// T808_0x0104 查询终端参数应答
type T808_0x0104 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 参数项列表
	Params []*Param
}

func (entity *T808_0x0104) MsgID() MsgID {
	return MsgT808_0x0104
}

func (entity *T808_0x0104) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入消息序列号
	writer.WriteUint16(entity.ReplyMsgSerialNo)

	// 写入参数个数
	writer.WriteByte(byte(len(entity.Params)))

	// 写入参数列表
	for _, param := range entity.Params {
		// 写入参数ID
		writer.WriteUint32(uint32(param.Id))

		// 写入参数长度
		writer.WriteByte(byte(len(param.Data)))

		// 写入参数数据
		writer.Write(param.Data)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x0104) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, fmt.Errorf("invalid body for T808_0x0104: %w (need >=3 bytes, got %d)", ErrInvalidBody, len(data))
	}
	reader := NewReader(data)

	// 读取消息序列号
	responseMsgSerialNo, err := reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read ReplyMsgSerialNo: %w", err)
	}
	entity.ReplyMsgSerialNo = responseMsgSerialNo

	// 读取参数个数
	paramNums, err := reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read paramNums: %w", err)
	}

	// 读取参数信息
	entity.Params = make([]*Param, 0, paramNums)
	for i := 0; i < int(paramNums); i++ {
		// 读取参数ID
		id, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Id: %w", i, err)
		}

		// 读取数据长度
		size, err := reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Size: %w", i, err)
		}

		// 读取数据内容
		value, err := reader.Read(int(size))
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Value: %w", i, err)
		}
		entity.Params = append(entity.Params, &Param{
			Id:   ParamID(id),
			Data: value,
		})
	}
	return len(data) - reader.Len(), nil
}
