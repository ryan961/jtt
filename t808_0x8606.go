package jtt

import (
	"time"
)

// RouteAttribute 路线属性定义
// 0-15 位
// 0 位：是否启用起始时间与结束时间的判断规则,0:否;1:是
// 1 位：保留
// 2 位：进路线是否报警给驾驶员,0:否;1:是
// 3 位：进路线是否报警给平台,0:否;1:是
// 4 位：出路线是否报警给驾驶员,0:否;1:是
// 5 位：出路线是否报警给平台,0:否;1:是
// 6~15 位：保留
type RouteAttribute uint16

// GetTimeRange 是否启用起始时间与结束时间的判断规则.0:否;1:是
func (attr RouteAttribute) GetTimeRange() bool {
	return GetBitUint16(uint16(attr), 0)
}

// SetTimeRange 设置是否启用起始时间与结束时间的判断规则.0:否;1:是
func (attr *RouteAttribute) SetTimeRange(value bool) {
	SetBitUint16((*uint16)(attr), 0, value)
}

// GetEnterAlertDriver 获取进路线是否报警给驾驶员,0:否;1:是
func (attr RouteAttribute) GetEnterAlertDriver() bool {
	return GetBitUint16(uint16(attr), 2)
}

// SetEnterAlertDriver 设置进路线是否报警给驾驶员,0:否;1:是
func (attr *RouteAttribute) SetEnterAlertDriver(value bool) {
	SetBitUint16((*uint16)(attr), 2, value)
}

// GetEnterAlertPlatform 获取进路线是否报警给平台,0:否;1:是
func (attr RouteAttribute) GetEnterAlertPlatform() bool {
	return GetBitUint16(uint16(attr), 3)
}

// SetEnterAlertPlatform 设置进路线是否报警给平台,0:否;1:是
func (attr *RouteAttribute) SetEnterAlertPlatform(value bool) {
	SetBitUint16((*uint16)(attr), 3, value)
}

// GetExitAlertDriver 获取出路线是否报警给驾驶员,0:否;1:是
func (attr RouteAttribute) GetExitAlertDriver() bool {
	return GetBitUint16(uint16(attr), 4)
}

// SetExitAlertDriver 设置出路线是否报警给驾驶员,0:否;1:是
func (attr *RouteAttribute) SetExitAlertDriver(value bool) {
	SetBitUint16((*uint16)(attr), 4, value)
}

// GetExitAlertPlatform 获取出路线是否报警给平台,0:否;1:是
func (attr RouteAttribute) GetExitAlertPlatform() bool {
	return GetBitUint16(uint16(attr), 5)
}

// SetExitAlertPlatform 设置出路线是否报警给平台,0:否;1:是
func (attr *RouteAttribute) SetExitAlertPlatform(value bool) {
	SetBitUint16((*uint16)(attr), 5, value)
}

// RouteSegmentAttribute 路段属性定义
// 0-7 位
// 0 位：是否设置路段行驶时间阈值,0:否;1:是
// 1 位：是否设置路段限速,0:否;1:是
// 2 位：中心纬度 0:北纬;1:南纬
// 3 位：中心经度 0:东经;1:西经
// 4~7 位：保留
type RouteSegmentAttribute byte

// GetTravelTimeThreshold 获取是否设置路段行驶时间阈值
func (attr RouteSegmentAttribute) GetTravelTimeThreshold() bool {
	return GetBitByte(byte(attr), 0)
}

// SetTravelTimeThreshold 设置路段行驶时间阈值
func (attr *RouteSegmentAttribute) SetTravelTimeThreshold(value bool) {
	SetBitByte((*byte)(attr), 0, value)
}

// GetSpeedLimit 获取路段限速
func (attr RouteSegmentAttribute) GetSpeedLimit() bool {
	return GetBitByte(byte(attr), 1)
}

// SetSpeedLimit 设置路段限速
func (attr *RouteSegmentAttribute) SetSpeedLimit(value bool) {
	SetBitByte((*byte)(attr), 1, value)
}

