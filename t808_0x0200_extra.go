package jtt

import "fmt"

// T808_0x0200_Extra_ID 附加消息 ID
//
// 0x01 - T808_0x0200_Extra_ID_Mileage 里程
// 0x02 - T808_0x0200_Extra_ID_Fuel 油量
// 0x03 - T808_0x0200_Extra_ID_Speed 速度
// 0x04 - T808_0x0200_Extra_ID_AlarmConfirm 报警确认
// 0x05~0x10 - 保留
// 0x11 - T808_0x0200_Extra_ID_SpeedLimit 超速报警
// 0x12 - T808_0x0200_Extra_ID_Region 进出区域报警
// 0x13 - T808_0x0200_Extra_ID_Route 路段行驶时间报警
// 0x14~0x24 - 保留
// 0x25 - T808_0x0200_Extra_ID_ExtSignal 扩展车辆信号状态位
// 0x2A - T808_0x0200_Extra_ID_IO IO状态位
// 0x2B - T808_0x0200_Extra_ID_Analog 模拟量
// 0x30 - T808_0x0200_Extra_ID_Signal 无线通信网络信号强度
// 0x31 - T808_0x0200_Extra_ID_Satellite GNSS定位卫星数
// 0xE1~0xFF - 自定义
type T808_0x0200_Extra_ID byte

func (id T808_0x0200_Extra_ID) String() string {
	return fmt.Sprintf("0x%02X", uint8(id))
}

const (
	// T808_0x0200_Extra_ID_Mileage 里程
	T808_0x0200_Extra_ID_Mileage T808_0x0200_Extra_ID = 0x01
	// T808_0x0200_Extra_ID_Fuel 油量
	T808_0x0200_Extra_ID_Fuel T808_0x0200_Extra_ID = 0x02
	// T808_0x0200_Extra_ID_Speed 速度
	T808_0x0200_Extra_ID_Speed T808_0x0200_Extra_ID = 0x03
	// T808_0x0200_Extra_ID_AlarmConfirm 报警确认
	T808_0x0200_Extra_ID_AlarmConfirm T808_0x0200_Extra_ID = 0x04
	// T808_0x0200_Extra_ID_SpeedLimit 超速报警
	T808_0x0200_Extra_ID_SpeedLimit T808_0x0200_Extra_ID = 0x11
	// T808_0x0200_Extra_ID_Region 进出区域报警
	T808_0x0200_Extra_ID_Region T808_0x0200_Extra_ID = 0x12
	// T808_0x0200_Extra_ID_Route 路段行驶时间报警
	T808_0x0200_Extra_ID_Route T808_0x0200_Extra_ID = 0x13
	// T808_0x0200_Extra_ID_ExtSignal 扩展车辆信号状态位
	T808_0x0200_Extra_ID_ExtSignal T808_0x0200_Extra_ID = 0x25
	// T808_0x0200_Extra_ID_IO IO状态位
	T808_0x0200_Extra_ID_IO T808_0x0200_Extra_ID = 0x2A
	// T808_0x0200_Extra_ID_Analog 模拟量
	T808_0x0200_Extra_ID_Analog T808_0x0200_Extra_ID = 0x2B
	// T808_0x0200_Extra_ID_Signal 无线通信网络信号强度
	T808_0x0200_Extra_ID_Signal T808_0x0200_Extra_ID = 0x30
	// T808_0x0200_Extra_ID_Satellite GNSS定位卫星数
	T808_0x0200_Extra_ID_Satellite T808_0x0200_Extra_ID = 0x31
)

// T808_0x0200_Extra 附加信息
type T808_0x0200_Extra struct {
	Id   T808_0x0200_Extra_ID
	Data []byte
}

// GetMileage 里程
func (e *T808_0x0200_Extra) GetMileage() (uint32, error) {
	if e.Id != T808_0x0200_Extra_ID_Mileage {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetMileage", e.Id.String(), e.Id)
	}

	render := NewReader(e.Data)
	mileage, err := render.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read mileage: %w", err)
	}
	return mileage, nil
}

// GetFuel 油量（WORD，1/10L，对应车上油量读取数）
func (e *T808_0x0200_Extra) GetFuel() (uint16, error) {
	if e.Id != T808_0x0200_Extra_ID_Fuel {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetFuel", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read fuel: %w", err)
	}
	return v, nil
}

