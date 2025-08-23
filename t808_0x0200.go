package jtt

import (
	"fmt"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

// LatitudeType 纬度类型
type LatitudeType int

const (
	_ LatitudeType = iota
	// NorthLatitudeType 北纬
	NorthLatitudeType = 0
	// SouthLatitudeType 南纬
	SouthLatitudeType = 1
)

// LongitudeType 经度类型
type LongitudeType int

const (
	_ LongitudeType = iota
	// EastLongitudeType 东经
	EastLongitudeType = 0
	// WestLongitudeType 西经
	WestLongitudeType = 1
)

// T808_0x0200 汇报位置
type T808_0x0200 struct {
	// 警告
	Alarm T808_0x0200_Alarm `json:"alarm"`
	// 状态
	Status T808_0x0200_Status `json:"status"`
	// 纬度
	Lat decimal.Decimal `json:"lat"`
	// 经度
	Lng decimal.Decimal `json:"lng"`
	// 海拔高度
	// 单位：米
	Altitude uint16 `json:"altitude"`
	// 速度
	// 单位：1/10km/h
	Speed uint16 `json:"speed"`
	// 方向
	// 0-359，正北为 0，顺时针
	Direction uint16 `json:"direction"`
	// 时间
	Time time.Time `json:"time"`
	// 附加信息
	Extras []T808_0x0200_Extra `json:"extras,omitempty"`
}

func (msg *T808_0x0200) MsgID() MsgID {
	return MsgT808_0x0200
}

func (msg *T808_0x0200) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入警告标志
	writer.WriteUint32(uint32(msg.Alarm))

	// 计算经纬度
	mul := decimal.NewFromFloat(1000000)
	lat := msg.Lat.Mul(mul).IntPart()
	if lat < 0 {
		msg.Status.SetSouthLatitude(true)
	}
	lng := msg.Lng.Mul(mul).IntPart()
	if lng < 0 {
		msg.Status.SetWestLongitude(true)
	}

	// 写入状态信息
	writer.WriteUint32(uint32(msg.Status))

	// 写入纬度信息
	writer.WriteUint32(uint32(math.Abs(float64(lat))))

	// 写入经度信息
	writer.WriteUint32(uint32(math.Abs(float64(lng))))

	// 写入海拔高度
	writer.WriteUint16(msg.Altitude)

	// 写入速度信息
	writer.WriteUint16(msg.Speed)

	// 写入方向信息
	writer.WriteUint16(msg.Direction)

	// 写入时间信息
	writer.WriteBcdTime(msg.Time)

	// 写入附加信息
	for i := 0; i < len(msg.Extras); i++ {
		writer.WriteByte(byte(msg.Extras[i].Id))
		writer.WriteByte(byte(len(msg.Extras[i].Data)))
		writer.Write(msg.Extras[i].Data)
	}
	return writer.Bytes(), nil
}

func (msg *T808_0x0200) Decode(data []byte) (int, error) {
	if len(data) < 28 {
		return 0, fmt.Errorf("data length error: %d", len(data))
	}
	reader := NewReader(data)

	// 读取警告标志
	var err error
	alarm, err := reader.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read alarm: %w", err)
	}
	msg.Alarm = T808_0x0200_Alarm(alarm)

	// 读取状态信息
	status, err := reader.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read status: %w", err)
	}
	msg.Status = T808_0x0200_Status(status)

	// 读取纬度信息
	latitude, err := reader.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read latitude: %w", err)
	}

	// 读取经度信息
	longitude, err := reader.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read longitude: %w", err)
	}

	msg.Lat, msg.Lng = GetGeoPointForWGS84(
		latitude, msg.Status.GetLatitudeType() == SouthLatitudeType,
		longitude, msg.Status.GetLongitudeType() == WestLongitudeType,
	)

	// 读取海拔高度
	msg.Altitude, err = reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read altitude: %w", err)
	}

	// 读取行驶速度
	msg.Speed, err = reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read speed: %w", err)
	}

	// 读取行驶方向
	msg.Direction, err = reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read direction: %w", err)
	}

	// 读取上报时间
	msg.Time, err = reader.ReadBcdTime()
	if err != nil {
		return 0, fmt.Errorf("read time: %w", err)
	}

	// 解码附加信息
	extras := make([]T808_0x0200_Extra, 0)
	buffer := data[len(data)-reader.Len():]
	for len(buffer) >= 2 {
		id, length := buffer[0], int(buffer[1])
		buffer = buffer[2:]
		if len(buffer) < length {
			return 0, fmt.Errorf("invalid extra length: %d, buffer length: %d", length, len(buffer))
		}

		extras = append(extras, T808_0x0200_Extra{
			Id:   T808_0x0200_Extra_ID(id),
			Data: buffer[:length],
		})
		buffer = buffer[length:]
	}
	if len(extras) > 0 {
		msg.Extras = extras
	}
	return len(data) - reader.Len(), nil
}

