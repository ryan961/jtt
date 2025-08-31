package jtt

import (
	"fmt"
	"time"
)

// T808_0x8600 设置圆形区域
// 2013版本和2019版本通用（2019版本新增夜间最高速度和区域名称）
type T808_0x8600 struct {
	AreaCount      byte                     // 区域总数
	AreaSettingTag byte                     // 区域设置属性：0-更新区域；1-追加区域；2-修改区域
	CircleAreas    []T808_0x8600_CircleArea // 圆形区域项

	protocolVersion VersionType // 协议版本
}

// SetProtocolVersion 设置协议版本
func (entity *T808_0x8600) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

// GetProtocolVersion 获取协议版本
func (entity *T808_0x8600) GetProtocolVersion() VersionType {
	return entity.protocolVersion
}

// MsgID 获取消息ID
func (entity *T808_0x8600) MsgID() MsgID {
	return MsgT808_0x8600
}

// AreaAttribute 区域属性定义
// 0-15 位
// 0 位：是否启用起始时间与结束时间的判断规则,0:否;1:是
// 1 位：是否启用最高速度、超速持续时间和夜间最高速度的判断规则,0:否;1:是
// 2 位：进区域是否报警给驾驶员，0：不报警；1：报警
// 3 位：进区域是否报警给平台，0：不报警；1：报警
// 4 位：出区域是否报警给驾驶员，0：不报警；1：报警
// 5 位：出区域是否报警给平台，0：不报警；1：报警
// 6 位：中心点纬度，0：北纬；1：南纬
// 7 位：中心点经度，0：东经；1：西经
// 8 位：是否允许开门，0：不允许；1：允许
// 9-13 位：保留
// 14 位：进区域是否开启通信模块，0：不开启；1：开启
// 15 位：进区域是否采集 GNSS 详细定位数据，0：不采集；1：采集
type AreaAttribute uint16

// GetTimeRange 获取区域是否包含时间范围
func (attr AreaAttribute) GetTimeRange() bool {
	return GetBitUint16(uint16(attr), 0)
}

// SetTimeRange 设置区域是否包含时间范围
func (attr *AreaAttribute) SetTimeRange(value bool) {
	SetBitUint16((*uint16)(attr), 0, value)
}

// GetSpeedLimit 获取区域是否包含速度限制
func (attr AreaAttribute) GetSpeedLimit() bool {
	return GetBitUint16(uint16(attr), 1)
}

// SetSpeedLimit 设置区域是否限速
func (attr *AreaAttribute) SetSpeedLimit(value bool) {
	SetBitUint16((*uint16)(attr), 1, value)
}

// GetEnterReportDriver 获取进区域是否报警给驾驶员
func (attr AreaAttribute) GetEnterReportDriver() bool {
	return GetBitUint16(uint16(attr), 2)
}

// SetEnterReportDriver 设置进区域是否报警给驾驶员
func (attr *AreaAttribute) SetEnterReportDriver(value bool) {
	SetBitUint16((*uint16)(attr), 2, value)
}

// GetEnterReportPlatform 获取进区域是否报警给平台
func (attr AreaAttribute) GetEnterReportPlatform() bool {
	return GetBitUint16(uint16(attr), 3)
}

// SetEnterReportPlatform 设置进区域是否报警给平台
func (attr *AreaAttribute) SetEnterReportPlatform(value bool) {
	SetBitUint16((*uint16)(attr), 3, value)
}

// GetExitReportDriver 获取出区域是否报警给驾驶员
func (attr AreaAttribute) GetExitReportDriver() bool {
	return GetBitUint16(uint16(attr), 4)
}

// SetExitReportDriver 设置出区域是否报警给驾驶员
func (attr *AreaAttribute) SetExitReportDriver(value bool) {
	SetBitUint16((*uint16)(attr), 4, value)
}

// GetExitReportPlatform 获取出区域是否报警给平台
func (attr AreaAttribute) GetExitReportPlatform() bool {
	return GetBitUint16(uint16(attr), 5)
}

// SetExitReportPlatform 设置出区域是否报警给平台
func (attr *AreaAttribute) SetExitReportPlatform(value bool) {
	SetBitUint16((*uint16)(attr), 5, value)
}

