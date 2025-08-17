package jtt

import "fmt"

// T808_0x8003 补传分包请求
type T808_0x8003 struct {
	// 原始消息流水号，对应要求补传的原始消息第一包的消息流水号
	OriginalMsgSerialNo uint16 `json:"originalMsgSerialNo"`
	// 重传包总数
	TotalCount uint16 `json:"totalCount"`
	// 重传包ID列表，重传包序号顺序排列	，BYTE[2*n]，如"包 ID1 包 ID2......包 IDn",n为重传包的总数
	PackageIDList []uint16 `json:"packageIDList"`
}

func (entity *T808_0x8003) MsgID() MsgID { return MsgT808_0x8003 }

func (entity *T808_0x8003) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入原始消息流水号
	writer.WriteUint16(entity.OriginalMsgSerialNo)

	// 写入重传包总数
	writer.WriteUint16(entity.TotalCount)

	// 写入重传包 ID 列表
	for _, id := range entity.PackageIDList {
		writer.WriteUint16(id)
	}
	return writer.Bytes(), nil
}

func (entity *T808_0x8003) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("invalid body for T808_0x8003: %w (need >=4 bytes, got %d)", ErrInvalidBody, len(data))
	}
	reader := NewReader(data)

	// 读取原始消息流水号
	var err error
	entity.OriginalMsgSerialNo, err = reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read OriginalMsgSerialNo: %w", err)
	}

	// 读取重传包总数
	entity.TotalCount, err = reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read TotalCount: %w", err)
	}

	// 读取重传包 ID 列表
	entity.PackageIDList = make([]uint16, entity.TotalCount)
	for i := uint16(0); i < entity.TotalCount; i++ {
		entity.PackageIDList[i], err = reader.ReadUint16()
		if err != nil {
			return 0, fmt.Errorf("read PackageIDList[%d]: %w", i, err)
		}
	}
	return len(data) - reader.Len(), nil
}