// T808_0x0200_Status 状态位
// 位定义：
//
//	bit00：ACC 状态（0 关；1 开）
//	bit01：定位状态（0 未定位；1 定位）
//	bit02：纬度（0 北纬；1 南纬）
//	bit03：经度（0 东经；1 西经）
//	bit04：运营/停运（0 运营状态；1 停运状态）
//	bit05：经纬度加密（0 未经保密插件加密；1 已加密）
//	bit06-bit07：保留
//	bit08-bit09：载重状态 00 空车；01 半载；10 保留；11 满载
//	bit10：车辆油路（0 正常；1 断开）
//	bit11：车辆电路（0 正常；1 断开）
//	bit12：车门加锁（0 解锁；1 加锁）
//	bit13：门1（前门）（0 关；1 开）
//	bit14：门2（中门）（0 关；1 开）
//	bit15：门3（后门）（0 关；1 开）
//	bit16：门4（驾驶席门）（0 关；1 开）
//	bit17：门5（自定义）（0 关；1 开）
//	bit18：使用 GPS 卫星定位（0 未使用；1 使用）
//	bit19：使用 北斗 卫星定位（0 未使用；1 使用）
//	bit20：使用 GLONASS 卫星定位（0 未使用；1 使用）
//	bit21：使用 Galileo 卫星定位（0 未使用；1 使用）
//	bit22-bit31：保留
type T808_0x0200_Status uint32

// GetAccState 获取Acc状态
func (status T808_0x0200_Status) GetAccState() bool {
	return GetBitUint32(uint32(status), 0)
}

// SetAccState 设置 ACC 状态
func (status *T808_0x0200_Status) SetAccState(b bool) {
	SetBitUint32((*uint32)(status), 0, b)
}

// Positioning 是否正在定位
func (status T808_0x0200_Status) Positioning() bool {
	return GetBitUint32(uint32(status), 1)
}

// SetPositioning 设置定位状态
func (status *T808_0x0200_Status) SetPositioning(b bool) {
	SetBitUint32((*uint32)(status), 1, b)
}

// SetSouthLatitude 设置南纬
func (status *T808_0x0200_Status) SetSouthLatitude(b bool) {
	SetBitUint32((*uint32)(status), 2, b)
}

// SetWestLongitude 设置西经
func (status *T808_0x0200_Status) SetWestLongitude(b bool) {
	SetBitUint32((*uint32)(status), 3, b)
}

// GetLatitudeType 获取纬度类型
func (status T808_0x0200_Status) GetLatitudeType() LatitudeType {
	if GetBitUint32(uint32(status), 2) {
		return SouthLatitudeType
	}
	return NorthLatitudeType
}

// GetLongitudeType 获取经度类型
func (status T808_0x0200_Status) GetLongitudeType() LongitudeType {
	if GetBitUint32(uint32(status), 3) {
		return WestLongitudeType
	}
	return EastLongitudeType
}

// OperatingStopped [bit4] 是否停运状态（1 停运，0 运营）
func (status T808_0x0200_Status) OperatingStopped() bool { return GetBitUint32(uint32(status), 4) }
func (status *T808_0x0200_Status) SetOperatingStopped(b bool) {
	SetBitUint32((*uint32)(status), 4, b)
}

