package jtt

import (
	"fmt"
	"time"
)

// T808_0x8604 多边形区域项
// 2013版本和2019版本通用（2019版本新增夜间最高速度和区域名称）
type T808_0x8604 struct {
	AreaID        uint32              // 区域ID
	AreaAttribute AreaAttribute       // 区域属性
	StartTime     time.Time           // 起始时间，若区域属性0位为0则没有该字段
	EndTime       time.Time           // 结束时间，若区域属性0位为0则没有该字段
	MaxSpeed      uint16              // 最高速度(单位：公里/小时)，若区域属性1位为0则没有该字段
	SpeedDuration byte                // 超速持续时间(单位：秒)，若区域属性1位为0则没有该字段
	PointCount    uint16              // 顶点数
	Points        []T808_0x8604_Point // 顶点列表
	NightMaxSpeed uint16              // 夜间最高速度(单位：公里/小时)，若区域属性1位为0则没有该字段（2019版本）
	AreaName      string              // 区域名称（2019版本）

	protocolVersion VersionType // 协议版本
}

// SetProtocolVersion 设置协议版本
func (entity *T808_0x8604) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

// GetProtocolVersion 获取协议版本
func (entity *T808_0x8604) GetProtocolVersion() VersionType {
	return entity.protocolVersion
}

// MsgID 获取消息ID
func (entity *T808_0x8604) MsgID() MsgID {
	return MsgT808_0x8604
}

// T808_0x8604_Point 多边形区域顶点
type T808_0x8604_Point struct {
	Lat uint32 // 顶点纬度(单位：1/10^6 度)
	Lng uint32 // 顶点经度(单位：1/10^6 度)
}

// Encode 编码消息
func (entity *T808_0x8604) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入区域ID
	writer.WriteUint32(entity.AreaID)

	// 写入区域属性
	writer.WriteUint16(uint16(entity.AreaAttribute))

	// 根据区域属性，写入起始和结束时间
	if entity.AreaAttribute.GetTimeRange() {
		writer.WriteBcdTime(entity.StartTime)
		writer.WriteBcdTime(entity.EndTime)
	}

	// 根据区域属性，写入最高速度和超速持续时间
	if entity.AreaAttribute.GetSpeedLimit() {
		// 写入最高速度
		writer.WriteUint16(entity.MaxSpeed)
		// 写入超速持续时间
		writer.WriteByte(entity.SpeedDuration)
	}

	// 写入顶点数
	writer.WriteUint16(entity.PointCount)

	// 写入顶点列表
	for _, point := range entity.Points {
		// 写入顶点纬度
		writer.WriteUint32(point.Lat)
		// 写入顶点经度
		writer.WriteUint32(point.Lng)
	}

	// 根据协议版本，写入夜间最高速度和区域名称
	if entity.protocolVersion == Version2019 {
		// 根据区域属性，写入夜间最高速度和区域名称
		if entity.AreaAttribute.GetSpeedLimit() {
			// 写入夜间最高速度
			writer.WriteUint16(entity.NightMaxSpeed)
		}

		// 名称长度
		length, err := GB18030Length(entity.AreaName)
		if err != nil {
			return nil, fmt.Errorf("get area name length: %w", err)
		}
		writer.WriteWord(uint16(length))

		// 名称
		err = writer.WriteString(entity.AreaName)
		if err != nil {
			return nil, fmt.Errorf("write area name: %w", err)
		}
	}

	return writer.Bytes(), nil
}

// Decode 解码消息
func (entity *T808_0x8604) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取区域ID
	areaID, err := reader.ReadUint32()
	if err != nil {
		return 0, fmt.Errorf("read area id: %w", err)
	}
	entity.AreaID = areaID

	// 读取区域属性
	areaAttr, err := reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read area attribute: %w", err)
	}
	entity.AreaAttribute = AreaAttribute(areaAttr)

	// 根据区域属性，读取起始和结束时间
	if entity.AreaAttribute.GetTimeRange() {
		startTime, err := reader.ReadBcdTime()
		if err != nil {
			return 0, fmt.Errorf("read start time: %w", err)
		}
		entity.StartTime = startTime

		endTime, err := reader.ReadBcdTime()
		if err != nil {
			return 0, fmt.Errorf("read end time: %w", err)
		}
		entity.EndTime = endTime
	}

	// 根据区域属性，读取最高速度和超速持续时间
	if entity.AreaAttribute.GetSpeedLimit() {
		maxSpeed, err := reader.ReadUint16()
		if err != nil {
			return 0, fmt.Errorf("read max speed: %w", err)
		}
		entity.MaxSpeed = maxSpeed

		speedDuration, err := reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read speed duration: %w", err)
		}
		entity.SpeedDuration = speedDuration
	}

	// 读取顶点数
	pointCount, err := reader.ReadUint16()
	if err != nil {
		return 0, fmt.Errorf("read point count: %w", err)
	}
	entity.PointCount = pointCount

	// 读取顶点列表
	entity.Points = make([]T808_0x8604_Point, pointCount)
	for i := 0; i < int(pointCount); i++ {
		// 读取顶点纬度
		lat, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read lat: %w", err)
		}
		entity.Points[i].Lat = lat

		// 读取顶点经度
		lng, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read lng: %w", err)
		}
		entity.Points[i].Lng = lng
	}

	// 根据协议版本，读取夜间最高速度和区域名称
	if entity.protocolVersion == Version2019 {
		// 根据区域属性，读取夜间最高速度和区域名称
		if entity.AreaAttribute.GetSpeedLimit() {
			nightMaxSpeed, err := reader.ReadUint16()
			if err != nil {
				return 0, fmt.Errorf("read night max speed: %w", err)
			}
			entity.NightMaxSpeed = nightMaxSpeed
		}

		// 读取名称长度
		length, err := reader.ReadWord()
		if err != nil {
			return 0, fmt.Errorf("read area name length: %w", err)
		}

		// 读取名称
		areaName, err := reader.ReadString(int(length))
		if err != nil {
			return 0, fmt.Errorf("read area name: %w", err)
		}
		entity.AreaName = areaName
	}
	return len(data) - reader.Len(), nil
}
