package jtt

// T808_0x8702 上报驾驶员身份信息请求
type T808_0x8702 struct{}

func (m *T808_0x8702) MsgID() MsgID { return MsgT808_0x8702 }

func (m *T808_0x8702) Encode() ([]byte, error) { return []byte{}, nil }

func (m *T808_0x8702) Decode(data []byte) (int, error) { return len(data), nil }
