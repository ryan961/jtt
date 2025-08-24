package jtt

import "fmt"

// T808_0x0108 终端升级结果通知
type T808_0x0108 struct {
	// 升级类型
	//	UpgradeTypeTerminal: 终端
	//	UpgradeTypeICCard: 道路运输 IC 卡读卡器
	//	UpgradeTypeBDModule: 北斗卫星定位模块
	Type UpgradeType
	// 升级结果
	//	0: 成功；1: 失败；2: 取消
	Result byte
}

func (m *T808_0x0108) MsgID() MsgID { return MsgT808_0x0108 }

func (m *T808_0x0108) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(byte(m.Type))
	w.WriteByte(m.Result)
	return w.Bytes(), nil
}

func (m *T808_0x0108) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid body for T808_0x0108: %w (need >=2 bytes, got %d)", ErrInvalidBody, len(data))
	}
	r := NewReader(data)
	b, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read Type: %w", err)
	}
	m.Type = UpgradeType(b)
	b, err = r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read Result: %w", err)
	}
	m.Result = b
	return len(data) - r.Len(), nil
}