// GetSpeed 行驶记录功能获取的速度（WORD，1/10km/h）
func (e *T808_0x0200_Extra) GetSpeed() (uint16, error) {
	if e.Id != T808_0x0200_Extra_ID_Speed {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetSpeed", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read speed: %w", err)
	}
	return v, nil
}

// GetAlarmConfirmId 报警确认ID（WORD，从1开始计数）
func (e *T808_0x0200_Extra) GetAlarmConfirmId() (uint16, error) {
	if e.Id != T808_0x0200_Extra_ID_AlarmConfirm {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetAlarmConfirmId", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read alarm confirm id: %w", err)
	}
	return v, nil
}

// OverspeedLocationType 位置类型
// 0: 无特定位置；1: 圆形区域；2: 矩形区域；3: 多边形区域；4: 路段
type OverspeedLocationType byte

const (
	OverspeedLocNone      OverspeedLocationType = 0
	OverspeedLocCircle    OverspeedLocationType = 1
	OverspeedLocRectangle OverspeedLocationType = 2
	OverspeedLocPolygon   OverspeedLocationType = 3
	OverspeedLocRoute     OverspeedLocationType = 4
)

// T808_0x0200_Extra_Overspeed 超速报警附加信息
// 长度 1 或 5 字节：
//
//	byte0: 位置类型；当位置类型!=0 时，后续4字节为区域/路段ID（DWORD）
type T808_0x0200_Extra_Overspeed struct {
	LocationType  OverspeedLocationType
	HasId         bool
	RegionRouteId uint32
}

// GetOverspeedInfo 解析0x11 超速报警附加信息
func (e *T808_0x0200_Extra) GetOverspeedInfo() (*T808_0x0200_Extra_Overspeed, error) {
	if e.Id != T808_0x0200_Extra_ID_SpeedLimit {
		return nil, fmt.Errorf("invalid extra id(%s/%d) for GetOverspeedInfo", e.Id.String(), e.Id)
	}
	if len(e.Data) < 1 {
		return nil, fmt.Errorf("overspeed extra too short: %d", len(e.Data))
	}
	r := NewReader(e.Data)
	lt, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read location type: %w", err)
	}
	info := &T808_0x0200_Extra_Overspeed{LocationType: OverspeedLocationType(lt)}
	if info.LocationType != OverspeedLocNone {
		id, err := r.ReadUint32()
		if err != nil {
			return nil, fmt.Errorf("read region/route id: %w", err)
		}
		info.HasId = true
		info.RegionRouteId = id
	}
	return info, nil
}

// RegionLocationType 表29 位置类型（1: 圆形；2: 矩形；3: 多边形；4: 路线）
type RegionLocationType byte

const (
	RegionLocCircle    RegionLocationType = 1
	RegionLocRectangle RegionLocationType = 2
	RegionLocPolygon   RegionLocationType = 3
	RegionLocRoute     RegionLocationType = 4
)

// T808_0x0200_Extra_Region 进出区域/路线报警附加信息，长度6
// byte0: 位置类型；byte1-4: 区域或线路ID（DWORD）；byte5: 方向（0进，1出）
type T808_0x0200_Extra_Region struct {
	LocationType RegionLocationType
	Id           uint32
	DirectionOut bool // false: 进; true: 出
}

// GetRegionAlarmInfo 解析0x12 进出区域/路线报警附加信息
func (e *T808_0x0200_Extra) GetRegionAlarmInfo() (*T808_0x0200_Extra_Region, error) {
	if e.Id != T808_0x0200_Extra_ID_Region {
		return nil, fmt.Errorf("invalid extra id(%s/%d) for GetRegionAlarmInfo", e.Id.String(), e.Id)
	}
	if len(e.Data) != 6 {
		return nil, fmt.Errorf("region extra invalid length: %d", len(e.Data))
	}
	r := NewReader(e.Data)
	lt, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read location type: %w", err)
	}
	id, err := r.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("read region/route id: %w", err)
	}
	dir, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read direction: %w", err)
	}
	return &T808_0x0200_Extra_Region{
		LocationType: RegionLocationType(lt),
		Id:           id,
		DirectionOut: dir == 1,
	}, nil
}

