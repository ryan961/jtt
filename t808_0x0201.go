package jtt

import "fmt"

// T808_0x0201 位置信息查询应答
type T808_0x0201 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16
	// 位置信息汇报
	LocationInfo *T808_0x0200
}

func (m *T808_0x0201) MsgID() MsgID { return MsgT808_0x0201 }

func (m *T808_0x0201) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteUint16(m.ReplyMsgSerialNo)
	if m.LocationInfo != nil {
		if data, err := m.LocationInfo.Encode(); err != nil {
			return nil, fmt.Errorf("encode LocationInfo: %w", err)
		} else {
			w.Write(data)
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x0201) Decode(data []byte) (int, error) {
	r := NewReader(data)
	var err error
	// 应答流水号
	m.ReplyMsgSerialNo, err = r.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read ReplyMsgSerialNo: %w", err)
	}
	// 位置信息
	bts, err := r.Read(r.Len())
	if err != nil {
		return 0, fmt.Errorf("read LocationInfo: %w", err)
	}
	locationInfo := &T808_0x0200{}
	if _, err := locationInfo.Decode(bts); err != nil {
		return 0, fmt.Errorf("decode LocationInfo: %w", err)
	}
	m.LocationInfo = locationInfo
	return len(data) - r.Len(), nil
}