// GetCenterLat 获取中心点纬度，0：北纬；1：南纬
func (attr AreaAttribute) GetCenterLat() int {
	if GetBitUint16(uint16(attr), 6) {
		return 1
	}
	return 0
}

// SetCenterLat 设置中心点纬度，0：北纬；1：南纬
func (attr *AreaAttribute) SetCenterLat(value int) {
	if value == 1 {
		SetBitUint16((*uint16)(attr), 6, true)
	} else {
		SetBitUint16((*uint16)(attr), 6, false)
	}
}

// GetCenterLng 获取中心点经度，0：东经；1：西经
func (attr AreaAttribute) GetCenterLng() int {
	if GetBitUint16(uint16(attr), 7) {
		return 1
	}
	return 0
}

// SetCenterLng 设置中心点经度，0：东经；1：西经
func (attr *AreaAttribute) SetCenterLng(value int) {
	if value == 1 {
		SetBitUint16((*uint16)(attr), 7, true)
	} else {
		SetBitUint16((*uint16)(attr), 7, false)
	}
}

// GetOpenDoor 获取是否允许开门
func (attr AreaAttribute) GetOpenDoor() bool {
	return GetBitUint16(uint16(attr), 8)
}

// SetOpenDoor 设置是否允许开门
func (attr *AreaAttribute) SetOpenDoor(value bool) {
	SetBitUint16((*uint16)(attr), 8, value)
}

// GetEnterOpenCommModule 获取进区域是否开启通信模块
func (attr AreaAttribute) GetEnterOpenCommModule() bool {
	return GetBitUint16(uint16(attr), 14)
}

// SetEnterOpenCommModule 设置进区域是否开启通信模块
func (attr *AreaAttribute) SetEnterOpenCommModule(value bool) {
	SetBitUint16((*uint16)(attr), 14, value)
}

// GetEnterCollectGnssDetail 获取进区域是否采集 GNSS 详细定位数据
func (attr AreaAttribute) GetEnterCollectGnssDetail() bool {
	return GetBitUint16(uint16(attr), 15)
}

// SetEnterCollectGnssDetail 设置进区域是否采集 GNSS 详细定位数据
func (attr *AreaAttribute) SetEnterCollectGnssDetail(value bool) {
	SetBitUint16((*uint16)(attr), 15, value)
}

// T808_0x8600_CircleArea 圆形区域项
type T808_0x8600_CircleArea struct {
	AreaID        uint32        // 区域ID
	AreaAttribute AreaAttribute // 区域属性
	CenterLat     uint32        // 中心点纬度(单位：1/10^6 度)
	CenterLng     uint32        // 中心点经度(单位：1/10^6 度)
	Radius        uint32        // 半径(单位：米)
	StartTime     time.Time     // 起始时间，若区域属性0位为0则没有该字段
	EndTime       time.Time     // 结束时间，若区域属性0位为0则没有该字段
	MaxSpeed      uint16        // 最高速度(单位：公里/小时)，若区域属性1位为0则没有该字段
	SpeedDuration byte          // 超速持续时间(单位：秒)，若区域属性1位为0则没有该字段
	NightMaxSpeed uint16        // 夜间最高速度(单位：公里/小时)，若区域属性1位为0则没有该字段（2019版本）
	AreaName      string        // 区域名称（2019版本）
}

// Encode 编码消息
func (entity *T808_0x8600) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入区域设置属性
	writer.WriteByte(entity.AreaSettingTag)

	// 写入区域总数
	writer.WriteByte(entity.AreaCount)

	// 写入圆形区域项
	for _, area := range entity.CircleAreas {
		// 写入区域ID
		writer.WriteUint32(area.AreaID)

		// 写入区域属性
		writer.WriteUint16(uint16(area.AreaAttribute))

		// 写入中心点纬度
		writer.WriteUint32(area.CenterLat)

		// 写入中心点经度
		writer.WriteUint32(area.CenterLng)

		// 写入半径
		writer.WriteUint32(area.Radius)

		// 根据区域属性，写入起始和结束时间
		if area.AreaAttribute.GetTimeRange() {
			writer.WriteBcdTime(area.StartTime)
			writer.WriteBcdTime(area.EndTime)
		}

		// 根据区域属性，写入最高速度和超速持续时间
		if area.AreaAttribute.GetSpeedLimit() {
			// 写入最高速度
			writer.WriteUint16(area.MaxSpeed)
			// 写入超速持续时间
			writer.WriteByte(area.SpeedDuration)

			// 根据协议版本，写入夜间最高速度
			if entity.protocolVersion == Version2019 {
				writer.WriteUint16(area.NightMaxSpeed)
			}
		}

		// 根据协议版本，写入区域名称
		if entity.protocolVersion == Version2019 {
			// 写入名称长度
			length, err := GB18030Length(area.AreaName)
			if err != nil {
				return nil, fmt.Errorf("write area name length: %w", err)
			}
			writer.WriteWord(uint16(length))

			// 写入名称
			err = writer.WriteString(area.AreaName)
			if err != nil {
				return nil, fmt.Errorf("write area name: %w", err)
			}
		}
	}

	return writer.Bytes(), nil
}