// T808_0x0200_Extra_RouteTime 路段行驶时间不足/过长报警附加信息，长度7
// byte0-3: 路段ID（DWORD）；byte4-5: 路段行驶时间（WORD, s）；byte6: 结果（0不足；1 过长）
type T808_0x0200_Extra_RouteTime struct {
	RouteId  uint32
	TimeSec  uint16
	Overlong bool // false: 不足; true: 过长
}

// GetRouteTimeAlarmInfo 解析0x13 路段行驶时间报警附加信息
func (e *T808_0x0200_Extra) GetRouteTimeAlarmInfo() (*T808_0x0200_Extra_RouteTime, error) {
	if e.Id != T808_0x0200_Extra_ID_Route {
		return nil, fmt.Errorf("invalid extra id(%s/%d) for GetRouteTimeAlarmInfo", e.Id.String(), e.Id)
	}
	if len(e.Data) != 7 {
		return nil, fmt.Errorf("route time extra invalid length: %d", len(e.Data))
	}
	r := NewReader(e.Data)
	rid, err := r.ReadUint32()
	if err != nil {
		return nil, fmt.Errorf("read route id: %w", err)
	}
	ts, err := r.ReadUint16()
	if err != nil {
		return nil, fmt.Errorf("read route time: %w", err)
	}
	res, err := r.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read result: %w", err)
	}
	return &T808_0x0200_Extra_RouteTime{RouteId: rid, TimeSec: ts, Overlong: res == 1}, nil
}

// T808_0x0200_Extra_ExtSignalBits 扩展车辆信号状态位（WORD）
// 位定义：
//
//	bit0  近光灯信号
//	bit1  远光灯信号
//	bit2  右转向信号
//	bit3  左转向信号
//	bit4  制动信号
//	bit5  反转信号
//	bit6  雾灯信号
//	bit7  示廓灯
//	bit8  喇叭信号
//	bit9  空调状态
//	bit10 门磁状态
//	bit11 缓速器工作
//	bit12 ABS 工作
//	bit13 加热器工作
//	bit14 离合器状态
//	bit15 保留
type T808_0x0200_Extra_ExtSignalBits uint16