// CoordEncrypted [bit5] 经纬度是否已加密
func (status T808_0x0200_Status) CoordEncrypted() bool { return GetBitUint32(uint32(status), 5) }
func (status *T808_0x0200_Status) SetCoordEncrypted(b bool) {
	SetBitUint32((*uint32)(status), 5, b)
}

// T808_0x0200_Status_LoadStatus 载重状态（占用 bit8-bit9）
//
//	00 T808_0x0200_Status_LoadEmpty 空车
//	01 T808_0x0200_Status_LoadHalf 半载
//	10 T808_0x0200_Status_LoadReserved 保留
//	11 T808_0x0200_Status_LoadFull 满载
type T808_0x0200_Status_LoadStatus byte

const (
	T808_0x0200_Status_LoadEmpty    T808_0x0200_Status_LoadStatus = 0 // 00 空车
	T808_0x0200_Status_LoadHalf     T808_0x0200_Status_LoadStatus = 1 // 01 半载
	T808_0x0200_Status_LoadReserved T808_0x0200_Status_LoadStatus = 2 // 10 保留
	T808_0x0200_Status_LoadFull     T808_0x0200_Status_LoadStatus = 3 // 11 满载
)

// GetLoadStatus 获取载重状态（bit8-bit9）
func (status T808_0x0200_Status) GetLoadStatus() T808_0x0200_Status_LoadStatus {
	b8 := GetBitUint32(uint32(status), 8)
	b9 := GetBitUint32(uint32(status), 9)
	v := 0
	if b8 {
		v |= 1
	}
	if b9 {
		v |= 2
	}
	return T808_0x0200_Status_LoadStatus(v)
}

// SetLoadStatus 设置载重状态（bit8-bit9）
func (status *T808_0x0200_Status) SetLoadStatus(ls T808_0x0200_Status_LoadStatus) {
	// bit8 为低位，bit9 为高位（00,01,10,11）
	SetBitUint32((*uint32)(status), 8, (ls&1) == 1)
	SetBitUint32((*uint32)(status), 9, (ls&2) == 2)
}

// OilCircuitDisconnected [bit10] 车辆油路是否断开
func (status T808_0x0200_Status) OilCircuitDisconnected() bool {
	return GetBitUint32(uint32(status), 10)
}
func (status *T808_0x0200_Status) SetOilCircuitDisconnected(b bool) {
	SetBitUint32((*uint32)(status), 10, b)
}

// ElectricCircuitDisconnected [bit11] 车辆电路是否断开
func (status T808_0x0200_Status) ElectricCircuitDisconnected() bool {
	return GetBitUint32(uint32(status), 11)
}
func (status *T808_0x0200_Status) SetElectricCircuitDisconnected(b bool) {
	SetBitUint32((*uint32)(status), 11, b)
}

// DoorLocked [bit12] 车门是否加锁
func (status T808_0x0200_Status) DoorLocked() bool { return GetBitUint32(uint32(status), 12) }
func (status *T808_0x0200_Status) SetDoorLocked(b bool) {
	SetBitUint32((*uint32)(status), 12, b)
}

// Door1Open [bit13] 前门是否打开
func (status T808_0x0200_Status) Door1Open() bool { return GetBitUint32(uint32(status), 13) }
func (status *T808_0x0200_Status) SetDoor1Open(b bool) {
	SetBitUint32((*uint32)(status), 13, b)
}

// Door2Open [bit14] 中门是否打开
func (status T808_0x0200_Status) Door2Open() bool { return GetBitUint32(uint32(status), 14) }
func (status *T808_0x0200_Status) SetDoor2Open(b bool) {
	SetBitUint32((*uint32)(status), 14, b)
}

// Door3Open [bit15] 后门是否打开
func (status T808_0x0200_Status) Door3Open() bool { return GetBitUint32(uint32(status), 15) }
func (status *T808_0x0200_Status) SetDoor3Open(b bool) {
	SetBitUint32((*uint32)(status), 15, b)
}

// Door4Open [bit16] 驾驶席门是否打开
func (status T808_0x0200_Status) Door4Open() bool { return GetBitUint32(uint32(status), 16) }
func (status *T808_0x0200_Status) SetDoor4Open(b bool) {
	SetBitUint32((*uint32)(status), 16, b)
}

