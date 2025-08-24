package jtt

import (
	"fmt"
	"math"
)

// T808_0x0107 查询终端属性应答
type T808_0x0107 struct {
	// 终端类型 WORD
	//
	//	位定义
	//
	//	bit0：客运车辆（0:禁用；1:启用）
	//	bit1：危化品车辆（0:禁用；1:启用）
	//	bit2：普通货运车辆（0:禁用；1:启用）
	//	bit3：租赁车辆（0:禁用；1:启用）
	//	bit6：支持硬盘录像（0:禁用；1:启用）
	//	bit7：0:一体机；1:分体机
	TerminalTypes TerminalTypes
	// 制造商 ID 固定 5 字节
	ManufacturerID string
	// 终端型号 固定长度，2013<=20 右补0x00；2019<=30 左补0x00
	TerminalModel string
	// 终端 ID 固定长度，2013<=7 右补0x00；2019=30
	TerminalID string
	// 终端 SIM 卡 ICCID，BCD[10]
	ICCID string
	// 终端硬件版本号
	HWVersion string
	// 终端固件版本号
	FWVersion string
	// GNSS 模块属性 BYTE
	//
	// 位定义：
	//
	//	bit0：GPS（0:禁用；1:启用）
	//	bit1：北斗（0:禁用；1:启用）
	//	bit2：GLONASS（0:禁用；1:启用）
	//	bit3：Galileo（0:禁用；1:启用）
	GNSSAttrs GNSSAttrs
	// 通信模块属性 BYTE
	//
	// 位定义：
	//
	//	bito,0:不支持GPRS通信,1:支持GPRS通信;
	//	bitl,0:不支持CDMA通信,1:支持CDMA通信;
	//	bit2,0:不支持TD-SCDMA通信,1:支持TD-SCDMA通信;
	//	bit3,0:不支持WCDMA通信,1:支持WCDMA通信;
	//	bit4,0:不支持CDMA2000通信,1:支持CDMA2000通信。
	//	bit5,0:不支持TD-LTE通信,1:支持TD-LTE通信;
	//	bit7,0:不支持其他通信方式,1:支持其他通信方式。
	CommAttrs CommAttrs

	// 协议版本号：-1=2011, 0=2013(默认), 1=2019
	protocolVersion VersionType
}

func (m *T808_0x0107) MsgID() MsgID { return MsgT808_0x0107 }

// ProtocolVersion 返回当前消息体所采用的协议版本。
func (m *T808_0x0107) ProtocolVersion() VersionType { return m.protocolVersion }

// SetProtocolVersion 设置该消息体的协议版本。
func (m *T808_0x0107) SetProtocolVersion(v VersionType) { m.protocolVersion = v }

