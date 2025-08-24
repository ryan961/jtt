package jtt

import "fmt"

// ConfirmedAlarmTypes 人工确认报警类型
//
// 位定义
//
//	bit0  : 确认紧急报警
//	bit1~2: 保留
//	bit3  : 确认危险预警
//	bit4~19: 保留
//	bit20 : 确认进出区域报警
//	bit21 : 确认进出路线报警
//	bit22 : 确认路段行驶时间不足/过长报警
//	bit23~26: 保留
//	bit27 : 确认车辆非法点火报警
//	bit28 : 确认车辆非法位移报警
//	bit29~31: 保留
type ConfirmedAlarmTypes uint32

// GetEmergency 获取紧急报警
func (c ConfirmedAlarmTypes) GetEmergency() bool { return GetBitUint32(uint32(c), 0) }

// SetEmergency 设置紧急报警
func (c *ConfirmedAlarmTypes) SetEmergency(v bool) { SetBitUint32((*uint32)(c), 0, v) }

// GetDangerWarning 获取危险预警
func (c ConfirmedAlarmTypes) GetDangerWarning() bool { return GetBitUint32(uint32(c), 3) }

// SetDangerWarning 设置危险预警
func (c *ConfirmedAlarmTypes) SetDangerWarning(v bool) { SetBitUint32((*uint32)(c), 3, v) }

// GetInOutRegion 获取进出区域报警
func (c ConfirmedAlarmTypes) GetInOutRegion() bool { return GetBitUint32(uint32(c), 20) }

// SetInOutRegion 设置进出区域报警
func (c *ConfirmedAlarmTypes) SetInOutRegion(v bool) { SetBitUint32((*uint32)(c), 20, v) }

// GetInOutRoute 获取进出路线报警
func (c ConfirmedAlarmTypes) GetInOutRoute() bool { return GetBitUint32(uint32(c), 21) }

// SetInOutRoute 设置进出路线报警
func (c *ConfirmedAlarmTypes) SetInOutRoute(v bool) { SetBitUint32((*uint32)(c), 21, v) }

// GetRouteTimeTooShortOrTooLong 获取路段行驶时间不足/过长报警
func (c ConfirmedAlarmTypes) GetRouteTimeTooShortOrTooLong() bool {
	return GetBitUint32(uint32(c), 22)
}

// SetRouteTimeTooShortOrTooLong 设置路段行驶时间不足/过长报警
func (c *ConfirmedAlarmTypes) SetRouteTimeTooShortOrTooLong(v bool) {
	SetBitUint32((*uint32)(c), 22, v)
}

// GetIllegalIgnition 获取车辆非法点火报警
func (c ConfirmedAlarmTypes) GetIllegalIgnition() bool { return GetBitUint32(uint32(c), 27) }

// SetIllegalIgnition 设置车辆非法点火报警
func (c *ConfirmedAlarmTypes) SetIllegalIgnition(v bool) { SetBitUint32((*uint32)(c), 27, v) }

// GetIllegalDisplacement 获取车辆非法位移报警
func (c ConfirmedAlarmTypes) GetIllegalDisplacement() bool { return GetBitUint32(uint32(c), 28) }

// SetIllegalDisplacement 设置车辆非法位移报警
func (c *ConfirmedAlarmTypes) SetIllegalDisplacement(v bool) { SetBitUint32((*uint32)(c), 28, v) }

// T808_0x8203 人工确认报警消息
type T808_0x8203 struct {
	// 报警消息流水号。需人工确认的报警消息流水号；0 表示该报警类型所有消息
	AlarmMsgSerialNo uint16
	// 人工确认报警类型 DWORD
	//
	// 位定义
	//
	//	bit0  : 确认紧急报警
	//	bit1~2: 保留
	//	bit3  : 确认危险预警
	//	bit4~19: 保留
	//	bit20 : 确认进出区域报警
	//	bit21 : 确认进出路线报警
	//	bit22 : 确认路段行驶时间不足/过长报警
	//	bit23~26: 保留
	//	bit27 : 确认车辆非法点火报警
	//	bit28 : 确认车辆非法位移报警
	//	bit29~31: 保留
	ConfirmedAlarmTypes ConfirmedAlarmTypes
}

func (m *T808_0x8203) MsgID() MsgID { return MsgT808_0x8203 }

func (m *T808_0x8203) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteWord(m.AlarmMsgSerialNo)
	w.WriteDWord(uint32(m.ConfirmedAlarmTypes))
	return w.Bytes(), nil
}

func (m *T808_0x8203) Decode(data []byte) (int, error) {
	if len(data) < 6 {
		return 0, fmt.Errorf("invalid body for T808_0x8203: %w (need >=6 bytes, got %d)", ErrInvalidBody, len(data))
	}
	r := NewReader(data)
	var err error
	if m.AlarmMsgSerialNo, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read AlarmMsgSerialNo: %w", err)
	}
	u, err := r.ReadDWord()
	if err != nil {
		return 0, fmt.Errorf("read Types: %w", err)
	}
	m.ConfirmedAlarmTypes = ConfirmedAlarmTypes(u)
	return len(data) - r.Len(), nil
}