// Door5Open [bit17] 自定义门是否打开
func (status T808_0x0200_Status) Door5Open() bool { return GetBitUint32(uint32(status), 17) }
func (status *T808_0x0200_Status) SetDoor5Open(b bool) {
	SetBitUint32((*uint32)(status), 17, b)
}

// UseGPS [bit18] 是否使用 GPS 卫星定位
func (status T808_0x0200_Status) UseGPS() bool { return GetBitUint32(uint32(status), 18) }
func (status *T808_0x0200_Status) SetUseGPS(b bool) {
	SetBitUint32((*uint32)(status), 18, b)
}

// UseBeiDou [bit19] 是否使用 北斗 卫星定位
func (status T808_0x0200_Status) UseBeiDou() bool { return GetBitUint32(uint32(status), 19) }
func (status *T808_0x0200_Status) SetUseBeiDou(b bool) {
	SetBitUint32((*uint32)(status), 19, b)
}

// UseGLONASS [bit20] 是否使用 GLONASS 卫星定位
func (status T808_0x0200_Status) UseGLONASS() bool { return GetBitUint32(uint32(status), 20) }
func (status *T808_0x0200_Status) SetUseGLONASS(b bool) {
	SetBitUint32((*uint32)(status), 20, b)
}

// UseGalileo [bit21] 是否使用 Galileo 卫星定位
func (status T808_0x0200_Status) UseGalileo() bool { return GetBitUint32(uint32(status), 21) }
func (status *T808_0x0200_Status) SetUseGalileo(b bool) {
	SetBitUint32((*uint32)(status), 21, b)
}

// T808_0x0200_Alarm 报警标志位
// 每一位表示一个报警/预警状态，位为1表示触发。
// 处理规则：未注明的均为“标志维持至报警条件解除”；注明“收到应答后清零”的在平台应答后清零。
// 位定义：
//
//	bit00：紧急报警（触动报警开关后触发）——收到应答后清零
//	bit01：超速报警
//	bit02：疲劳驾驶
//	bit03：危险预警——收到应答后清零
//	bit04：GNSS 模块发生故障
//	bit05：GNSS 天线未接或被剪断
//	bit06：GNSS 天线短路
//	bit07：终端主电源欠压
//	bit08：终端主电源掉电
//	bit09：终端 LCD 或显示器故障
//	bit10：TTS 模块故障
//	bit11：摄像头故障
//	bit12：道路运输证 IC 卡模块故障
//	bit13：超速预警
//	bit14：疲劳驾驶预警
//	bit15-bit17：保留
//	bit18：当天累计驾驶超时
//	bit19：超时停车
//	bit20：进出区域——收到应答后清零
//	bit21：进出路线——收到应答后清零
//	bit22：路段行驶时间不足/过长——收到应答后清零
//	bit23：路线偏离报警
//	bit24：车辆 VSS 故障
//	bit25：车辆油量异常
//	bit26：车辆被盗（通过车辆防盗器）
//	bit27：车辆非法点火——收到应答后清零
//	bit28：车辆非法位移——收到应答后清零
//	bit29：碰撞预警
//	bit30：侧翻预警
//	bit31：非法开门报警（终端未设置区域时，不判断非法开门）
type T808_0x0200_Alarm uint32

// Emergency [bit0] 紧急报警（触动报警开关后触发）。收到应答后清零。
func (alarm T808_0x0200_Alarm) Emergency() bool      { return GetBitUint32(uint32(alarm), 0) }
func (alarm *T808_0x0200_Alarm) SetEmergency(b bool) { SetBitUint32((*uint32)(alarm), 0, b) }

// Overspeed [bit1] 超速报警。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) Overspeed() bool      { return GetBitUint32(uint32(alarm), 1) }
func (alarm *T808_0x0200_Alarm) SetOverspeed(b bool) { SetBitUint32((*uint32)(alarm), 1, b) }

// Fatigue [bit2] 疲劳驾驶。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) Fatigue() bool      { return GetBitUint32(uint32(alarm), 2) }
func (alarm *T808_0x0200_Alarm) SetFatigue(b bool) { SetBitUint32((*uint32)(alarm), 2, b) }

