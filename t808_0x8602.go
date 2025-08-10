package jtt

import (
	"time"
)

// T808_0x8602_RectangleArea 矩形区域项
type T808_0x8602_RectangleArea struct {
	AreaID         uint32        // 区域ID
	AreaAttribute  AreaAttribute // 区域属性
	LeftTopLat     uint32        // 左上点纬度(单位：1/10^6 度)
	LeftTopLng     uint32        // 左上点经度(单位：1/10^6 度)
	RightBottomLat uint32        // 右下点纬度(单位：1/10^6 度)
	RightBottomLng uint32        // 右下点经度(单位：1/10^6 度)
	StartTime      time.Time     // 起始时间，若区域属性0位为0则没有该字段
	EndTime        time.Time     // 结束时间，若区域属性0位为0则没有该字段
	MaxSpeed       uint16        // 最高速度(单位：公里/小时)，若区域属性1位为0则没有该字段
	SpeedDuration  byte          // 超速持续时间(单位：秒)，若区域属性1位为0则没有该字段
	NightMaxSpeed  uint16        // 夜间最高速度(单位：公里/小时)，若区域属性1位为0则没有该字段（2019版本）
	AreaName       string        // 区域名称（2019版本）
}

// T808_0x8602 设置矩形区域
// 2013版本和2019版本通用（2019版本新增夜间最高速度和区域名称）
type T808_0x8602 struct {
	AreaCount      byte                        // 区域总数
	AreaSettingTag byte                        // 区域设置属性：0-更新区域；1-追加区域；2-修改区域
	RectangleAreas []T808_0x8602_RectangleArea // 矩形区域项

	protocolVersion VersionType // 协议版本
}

// SetProtocolVersion 设置协议版本
func (entity *T808_0x8602) SetProtocolVersion(protocolVersion VersionType) {
	entity.protocolVersion = protocolVersion
}

// GetProtocolVersion 获取协议版本
func (entity *T808_0x8602) GetProtocolVersion() VersionType {
	return entity.protocolVersion
}

// MsgID 获取消息ID
func (entity *T808_0x8602) MsgID() MsgID {
	return MsgT808_0x8602
}

// Encode 编码消息
func (entity *T808_0x8602) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入区域设置属性
	writer.WriteByte(entity.AreaSettingTag)

	// 写入区域总数
	writer.WriteByte(entity.AreaCount)

	// 写入矩形区域项
	for _, area := range entity.RectangleAreas {
		// 写入区域ID
		writer.WriteUint32(area.AreaID)

		// 写入区域属性
		writer.WriteUint16(uint16(area.AreaAttribute))

		// 写入左上点纬度
		writer.WriteUint32(area.LeftTopLat)

		// 写入左上点经度
		writer.WriteUint32(area.LeftTopLng)

		// 写入右下点纬度
		writer.WriteUint32(area.RightBottomLat)

		// 写入右下点经度
		writer.WriteUint32(area.RightBottomLng)

		// 根据区域属性，写入起始和结束时间
		if area.AreaAttribute.GetTimeRange() {
			writer.WriteBcdTime(area.StartTime)
			writer.WriteBcdTime(area.EndTime)
		}

		// 根据区域属性，写入最高速度和超速持续时间
		if area.AreaAttribute.GetSpeedLimit() {
			writer.WriteUint16(area.MaxSpeed)
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
				return nil, err
			}
			writer.WriteWord(uint16(length))

			// 写入名称
			err = writer.WriteString(area.AreaName)
			if err != nil {
				return nil, err
			}
		}
	}

	return writer.Bytes(), nil
}

// Decode 解码消息
func (entity *T808_0x8602) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	var err error
	// 读取区域设置属性
	entity.AreaSettingTag, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取区域总数
	entity.AreaCount, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取矩形区域项
	entity.RectangleAreas = make([]T808_0x8602_RectangleArea, entity.AreaCount)
	for i := 0; i < int(entity.AreaCount); i++ {
		// 读取区域ID
		areaID, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RectangleAreas[i].AreaID = areaID

		// 读取区域属性
		areaAttr, err := reader.ReadUint16()
		if err != nil {
			return 0, err
		}
		entity.RectangleAreas[i].AreaAttribute = AreaAttribute(areaAttr)

		// 读取左上点纬度
		leftTopLat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RectangleAreas[i].LeftTopLat = leftTopLat

		// 读取左上点经度
		leftTopLng, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RectangleAreas[i].LeftTopLng = leftTopLng

		// 读取右下点纬度
		rightBottomLat, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RectangleAreas[i].RightBottomLat = rightBottomLat

		// 读取右下点经度
		rightBottomLng, err := reader.ReadUint32()
		if err != nil {
			return 0, err
		}
		entity.RectangleAreas[i].RightBottomLng = rightBottomLng

		// 根据区域属性，读取起始和结束时间
		if entity.RectangleAreas[i].AreaAttribute.GetTimeRange() {
			startTime, err := reader.ReadBcdTime()
			if err != nil {
				return 0, err
			}
			entity.RectangleAreas[i].StartTime = startTime

			endTime, err := reader.ReadBcdTime()
			if err != nil {
				return 0, err
			}
			entity.RectangleAreas[i].EndTime = endTime
		}

		// 根据区域属性，读取最高速度和超速持续时间
		if entity.RectangleAreas[i].AreaAttribute.GetSpeedLimit() {
			maxSpeed, err := reader.ReadUint16()
			if err != nil {
				return 0, err
			}
			entity.RectangleAreas[i].MaxSpeed = maxSpeed

			speedDuration, err := reader.ReadByte()
			if err != nil {
				return 0, err
			}
			entity.RectangleAreas[i].SpeedDuration = speedDuration

			// 根据协议版本，读取夜间最高速度
			if entity.protocolVersion == Version2019 {
				nightMaxSpeed, err := reader.ReadUint16()
				if err != nil {
					return 0, err
				}
				entity.RectangleAreas[i].NightMaxSpeed = nightMaxSpeed
			}
		}

		// 根据协议版本，读取区域名称
		if entity.protocolVersion == Version2019 {
			length, err := reader.ReadWord()
			if err != nil {
				return 0, err
			}
			if length > 0 {
				areaName, err := reader.ReadString(int(length))
				if err != nil {
					return 0, err
				}
				entity.RectangleAreas[i].AreaName = areaName
			}
		}
	}
	return len(data) - reader.Len(), nil
}
