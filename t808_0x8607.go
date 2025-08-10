package jtt

// T808_0x8607 删除路线
type T808_0x8607 struct {
	RouteCount byte     // 路线数量
	RouteIDs   []uint32 // 路线ID列表
}

// MsgID 获取消息ID
func (entity *T808_0x8607) MsgID() MsgID {
	return MsgT808_0x8607
}

// Encode 编码消息
func (entity *T808_0x8607) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入路线数量
	writer.WriteByte(entity.RouteCount)

	// 写入路线ID列表
	if entity.RouteCount == 0 {
		// 删除所有路线
		return writer.Bytes(), nil
	}

	// 写入指定的路线ID列表
	for _, routeID := range entity.RouteIDs {
		writer.WriteUint32(routeID)
	}

	return writer.Bytes(), nil
}

// Decode 解码消息
func (entity *T808_0x8607) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取路线数量
	var err error
	entity.RouteCount, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取路线ID列表
	if entity.RouteCount == 0 {
		// 删除所有路线
		return len(data) - reader.Len(), nil
	}

	// 读取指定的路线ID列表
	entity.RouteIDs = make([]uint32, entity.RouteCount)
	for i := 0; i < int(entity.RouteCount); i++ {
		routeID, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RouteIDs[i] = routeID
	}

	return len(data) - reader.Len(), nil
}
