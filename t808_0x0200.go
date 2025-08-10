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

// T808_0x0200_Status 位置状态
type T808_0x0200_Status uint32

// GetAccState 获取Acc状态
func (status T808_0x0200_Status) GetAccState() bool {
	return GetBitUint32(uint32(status), 0)
}

// Positioning 是否正在定位
func (status T808_0x0200_Status) Positioning() bool {
	return GetBitUint32(uint32(status), 1)
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

// T808_0x0200 汇报位置
type T808_0x0200 struct {
	// 警告
	Alarm uint32 `json:"alarm"`
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
	writer.WriteUint32(msg.Alarm)

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
	msg.Alarm, err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取状态信息
	status, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	msg.Status = T808_0x0200_Status(status)

	// 读取纬度信息
	latitude, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	// 读取经度信息
	longitude, err := reader.ReadUint32()
	if err != nil {
		return 0, err
	}

	msg.Lat, msg.Lng = GetGeoPointForWGS84(
		latitude, msg.Status.GetLatitudeType() == SouthLatitudeType,
		longitude, msg.Status.GetLongitudeType() == WestLongitudeType,
	)

	// 读取海拔高度
	msg.Altitude, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取行驶速度
	msg.Speed, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取行驶方向
	msg.Direction, err = reader.ReadUint16()
	if err != nil {
		return 0, err
	}

	// 读取上报时间
	msg.Time, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
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
			Id:     Type(id),
			Length: byte(length),
			Data:   buffer[:length],
		})
		buffer = buffer[length:]
	}
	if len(extras) > 0 {
		msg.Extras = extras
	}
	return len(data) - reader.Len(), nil
}