func (b T808_0x0200_Extra_ExtSignalBits) LowBeam() bool             { return GetBitUint16(uint16(b), 0) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetLowBeam(v bool)        { SetBitUint16((*uint16)(b), 0, v) }
func (b T808_0x0200_Extra_ExtSignalBits) HighBeam() bool            { return GetBitUint16(uint16(b), 1) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetHighBeam(v bool)       { SetBitUint16((*uint16)(b), 1, v) }
func (b T808_0x0200_Extra_ExtSignalBits) RightTurn() bool           { return GetBitUint16(uint16(b), 2) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetRightTurn(v bool)      { SetBitUint16((*uint16)(b), 2, v) }
func (b T808_0x0200_Extra_ExtSignalBits) LeftTurn() bool            { return GetBitUint16(uint16(b), 3) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetLeftTurn(v bool)       { SetBitUint16((*uint16)(b), 3, v) }
func (b T808_0x0200_Extra_ExtSignalBits) Brake() bool               { return GetBitUint16(uint16(b), 4) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetBrake(v bool)          { SetBitUint16((*uint16)(b), 4, v) }
func (b T808_0x0200_Extra_ExtSignalBits) Reverse() bool             { return GetBitUint16(uint16(b), 5) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetReverse(v bool)        { SetBitUint16((*uint16)(b), 5, v) }
func (b T808_0x0200_Extra_ExtSignalBits) FogLight() bool            { return GetBitUint16(uint16(b), 6) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetFogLight(v bool)       { SetBitUint16((*uint16)(b), 6, v) }
func (b T808_0x0200_Extra_ExtSignalBits) PositionLamp() bool        { return GetBitUint16(uint16(b), 7) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetPositionLamp(v bool)   { SetBitUint16((*uint16)(b), 7, v) }
func (b T808_0x0200_Extra_ExtSignalBits) Horn() bool                { return GetBitUint16(uint16(b), 8) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetHorn(v bool)           { SetBitUint16((*uint16)(b), 8, v) }
func (b T808_0x0200_Extra_ExtSignalBits) AirConditioner() bool      { return GetBitUint16(uint16(b), 9) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetAirConditioner(v bool) { SetBitUint16((*uint16)(b), 9, v) }
func (b T808_0x0200_Extra_ExtSignalBits) DoorMagnet() bool          { return GetBitUint16(uint16(b), 10) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetDoorMagnet(v bool)     { SetBitUint16((*uint16)(b), 10, v) }
func (b T808_0x0200_Extra_ExtSignalBits) Retarder() bool            { return GetBitUint16(uint16(b), 11) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetRetarder(v bool)       { SetBitUint16((*uint16)(b), 11, v) }
func (b T808_0x0200_Extra_ExtSignalBits) ABS() bool                 { return GetBitUint16(uint16(b), 12) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetABS(v bool)            { SetBitUint16((*uint16)(b), 12, v) }
func (b T808_0x0200_Extra_ExtSignalBits) Heater() bool              { return GetBitUint16(uint16(b), 13) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetHeater(v bool)         { SetBitUint16((*uint16)(b), 13, v) }
func (b T808_0x0200_Extra_ExtSignalBits) Clutch() bool              { return GetBitUint16(uint16(b), 14) }
func (b *T808_0x0200_Extra_ExtSignalBits) SetClutch(v bool)         { SetBitUint16((*uint16)(b), 14, v) }

// GetExtSignalBits 解析0x25 扩展车辆信号状态位（WORD）
func (e *T808_0x0200_Extra) GetExtSignalBits() (T808_0x0200_Extra_ExtSignalBits, error) {
	if e.Id != T808_0x0200_Extra_ID_ExtSignal {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetExtSignalBits", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read ext signal bits: %w", err)
	}
	return T808_0x0200_Extra_ExtSignalBits(v), nil
}

// T808_0x0200_Extra_IOBits IO状态位（WORD）
// bit0: 深度休眠；bit1: 休眠；2-15: 保留
type T808_0x0200_Extra_IOBits uint16

func (b T808_0x0200_Extra_IOBits) DeepSleep() bool { return GetBitUint16(uint16(b), 0) }
func (b T808_0x0200_Extra_IOBits) Sleep() bool     { return GetBitUint16(uint16(b), 1) }

// GetIOStatus 解析0x2A IO状态位（WORD）
func (e *T808_0x0200_Extra) GetIOStatus() (T808_0x0200_Extra_IOBits, error) {
	if e.Id != T808_0x0200_Extra_ID_IO {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetIOStatus", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read io bits: %w", err)
	}
	return T808_0x0200_Extra_IOBits(v), nil
}

// GetAnalog 模拟量（DWORD）：bit0-15 AD0，bit16-31 AD1
func (e *T808_0x0200_Extra) GetAnalog() (ad0 uint16, ad1 uint16, err error) {
	if e.Id != T808_0x0200_Extra_ID_Analog {
		return 0, 0, fmt.Errorf("invalid extra id(%s/%d) for GetAnalog", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadUint32()
	if err != nil {
		return 0, 0, fmt.Errorf("read analog: %w", err)
	}
	ad0 = uint16(v & 0xFFFF)
	ad1 = uint16((v >> 16) & 0xFFFF)
	return ad0, ad1, nil
}

// GetSignalStrength 无线通信网络信号强度（BYTE）
func (e *T808_0x0200_Extra) GetSignalStrength() (byte, error) {
	if e.Id != T808_0x0200_Extra_ID_Signal {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetSignalStrength", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read signal strength: %w", err)
	}
	return v, nil
}

// GetSatelliteCount GNSS定位卫星数（BYTE）
func (e *T808_0x0200_Extra) GetSatelliteCount() (byte, error) {
	if e.Id != T808_0x0200_Extra_ID_Satellite {
		return 0, fmt.Errorf("invalid extra id(%s/%d) for GetSatelliteCount", e.Id.String(), e.Id)
	}
	r := NewReader(e.Data)
	v, err := r.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read satellite count: %w", err)
	}
	return v, nil
}

// -------------------- Setters --------------------

// SetMileage 设置 0x01 里程（DWORD）
func (e *T808_0x0200_Extra) SetMileage(mileage uint32) {
	e.Id = T808_0x0200_Extra_ID_Mileage
	w := NewWriter()
	w.WriteUint32(mileage)
	e.Data = w.Bytes()
}

// SetFuel 设置 0x02 油量（WORD，1/10L）
func (e *T808_0x0200_Extra) SetFuel(fuel uint16) {
	e.Id = T808_0x0200_Extra_ID_Fuel
	w := NewWriter()
	w.WriteUint16(fuel)
	e.Data = w.Bytes()
}

// SetSpeed 设置 0x03 行驶记录速度（WORD，1/10km/h）
func (e *T808_0x0200_Extra) SetSpeed(speed uint16) {
	e.Id = T808_0x0200_Extra_ID_Speed
	w := NewWriter()
	w.WriteUint16(speed)
	e.Data = w.Bytes()
}

// SetAlarmConfirmId 设置 0x04 报警确认ID（WORD）
func (e *T808_0x0200_Extra) SetAlarmConfirmId(id uint16) {
	e.Id = T808_0x0200_Extra_ID_AlarmConfirm
	w := NewWriter()
	w.WriteUint16(id)
	e.Data = w.Bytes()
}

// SetOverspeedInfo 设置 0x11 超速报警附加信息
// 当 LocationType != 0 时，序列化后续 4 字节区域/路段ID
func (e *T808_0x0200_Extra) SetOverspeedInfo(info T808_0x0200_Extra_Overspeed) {
	e.Id = T808_0x0200_Extra_ID_SpeedLimit
	w := NewWriter()
	w.WriteByte(byte(info.LocationType))
	if info.LocationType != OverspeedLocNone {
		w.WriteUint32(info.RegionRouteId)
	}
	e.Data = w.Bytes()
}

// SetRegionAlarmInfo 设置 0x12 进出区域/路线报警附加信息，长度6
func (e *T808_0x0200_Extra) SetRegionAlarmInfo(info T808_0x0200_Extra_Region) {
	e.Id = T808_0x0200_Extra_ID_Region
	w := NewWriter()
	w.WriteByte(byte(info.LocationType))
	w.WriteUint32(info.Id)
	if info.DirectionOut {
		w.WriteByte(1)
	} else {
		w.WriteByte(0)
	}
	e.Data = w.Bytes()
}

// SetRouteTimeAlarmInfo 设置 0x13 路段行驶时间报警附加信息，长度7
func (e *T808_0x0200_Extra) SetRouteTimeAlarmInfo(info T808_0x0200_Extra_RouteTime) {
	e.Id = T808_0x0200_Extra_ID_Route
	w := NewWriter()
	w.WriteUint32(info.RouteId)
	w.WriteUint16(info.TimeSec)
	if info.Overlong {
		w.WriteByte(1)
	} else {
		w.WriteByte(0)
	}
	e.Data = w.Bytes()
}

// SetExtSignalBits 设置 0x25 扩展车辆信号状态位（WORD）
func (e *T808_0x0200_Extra) SetExtSignalBits(bits T808_0x0200_Extra_ExtSignalBits) {
	e.Id = T808_0x0200_Extra_ID_ExtSignal
	w := NewWriter()
	w.WriteUint16(uint16(bits))
	e.Data = w.Bytes()
}

// SetIOStatus 设置 0x2A IO状态位（WORD）
func (e *T808_0x0200_Extra) SetIOStatus(bits T808_0x0200_Extra_IOBits) {
	e.Id = T808_0x0200_Extra_ID_IO
	w := NewWriter()
	w.WriteUint16(uint16(bits))
	e.Data = w.Bytes()
}

// SetAnalog 设置 0x2B 模拟量（DWORD）：bit0-15 AD0，bit16-31 AD1
func (e *T808_0x0200_Extra) SetAnalog(ad0 uint16, ad1 uint16) {
	e.Id = T808_0x0200_Extra_ID_Analog
	w := NewWriter()
	v := uint32(ad0) | (uint32(ad1) << 16)
	w.WriteUint32(v)
	e.Data = w.Bytes()
}

// SetSignalStrength 设置 0x30 无线通信网络信号强度（BYTE）
func (e *T808_0x0200_Extra) SetSignalStrength(s byte) {
	e.Id = T808_0x0200_Extra_ID_Signal
	w := NewWriter()
	w.WriteByte(s)
	e.Data = w.Bytes()
}

// SetSatelliteCount 设置 0x31 GNSS定位卫星数（BYTE）
func (e *T808_0x0200_Extra) SetSatelliteCount(n byte) {
	e.Id = T808_0x0200_Extra_ID_Satellite
	w := NewWriter()
	w.WriteByte(n)
	e.Data = w.Bytes()
}