func (m *T808_0x0107) Encode() ([]byte, error) {
	ver := m.protocolVersion
	if ver == 0 {
		ver = Version2013
	}

	w := NewWriter()
	w.WriteWord(uint16(m.TerminalTypes))

	// GB18030 字节长度
	gbLen := func(s string) (int, error) { return GB18030Length(s) }

	// ManufacturerID BYTE[5]
	if n, err := gbLen(m.ManufacturerID); err != nil {
		return nil, fmt.Errorf("get GB18030 length for ManufacturerID: %w", err)
	} else if n != 5 {
		return nil, fmt.Errorf("invalid ManufacturerID: %w (need 5 bytes, got %d)", ErrInvalidBody, n)
	}
	if err := w.WriteString(m.ManufacturerID, 5); err != nil {
		return nil, fmt.Errorf("write ManufacturerID: %w", err)
	}

	switch ver {
	case Version2013:
		// TerminalModel <=20，右补 0x00
		if n, err := gbLen(m.TerminalModel); err != nil {
			return nil, fmt.Errorf("get GB18030 length for TerminalModel: %w", err)
		} else if n > 20 {
			return nil, fmt.Errorf("invalid TerminalModel: %w (need 20 bytes, got %d)", ErrInvalidBody, n)
		}
		if err := w.WriteString(m.TerminalModel, 20); err != nil {
			return nil, fmt.Errorf("write TerminalModel: %w", err)
		}
		// TerminalID <=7，右补 0x00
		if n, err := gbLen(m.TerminalID); err != nil {
			return nil, fmt.Errorf("get GB18030 length for TerminalID: %w", err)
		} else if n > 7 {
			return nil, fmt.Errorf("invalid TerminalID: %w (need 7 bytes, got %d)", ErrInvalidBody, n)
		}
		if err := w.WriteString(m.TerminalID, 7); err != nil {
			return nil, fmt.Errorf("write TerminalID: %w", err)
		}

	case Version2019:
		// TerminalModel <=30，右补 0x00
		if n, err := gbLen(m.TerminalModel); err != nil {
			return nil, fmt.Errorf("get GB18030 length for TerminalModel: %w", err)
		} else if n > 30 {
			return nil, fmt.Errorf("invalid TerminalModel: %w (need 30 bytes, got %d)", ErrInvalidBody, n)
		}
		if err := w.WriteString(m.TerminalModel, 30); err != nil {
			return nil, fmt.Errorf("write TerminalModel: %w", err)
		}
		// TerminalID BYTE[30]
		if n, err := gbLen(m.TerminalID); err != nil {
			return nil, fmt.Errorf("get GB18030 length for TerminalID: %w", err)
		} else if n > 30 {
			return nil, fmt.Errorf("invalid TerminalID: %w (need 30 bytes, got %d)", ErrInvalidBody, n)
		}
		if err := w.WriteString(m.TerminalID, 30); err != nil {
			return nil, fmt.Errorf("write TerminalID: %w", err)
		}
	default:
		return nil, fmt.Errorf("encode 0x0107: unsupported protocolVersion: %d", ver)
	}
	// ICCID BCD[10]
	w.WriteBcd(m.ICCID, 10)
	// 硬件版本
	if n, err := GB18030Length(m.HWVersion); err != nil {
		return nil, fmt.Errorf("get GB18030 length for HWVersion: %w", err)
	} else {
		if n > math.MaxUint8 {
			return nil, fmt.Errorf("invalid HWVersion: %w (need %d bytes, got %d)", ErrInvalidBody, math.MaxUint8, n)
		}
		w.WriteByte(byte(n))
		if err := w.WriteString(m.HWVersion, n); err != nil {
			return nil, fmt.Errorf("write HWVersion: %w", err)
		}
	}
	// 固件版本
	if n, err := GB18030Length(m.FWVersion); err != nil {
		return nil, fmt.Errorf("get GB18030 length for FWVersion: %w", err)
	} else {
		if n > math.MaxUint8 {
			return nil, fmt.Errorf("invalid FWVersion: %w (need %d bytes, got %d)", ErrInvalidBody, math.MaxUint8, n)
		}
		w.WriteByte(byte(n))
		if err := w.WriteString(m.FWVersion, n); err != nil {
			return nil, fmt.Errorf("write FWVersion: %w", err)
		}
	}

	// GNSS 模块属性
	w.WriteByte(byte(m.GNSSAttrs))
	// 通信模块属性
	w.WriteByte(byte(m.CommAttrs))
	return w.Bytes(), nil
}

func (m *T808_0x0107) Decode(data []byte) (int, error) {
	ver := m.protocolVersion
	if ver == 0 {
		ver = Version2013
	}

	r := NewReader(data)

	terminalType, err := r.ReadWord()
	if err != nil {
		return 0, fmt.Errorf("read terminalType: %w", err)
	}
	m.TerminalTypes = TerminalTypes(terminalType)

	var b []byte
	switch ver {
	case Version2013:
		// 制造商 ID
		if b, err = r.Read(5); err != nil {
			return 0, fmt.Errorf("read ManufacturerID: %w", err)
		}
		m.ManufacturerID = BytesToString(b)
		// 终端型号
		if b, err = r.Read(20); err != nil {
			return 0, fmt.Errorf("read TerminalModel: %w", err)
		}
		m.TerminalModel = BytesToString(b)
		// 终端 ID
		if b, err = r.Read(7); err != nil {
			return 0, fmt.Errorf("read TerminalID: %w", err)
		}
		m.TerminalID = BytesToString(b)
	case Version2019:
		// 制造商 ID
		if b, err = r.Read(5); err != nil {
			return 0, fmt.Errorf("read ManufacturerID: %w", err)
		}
		m.ManufacturerID = BytesToString(b)
		// 终端型号
		if b, err = r.Read(30); err != nil {
			return 0, fmt.Errorf("read TerminalModel: %w", err)
		}
		m.TerminalModel = string(b)
		// 终端 ID
		if b, err = r.Read(30); err != nil {
			return 0, fmt.Errorf("read TerminalID: %w", err)
		}
		m.TerminalID = BytesToString(b)
	default:
		return 0, fmt.Errorf("decode 0x0107: unsupported protocolVersion: %d", ver)
	}
	if m.ICCID, err = r.ReadBcd(10); err != nil {
		return 0, fmt.Errorf("read ICCID: %w", err)
	}
	// 硬件版本
	var n byte
	if n, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read n: %w", err)
	}
	if n > 0 {
		if m.HWVersion, err = r.ReadString(int(n)); err != nil {
			return 0, fmt.Errorf("read HWVersion: %w", err)
		}
	}
	// 固件版本
	if n, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read n: %w", err)
	}
	if n > 0 {
		if m.FWVersion, err = r.ReadString(int(n)); err != nil {
			return 0, fmt.Errorf("read FWVersion: %w", err)
		}
	}
	// GNSS 模块属性
	var gnssByte byte
	if gnssByte, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read gnssByte: %w", err)
	}
	m.GNSSAttrs = GNSSAttrs(gnssByte)
	// 通信模块属性
	var commAttrByte byte
	if commAttrByte, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read commAttrByte: %w", err)
	}
	m.CommAttrs = CommAttrs(commAttrByte)
	return len(data) - r.Len(), nil
}