// Danger [bit3] 危险预警。收到应答后清零。
func (alarm T808_0x0200_Alarm) Danger() bool      { return GetBitUint32(uint32(alarm), 3) }
func (alarm *T808_0x0200_Alarm) SetDanger(b bool) { SetBitUint32((*uint32)(alarm), 3, b) }

// GnssFault [bit4] GNSS 模块发生故障。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) GnssFault() bool      { return GetBitUint32(uint32(alarm), 4) }
func (alarm *T808_0x0200_Alarm) SetGnssFault(b bool) { SetBitUint32((*uint32)(alarm), 4, b) }

// GnssAntennaDisconnect [bit5] GNSS 天线未接或被剪断。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) GnssAntennaDisconnect() bool { return GetBitUint32(uint32(alarm), 5) }
func (alarm *T808_0x0200_Alarm) SetGnssAntennaDisconnect(b bool) {
	SetBitUint32((*uint32)(alarm), 5, b)
}

// GnssAntennaShort [bit6] GNSS 天线短路。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) GnssAntennaShort() bool      { return GetBitUint32(uint32(alarm), 6) }
func (alarm *T808_0x0200_Alarm) SetGnssAntennaShort(b bool) { SetBitUint32((*uint32)(alarm), 6, b) }

// MainPowerUndervoltage [bit7] 终端主电源欠压。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) MainPowerUndervoltage() bool { return GetBitUint32(uint32(alarm), 7) }
func (alarm *T808_0x0200_Alarm) SetMainPowerUndervoltage(b bool) {
	SetBitUint32((*uint32)(alarm), 7, b)
}

// MainPowerDown [bit8] 终端主电源掉电。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) MainPowerDown() bool      { return GetBitUint32(uint32(alarm), 8) }
func (alarm *T808_0x0200_Alarm) SetMainPowerDown(b bool) { SetBitUint32((*uint32)(alarm), 8, b) }

// LcdFault [bit9] 终端 LCD 或显示器故障。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) LcdFault() bool      { return GetBitUint32(uint32(alarm), 9) }
func (alarm *T808_0x0200_Alarm) SetLcdFault(b bool) { SetBitUint32((*uint32)(alarm), 9, b) }

// TtsFault [bit10] TTS 模块故障。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) TtsFault() bool      { return GetBitUint32(uint32(alarm), 10) }
func (alarm *T808_0x0200_Alarm) SetTtsFault(b bool) { SetBitUint32((*uint32)(alarm), 10, b) }

// CameraFault [bit11] 摄像头故障。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) CameraFault() bool      { return GetBitUint32(uint32(alarm), 11) }
func (alarm *T808_0x0200_Alarm) SetCameraFault(b bool) { SetBitUint32((*uint32)(alarm), 11, b) }

// ICCardModuleFault [bit12] 道路运输证 IC 卡模块故障。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) ICCardModuleFault() bool      { return GetBitUint32(uint32(alarm), 12) }
func (alarm *T808_0x0200_Alarm) SetICCardModuleFault(b bool) { SetBitUint32((*uint32)(alarm), 12, b) }

// OverspeedWarn [bit13] 超速预警。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) OverspeedWarn() bool      { return GetBitUint32(uint32(alarm), 13) }
func (alarm *T808_0x0200_Alarm) SetOverspeedWarn(b bool) { SetBitUint32((*uint32)(alarm), 13, b) }

// FatigueWarn [bit14] 疲劳驾驶预警。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) FatigueWarn() bool      { return GetBitUint32(uint32(alarm), 14) }
func (alarm *T808_0x0200_Alarm) SetFatigueWarn(b bool) { SetBitUint32((*uint32)(alarm), 14, b) }

// CumulativeDrivingTimeout [bit18] 当天累计驾驶超时。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) CumulativeDrivingTimeout() bool {
	return GetBitUint32(uint32(alarm), 18)
}
func (alarm *T808_0x0200_Alarm) SetCumulativeDrivingTimeout(b bool) {
	SetBitUint32((*uint32)(alarm), 18, b)
}

// OvertimeParking [bit19] 超时停车。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) OvertimeParking() bool      { return GetBitUint32(uint32(alarm), 19) }
func (alarm *T808_0x0200_Alarm) SetOvertimeParking(b bool) { SetBitUint32((*uint32)(alarm), 19, b) }

