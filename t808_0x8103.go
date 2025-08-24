package jtt

import "fmt"

// T808_0x8103 设置终端参数
// 消息体：
// 参数总数N(1 byte) + N个参数项
// 每个参数项：参数ID(DWORD) + 参数长度(BYTE) + 参数值
// 参数ID和值的定义参考 params.go 中的 ParamID 与 Param
// 编码/解码风格与 0x0104(查询终端参数应答)保持一致。
type T808_0x8103 struct {
	// 参数项列表
	Params []*Param
}

func (entity *T808_0x8103) MsgID() MsgID { return MsgT808_0x8103 }

func (entity *T808_0x8103) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入参数个数
	writer.WriteByte(byte(len(entity.Params)))

	// 写入参数列表
	for _, p := range entity.Params {
		// 参数ID
		writer.WriteUint32(uint32(p.Id))
		// 参数长度
		writer.WriteByte(byte(len(p.Data)))
		// 参数值
		writer.Write(p.Data)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8103) Decode(data []byte) (int, error) {
	if len(data) < 1 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	// 读取参数个数
	n, err := reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read paramNums: %w", err)
	}

	entity.Params = make([]*Param, 0, n)
	for i := 0; i < int(n); i++ {
		id, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Id: %w", i, err)
		}
		ln, err := reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Size: %w", i, err)
		}
		val, err := reader.Read(int(ln))
		if err != nil {
			return 0, fmt.Errorf("read param[%d].Value: %w", i, err)
		}
		entity.Params = append(entity.Params, &Param{
			Id:   ParamID(id),
			Data: val,
		})
	}
	return len(data) - reader.Len(), nil
}