// TerminalTypes 终端类型
//
//	位定义
//
//	bit0,0:不适用客运车辆,1:适用客运车辆;
//	bitl,0:不适用危险品车辆,1:适用危险品车辆;
//	bit2,0:不适用普通货运车辆,1:适用普通货运车辆;
//	bit3,0:不适用出租车辆,1:适用出租车辆;
//	bit6,0:不支持硬盘录像,1:支持硬盘录像;
//	bit7,0:一体机,1:分体机。
type TerminalTypes uint16

// GetApplicableToPassengerVehicle 获取是否适用客运车辆
//
//	0:不适用客运车辆,1:适用客运车辆
func (t TerminalTypes) GetApplicableToPassengerVehicle() bool {
	return GetBitUint16(uint16(t), 0)
}

// SetApplicableToPassengerVehicle 设置是否适用客运车辆
//
//	0:不适用客运车辆,1:适用客运车辆
func (t *TerminalTypes) SetApplicableToPassengerVehicle(v bool) {
	SetBitUint16((*uint16)(t), 0, v)
}

// SetApplicableToDangerousGoodsVehicle 设置是否适用危险品车辆
//
//	0:不适用危险品车辆,1:适用危险品车辆
func (t *TerminalTypes) SetApplicableToDangerousGoodsVehicle(v bool) {
	SetBitUint16((*uint16)(t), 1, v)
}

// GetApplicableToDangerousGoodsVehicle 获取是否适用危险品车辆
//
//	0:不适用危险品车辆,1:适用危险品车辆
func (t TerminalTypes) GetApplicableToDangerousGoodsVehicle() bool {
	return GetBitUint16(uint16(t), 1)
}

// SetApplicableToOrdinaryCargoVehicle 设置是否适用普通货运车辆
//
//	0:不适用普通货运车辆,1:适用普通货运车辆
func (t *TerminalTypes) SetApplicableToOrdinaryCargoVehicle(v bool) {
	SetBitUint16((*uint16)(t), 2, v)
}

// GetApplicableToOrdinaryCargoVehicle 获取是否适用普通货运车辆
//
//	0:不适用普通货运车辆,1:适用普通货运车辆
func (t TerminalTypes) GetApplicableToOrdinaryCargoVehicle() bool {
	return GetBitUint16(uint16(t), 2)
}

// SetApplicableToRentalVehicle 设置是否适用出租车辆
//
//	0:不适用出租车辆,1:适用出租车辆
func (t *TerminalTypes) SetApplicableToRentalVehicle(v bool) {
	SetBitUint16((*uint16)(t), 3, v)
}

// GetApplicableToRentalVehicle 获取是否适用出租车辆
//
//	0:不适用出租车辆,1:适用出租车辆
func (t TerminalTypes) GetApplicableToRentalVehicle() bool {
	return GetBitUint16(uint16(t), 3)
}

// SetSupportHardDiskRecorder 设置是否支持硬盘录像
//
//	0:不支持硬盘录像,1:支持硬盘录像
func (t *TerminalTypes) SetSupportHardDiskRecorder(v bool) {
	SetBitUint16((*uint16)(t), 6, v)
}

// GetSupportHardDiskRecorder 获取是否支持硬盘录像
//
//	0:不支持硬盘录像,1:支持硬盘录像
func (t TerminalTypes) GetSupportHardDiskRecorder() bool {
	return GetBitUint16(uint16(t), 6)
}