// GetCenterLat 获取中心维度 0:北纬;1:南纬
func (attr RouteSegmentAttribute) GetCenterLat() int {
	if GetBitByte(byte(attr), 2) {
		return 1
	}
	return 0
}

// SetCenterLat 设置中心维度 0:北纬;1:南纬
func (attr *RouteSegmentAttribute) SetCenterLat(value int) {
	if value == 1 {
		SetBitByte((*byte)(attr), 2, true)
	} else {
		SetBitByte((*byte)(attr), 2, false)
	}
}

// GetCenterLng 获取中心经度 0:东经;1:西经
func (attr RouteSegmentAttribute) GetCenterLng() int {
	if GetBitByte(byte(attr), 3) {
		return 1
	}
	return 0
}

// SetCenterLng 设置中心经度 0:东经;1:西经
func (attr *RouteSegmentAttribute) SetCenterLng(value int) {
	if value == 1 {
		SetBitByte((*byte)(attr), 3, true)
	} else {
		SetBitByte((*byte)(attr), 3, false)
	}
}

// T808_0x8606_RoutePoint 路线拐点项
type T808_0x8606_RoutePoint struct {
	PointID                uint32                // 拐点ID
	SegmentID              uint32                // 路段ID
	PointLat               uint32                // 拐点纬度(单位：1/10^6 度)
	PointLng               uint32                // 拐点经度(单位：1/10^6 度)
	SegmentWidth           byte                  // 路段宽度(单位：米)
	SegmentAttribute       RouteSegmentAttribute // 路段属性
	TravelTimeThresholdMax uint16                // 路段行驶过长阈值(单位：秒)，若路段属性0位为0则没有该字段
	TravelTimeThresholdMin uint16                // 路段行驶过短阈值(单位：秒)，若路段属性0位为0则没有该字段
	MaxSpeed               uint16                // 最高速度(单位：公里/小时)，若路段属性1位为0则没有该字段
	SpeedDuration          byte                  // 超速持续时间(单位：秒)，若路段属性1位为0则没有该字段
	NightMaxSpeed          uint16                // 路段夜间最高速度(单位：公里/小时)，若路段属性1位为0则没有该字段（2019版本）
}

// T808_0x8606 设置路线
// 2013版本和2019版本通用（2019版本新增夜间最高速度和区域名称）
type T808_0x8606 struct {
	RouteID        uint32                   // 路线ID
	RouteAttribute RouteAttribute           // 路线属性
	StartTime      time.Time                // 起始时间，若路线属性0位为0则没有该字段
	EndTime        time.Time                // 结束时间，若路线属性0位为0则没有该字段
	PointCount     uint16                   // 拐点数量
	RoutePoints    []T808_0x8606_RoutePoint // 拐点列表
	RouteName      string                   // 路线名称（2019版本）

	protocolVersion VersionType // 协议版本
}

// SetProtocolVersion 设置协议版本
func (entity *T808_0x8606) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

// GetProtocolVersion 获取协议版本
func (entity *T808_0x8606) GetProtocolVersion() VersionType {
	return entity.protocolVersion
}

// MsgID 获取消息ID
func (entity *T808_0x8606) MsgID() MsgID {
	return MsgT808_0x8606
}