// InOutRegion [bit20] 进出区域。收到应答后清零。
func (alarm T808_0x0200_Alarm) InOutRegion() bool      { return GetBitUint32(uint32(alarm), 20) }
func (alarm *T808_0x0200_Alarm) SetInOutRegion(b bool) { SetBitUint32((*uint32)(alarm), 20, b) }

// InOutRoute [bit21] 进出路线。收到应答后清零。
func (alarm T808_0x0200_Alarm) InOutRoute() bool      { return GetBitUint32(uint32(alarm), 21) }
func (alarm *T808_0x0200_Alarm) SetInOutRoute(b bool) { SetBitUint32((*uint32)(alarm), 21, b) }

// RoadTimeAbnormal [bit22] 路段行驶时间不足/过长。收到应答后清零。
func (alarm T808_0x0200_Alarm) RoadTimeAbnormal() bool      { return GetBitUint32(uint32(alarm), 22) }
func (alarm *T808_0x0200_Alarm) SetRoadTimeAbnormal(b bool) { SetBitUint32((*uint32)(alarm), 22, b) }

// RouteDeviation [bit23] 路线偏离报警。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) RouteDeviation() bool      { return GetBitUint32(uint32(alarm), 23) }
func (alarm *T808_0x0200_Alarm) SetRouteDeviation(b bool) { SetBitUint32((*uint32)(alarm), 23, b) }

// VssFault [bit24] 车辆 VSS 故障。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) VssFault() bool      { return GetBitUint32(uint32(alarm), 24) }
func (alarm *T808_0x0200_Alarm) SetVssFault(b bool) { SetBitUint32((*uint32)(alarm), 24, b) }

// FuelAbnormal [bit25] 车辆油量异常。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) FuelAbnormal() bool      { return GetBitUint32(uint32(alarm), 25) }
func (alarm *T808_0x0200_Alarm) SetFuelAbnormal(b bool) { SetBitUint32((*uint32)(alarm), 25, b) }

// VehicleStolen [bit26] 车辆被盗（通过车辆防盗器）。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) VehicleStolen() bool      { return GetBitUint32(uint32(alarm), 26) }
func (alarm *T808_0x0200_Alarm) SetVehicleStolen(b bool) { SetBitUint32((*uint32)(alarm), 26, b) }

// IllegalIgnition [bit27] 车辆非法点火。收到应答后清零。
func (alarm T808_0x0200_Alarm) IllegalIgnition() bool      { return GetBitUint32(uint32(alarm), 27) }
func (alarm *T808_0x0200_Alarm) SetIllegalIgnition(b bool) { SetBitUint32((*uint32)(alarm), 27, b) }

// IllegalDisplacement [bit28] 车辆非法位移。收到应答后清零。
func (alarm T808_0x0200_Alarm) IllegalDisplacement() bool      { return GetBitUint32(uint32(alarm), 28) }
func (alarm *T808_0x0200_Alarm) SetIllegalDisplacement(b bool) { SetBitUint32((*uint32)(alarm), 28, b) }

// CollisionWarn [bit29] 碰撞预警。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) CollisionWarn() bool      { return GetBitUint32(uint32(alarm), 29) }
func (alarm *T808_0x0200_Alarm) SetCollisionWarn(b bool) { SetBitUint32((*uint32)(alarm), 29, b) }

// RolloverWarn [bit30] 侧翻预警。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) RolloverWarn() bool      { return GetBitUint32(uint32(alarm), 30) }
func (alarm *T808_0x0200_Alarm) SetRolloverWarn(b bool) { SetBitUint32((*uint32)(alarm), 30, b) }

// IllegalDoorOpen [bit31] 非法开门报警（终端未设置区域时，不判断非法开门）。维持至报警条件解除。
func (alarm T808_0x0200_Alarm) IllegalDoorOpen() bool      { return GetBitUint32(uint32(alarm), 31) }
func (alarm *T808_0x0200_Alarm) SetIllegalDoorOpen(b bool) { SetBitUint32((*uint32)(alarm), 31, b) }
