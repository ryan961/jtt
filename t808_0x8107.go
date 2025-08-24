package jtt

// T808_0x8107 查询终端属性（消息体为空）
type T808_0x8107 struct{}

func (m *T808_0x8107) MsgID() MsgID { return MsgT808_0x8107 }

func (m *T808_0x8107) Encode() ([]byte, error) { return []byte{}, nil }

func (m *T808_0x8107) Decode([]byte) (int, error) { return 0, nil }