// Decode 解码消息
func (entity *T808_0x8600) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	// 读取区域设置属性
	entity.AreaSettingTag, err = reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read area setting tag: %w", err)
	}

	// 读取区域总数
	entity.AreaCount, err = reader.ReadByte()
	if err != nil {
		return 0, fmt.Errorf("read area count: %w", err)
	}

	// 读取圆形区域项
	entity.CircleAreas = make([]T808_0x8600_CircleArea, entity.AreaCount)
	for i := 0; i < int(entity.AreaCount); i++ {
		// 读取区域ID
		areaID, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read area id: %w", err)
		}
		entity.CircleAreas[i].AreaID = areaID

		// 读取区域属性
		areaAttr, err := reader.ReadUint16()
		if err != nil {
			return 0, fmt.Errorf("read area attribute: %w", err)
		}
		entity.CircleAreas[i].AreaAttribute = AreaAttribute(areaAttr)

		// 读取中心点纬度
		centerLat, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read center lat: %w", err)
		}
		entity.CircleAreas[i].CenterLat = centerLat

		// 读取中心点经度
		centerLng, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read center lng: %w", err)
		}
		entity.CircleAreas[i].CenterLng = centerLng

		// 读取半径
		radius, err := reader.ReadUint32()
		if err != nil {
			return 0, fmt.Errorf("read radius: %w", err)
		}
		entity.CircleAreas[i].Radius = radius

		// 根据区域属性，读取起始和结束时间
		if entity.CircleAreas[i].AreaAttribute.GetTimeRange() {
			startTime, err := reader.ReadBcdTime()
			if err != nil {
				return 0, fmt.Errorf("read start time: %w", err)
			}
			entity.CircleAreas[i].StartTime = startTime

			endTime, err := reader.ReadBcdTime()
			if err != nil {
				return 0, fmt.Errorf("read end time: %w", err)
			}
			entity.CircleAreas[i].EndTime = endTime
		}

		// 根据区域属性，读取最高速度和超速持续时间
		if entity.CircleAreas[i].AreaAttribute.GetSpeedLimit() {
			maxSpeed, err := reader.ReadUint16()
			if err != nil {
				return 0, fmt.Errorf("read max speed: %w", err)
			}
			entity.CircleAreas[i].MaxSpeed = maxSpeed

			speedDuration, err := reader.ReadByte()
			if err != nil {
				return 0, fmt.Errorf("read speed duration: %w", err)
			}
			entity.CircleAreas[i].SpeedDuration = speedDuration

			// 根据协议版本，读取夜间最高速度
			if entity.protocolVersion == Version2019 {
				nightMaxSpeed, err := reader.ReadUint16()
				if err != nil {
					return 0, fmt.Errorf("read night max speed: %w", err)
				}
				entity.CircleAreas[i].NightMaxSpeed = nightMaxSpeed
			}
		}

		// 根据协议版本，读取区域名称
		if entity.protocolVersion == Version2019 {
			areaNameLength, err := reader.ReadWord()
			if err != nil {
				return 0, fmt.Errorf("read area name length: %w", err)
			}
			if areaNameLength > 0 {
				areaName, err := reader.ReadString(int(areaNameLength))
				if err != nil {
					return 0, fmt.Errorf("read area name: %w", err)
				}
				entity.CircleAreas[i].AreaName = areaName
			}
		}
	}

	return len(data) - reader.Len(), nil
}
