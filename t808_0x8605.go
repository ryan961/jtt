package jtt

// T808_0x8605 删除多边形区域
type T808_0x8605 struct {
	AreaCount byte     // 区域数,0：删除所有区域,不超过 125 个
	AreaIDs   []uint32 // 区域ID列表
}

// MsgID 获取消息ID
func (entity *T808_0x8605) MsgID() MsgID {
	return MsgT808_0x8605
}

// Encode 编码消息
func (entity *T808_0x8605) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入区域数
	writer.WriteByte(entity.AreaCount)

	// 写入区域ID列表
	if entity.AreaCount == 0 {
		// 删除所有区域
		return writer.Bytes(), nil
	}

	// 写入指定的区域ID列表
	for _, areaID := range entity.AreaIDs {
		writer.WriteUint32(areaID)
	}

	return writer.Bytes(), nil
}

// Decode 解码消息
func (entity *T808_0x8605) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取区域数
	var err error
	entity.AreaCount, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取区域ID列表
	if entity.AreaCount == 0 {
		// 删除所有区域
		return len(data) - reader.Len(), nil
	}

	// 读取指定的区域ID列表
	entity.AreaIDs = make([]uint32, entity.AreaCount)
	for i := 0; i < int(entity.AreaCount); i++ {
		areaID, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.AreaIDs[i] = areaID
	}

	return len(data) - reader.Len(), nil
}