// Encode 编码消息
func (entity *T808_0x8606) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入路线ID
	writer.WriteUint32(entity.RouteID)

	// 写入路线属性
	writer.WriteUint16(uint16(entity.RouteAttribute))

	// 根据路线属性，写入起始和结束时间
	if entity.RouteAttribute.GetTimeRange() {
		writer.WriteBcdTime(entity.StartTime)
		writer.WriteBcdTime(entity.EndTime)
	}

	// 写入拐点数量
	writer.WriteUint16(entity.PointCount)

	// 写入拐点列表
	for _, point := range entity.RoutePoints {
		// 写入拐点ID
		writer.WriteUint32(point.PointID)

		// 写入路段ID
		writer.WriteUint32(point.SegmentID)

		// 写入拐点纬度
		writer.WriteUint32(point.PointLat)

		// 写入拐点经度
		writer.WriteUint32(point.PointLng)

		// 写入路段宽度
		writer.WriteByte(point.SegmentWidth)

		// 写入路段属性
		writer.WriteByte(byte(point.SegmentAttribute))

		// 根据路段属性，写入行驶时间阈值
		if point.SegmentAttribute.GetTravelTimeThreshold() {
			writer.WriteUint16(point.TravelTimeThresholdMax)
			writer.WriteUint16(point.TravelTimeThresholdMin)
		}

		// 根据路段属性，写入最高速度和超速持续时间
		if point.SegmentAttribute.GetSpeedLimit() {
			writer.WriteUint16(point.MaxSpeed)
			writer.WriteByte(point.SpeedDuration)
		}

		// 根据协议版本，写入夜间最高速度
		if entity.protocolVersion == Version2019 {
			writer.WriteUint16(point.NightMaxSpeed)
		}
	}

	// 根据协议版本，写入路线名称
	if entity.protocolVersion == Version2019 {
		length, err := GB18030Length(entity.RouteName)
		if err != nil {
			return nil, err
		}
		// 写入名称长度
		writer.WriteWord(uint16(length))

		// 写入名称
		err = writer.WriteString(entity.RouteName)
		if err != nil {
			return nil, err
		}
	}

	return writer.Bytes(), nil
}

// Decode 解码消息
func (entity *T808_0x8606) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取路线ID
	var err error
	entity.RouteID, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取路线属性
	routeAttr, err := reader.ReadUint16()
	if err != nil {
		return 0, err
	}
	entity.RouteAttribute = RouteAttribute(routeAttr)

	// 根据路线属性，读取起始和结束时间
	if entity.RouteAttribute.GetTimeRange() {
		entity.StartTime, err = reader.ReadBcdTime()
		if err != nil {
			return 0, err
		}

		entity.EndTime, err = reader.ReadBcdTime()
		if err != nil {
			return 0, err
		}
	}

	// 读取拐点数量
	entity.PointCount, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取拐点列表
	entity.RoutePoints = make([]T808_0x8606_RoutePoint, entity.PointCount)
	for i := 0; i < int(entity.PointCount); i++ {
		// 读取拐点ID
		pointID, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RoutePoints[i].PointID = pointID

		// 读取路段ID
		segmentID, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RoutePoints[i].SegmentID = segmentID

		// 读取拐点纬度
		pointLat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RoutePoints[i].PointLat = pointLat

		// 读取拐点经度
		pointLng, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RoutePoints[i].PointLng = pointLng

		// 读取路段宽度
		segmentWidth, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.RoutePoints[i].SegmentWidth = segmentWidth

		// 读取路段属性
		segmentAttr, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		entity.RoutePoints[i].SegmentAttribute = RouteSegmentAttribute(segmentAttr)

		// 根据路段属性，读取行驶时间阈值
		if entity.RoutePoints[i].SegmentAttribute.GetTravelTimeThreshold() {
			travelTimeThresholdMax, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}
			entity.RoutePoints[i].TravelTimeThresholdMax = travelTimeThresholdMax

			travelTimeThresholdMin, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}
			entity.RoutePoints[i].TravelTimeThresholdMin = travelTimeThresholdMin
		}

		// 根据路段属性，读取最高速度和超速持续时间
		if entity.RoutePoints[i].SegmentAttribute.GetSpeedLimit() {
			maxSpeed, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}
			entity.RoutePoints[i].MaxSpeed = maxSpeed

			speedDuration, err := reader.ReadByte()
			if err != nil {
				return 0, err
			}
			entity.RoutePoints[i].SpeedDuration = speedDuration
		}

		// 根据协议版本，读取夜间最高速度
		if entity.protocolVersion == Version2019 {
			nightMaxSpeed, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}
			entity.RoutePoints[i].NightMaxSpeed = nightMaxSpeed
		}
	}

	// 根据协议版本，读取路线名称
	if entity.protocolVersion == Version2019 {
		routeNameLength, err := reader.ReadWord()
		if err != nil {
			return 0, err
		}
		routeName, err := reader.ReadString(int(routeNameLength))
		if err != nil {
			return 0, err
		}
		entity.RouteName = routeName
	}

	return len(data) - reader.Len(), nil
}
