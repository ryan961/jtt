package jtt

// T808_0x8204 链路检测
type T808_0x8204 struct{}

func (m *T808_0x8204) MsgID() MsgID { return MsgT808_0x8204 }

func (m *T808_0x8204) Encode() ([]byte, error) { return nil, nil }

func (m *T808_0x8204) Decode(data []byte) (int, error) { return 0, nil }
