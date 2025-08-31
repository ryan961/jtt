package jtt

// T808_0x8E12 驾驶员身份库信息查询
// 消息体为空
type T808_0x8E12 struct {
}

func (msg *T808_0x8E12) MsgID() MsgID {
	return MsgT808_0x8E12
}

func (msg *T808_0x8E12) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (msg *T808_0x8E12) Decode(data []byte) (int, error) {
	// 消息体为空
	return len(data), nil
}
