package jtt

// T808_0x0004 查询服务器时间请求
type T808_0x0004 struct{}

func (entity *T808_0x0004) MsgID() MsgID {
	return MsgT808_0x0004
}

func (entity *T808_0x0004) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (entity *T808_0x0004) Decode(data []byte) (int, error) {
	return len(data), nil
}