// SetIntegratedTerminal 设置是否一体机
//
//	0:一体机,1:分体机
func (t *TerminalTypes) SetIntegratedTerminal(v bool) {
	SetBitUint16((*uint16)(t), 7, v)
}

// GetIntegratedTerminal 获取是否一体机
//
//	0:一体机,1:分体机
func (t TerminalTypes) GetIntegratedTerminal() bool {
	return GetBitUint16(uint16(t), 7)
}

// CommAttrs 通信模块属性
//
//	位定义
//
//	bit0,0:不支持GPRS通信,1:支持GPRS通信;
//	bitl,0:不支持CDMA通信,1:支持CDMA通信;
//	bit2,0:不支持TD-SCDMA通信,1:支持TD-SCDMA通信;
//	bit3,0:不支持WCDMA通信,1:支持WCDMA通信;
//	bit4,0:不支持CDMA2000通信,1:支持CDMA2000通信。
//	bit5,0:不支持TD-LTE通信,1:支持TD-LTE通信;
//	bit7,0:不支持其他通信方式,1:支持其他通信方式。
type CommAttrs byte

// GetGPRS 获取是否支持 GPRS
//
//	0:不支持GPRS通信,1:支持GPRS通信
func (c CommAttrs) GetGPRS() bool {
	return GetBitByte(byte(c), 0)
}

// SetGPRS 设置是否支持 GPRS
//
//	0:不支持GPRS通信,1:支持GPRS通信
func (c *CommAttrs) SetGPRS(v bool) {
	SetBitByte((*byte)(c), 0, v)
}

// GetCDMA 获取是否支持 CDMA
//
//	0:不支持CDMA通信,1:支持CDMA通信
func (c CommAttrs) GetCDMA() bool {
	return GetBitByte(byte(c), 1)
}

// SetCDMA 设置是否支持 CDMA
//
//	0:不支持CDMA通信,1:支持CDMA通信
func (c *CommAttrs) SetCDMA(v bool) {
	SetBitByte((*byte)(c), 1, v)
}

// GetTDSCDMA 获取是否支持 TD-SCDMA
//
//	0:不支持TD-SCDMA通信,1:支持TD-SCDMA通信
func (c CommAttrs) GetTDSCDMA() bool {
	return GetBitByte(byte(c), 2)
}

// SetTDSCDMA 设置是否支持 TD-SCDMA
//
//	0:不支持TD-SCDMA通信,1:支持TD-SCDMA通信
func (c *CommAttrs) SetTDSCDMA(v bool) {
	SetBitByte((*byte)(c), 2, v)
}

// GetWCDMA 获取是否支持 WCDMA
//
//	0:不支持WCDMA通信,1:支持WCDMA通信
func (c CommAttrs) GetWCDMA() bool {
	return GetBitByte(byte(c), 3)
}

// SetWCDMA 设置是否支持 WCDMA
//
//	0:不支持WCDMA通信,1:支持WCDMA通信
func (c *CommAttrs) SetWCDMA(v bool) {
	SetBitByte((*byte)(c), 3, v)
}

// GetCDMA2000 获取是否支持 CDMA2000
//
//	0:不支持CDMA2000通信,1:支持CDMA2000通信
func (c CommAttrs) GetCDMA2000() bool {
	return GetBitByte(byte(c), 4)
}

// SetCDMA2000 设置是否支持 CDMA2000
//
//	0:不支持CDMA2000通信,1:支持CDMA2000通信
func (c *CommAttrs) SetCDMA2000(v bool) {
	SetBitByte((*byte)(c), 4, v)
}

// GetTDLTE 获取是否支持 TD-LTE
//
//	0:不支持TD-LTE通信,1:支持TD-LTE通信
func (c CommAttrs) GetTDLTE() bool {
	return GetBitByte(byte(c), 5)
}

// SetTDLTE 设置是否支持 TD-LTE
//
//	0:不支持TD-LTE通信,1:支持TD-LTE通信
func (c *CommAttrs) SetTDLTE(v bool) {
	SetBitByte((*byte)(c), 5, v)
}

// GetOther 获取是否支持其他通信方式
//
//	0:不支持其他通信方式,1:支持其他通信方式
func (c CommAttrs) GetOther() bool {
	return GetBitByte(byte(c), 7)
}

// SetOther 设置是否支持其他通信方式
//
//	0:不支持其他通信方式,1:支持其他通信方式
func (c *CommAttrs) SetOther(v bool) {
	SetBitByte((*byte)(c), 7, v)
}
