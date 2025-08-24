package jtt

import "fmt"

// T808_0x8106 查询指定终端参数（终端应使用 0x0104 消息进行应答）
type T808_0x8106 struct {
	// 参数ID列表
	IDs []ParamID
}

func (entity *T808_0x8106) MsgID() MsgID { return MsgT808_0x8106 }

func (entity *T808_0x8106) Encode() ([]byte, error) {
	writer := NewWriter()
	// 写入参数总数
	writer.WriteByte(byte(len(entity.IDs)))
	// 写入参数ID列表（DWORD each）
	for _, id := range entity.IDs {
		writer.WriteUint32(uint32(id))
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8106) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, fmt.Errorf("invalid body for T808_0x8106: %w (need >=1 bytes, got %d)", ErrInvalidBody, len(data))
	}
	reader := NewReader(data)
	// 读取参数总数
	n, err := reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read paramNums: %w", err)
	}
	entity.IDs = make([]ParamID, 0, int(n))
	for i := 0; i < int(n); i++ {
		v, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Id: %w", i, err)
		}
		entity.IDs = append(entity.IDs, ParamID(v))
	}
	return len(data) - reader.Len(), nil
}
