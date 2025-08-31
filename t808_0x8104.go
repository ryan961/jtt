package jtt

// T808_0x8104 查询终端参数
type T808_0x8104 struct{}

func (entity *T808_0x8104) MsgID() MsgID { return MsgT808_0x8104 }

func (entity *T808_0x8104) Encode() ([]byte, error) {
	return []byte{}, nil
}

func (entity *T808_0x8104) Decode(data []byte) (int, error) {
	return len(data), nil
}
