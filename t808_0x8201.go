package jtt

// T808_0x8201 位置信息查询应答
type T808_0x8201 struct{}

func (m *T808_0x8201) MsgID() MsgID { return MsgT808_0x8201 }

func (m *T808_0x8201) Encode() ([]byte, error) { return nil, nil }

func (m *T808_0x8201) Decode(data []byte) (int, error) { return 0, nil }
