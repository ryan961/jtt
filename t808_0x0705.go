package jtt

import (
	"fmt"
	"time"
)

// T808_0x0705 CAN总线数据上传
type T808_0x0705 struct {
	// 数据项个数，包含的CAN总线数据项个数，值大于0
	Count uint16
	// CAN总线数据接收时间，第1条CAN总线数据的接收时间，hh-mm-ss-msms
	ReceiveTime time.Time
	// CAN总线数据项
	Items []T808_0x0705_CANItem
}

func (m *T808_0x0705) MsgID() MsgID { return MsgT808_0x0705 }

// T808_0x0705_CANItem CAN总线数据项
type T808_0x0705_CANItem struct {
	// CAN ID (DWORD)
	// bit31: CAN通道号，0:CAN1，1:CAN2
	// bit30: 帧类型，0:标准帧，1:扩展帧
	// bit29: 数据采集方式，0:原始数据，1:采集区间的平均值
	// bit28-bit0: CAN总线ID
	CANID T808_0x0705_CANID
	// CAN数据，固定8字节
	CANData [8]byte
}

// T808_0x0705_CANID CAN ID
//
//	bit31: CAN通道号，0:CAN1，1:CAN2
//	bit30: 帧类型，0:标准帧，1:扩展帧
//	bit29: 数据采集方式，0:原始数据，1:采集区间的平均值
//	bit28-bit0: CAN总线ID
type T808_0x0705_CANID uint32

// 获取CAN通道号（bit31）
func (id T808_0x0705_CANID) GetChannel() byte {
	return byte(id >> 31)
}

// 设置CAN通道号（bit31）
func (id T808_0x0705_CANID) SetChannel(channel byte) T808_0x0705_CANID {
	return id | (T808_0x0705_CANID(channel) << 31)
}

// 设置帧类型（bit30）
func (id *T808_0x0705_CANID) SetFrameType(frameType byte) {
	*id = *id | (T808_0x0705_CANID(frameType) << 30)
}

// 设置数据采集方式（bit29）
func (id *T808_0x0705_CANID) SetDataType(dataType byte) {
	*id = *id | (T808_0x0705_CANID(dataType) << 29)
}

// 设置CAN总线ID（bit28-bit0）
func (id *T808_0x0705_CANID) SetCANID(canID uint32) {
	*id = *id | (T808_0x0705_CANID(canID) & 0x1FFFFFFF)
}

// 获取帧类型（bit30）
func (id T808_0x0705_CANID) GetFrameType() byte {
	return byte(id >> 30)
}

// 获取数据采集方式（bit29）
func (id T808_0x0705_CANID) GetDataType() byte {
	return byte(id >> 29)
}

// 获取CAN总线ID（bit28-bit0）
func (id T808_0x0705_CANID) GetCANID() uint32 {
	return uint32(id & 0x1FFFFFFF)
}

func (m *T808_0x0705) Encode() ([]byte, error) {
	w := NewWriter()

	// 写入数据项个数
	w.WriteWord(uint16(len(m.Items)))

	// 写入CAN总线数据接收时间，BCD[5] hh-mm-ss-msms
	timeStr := m.ReceiveTime.Format("15040505")
	if len(timeStr) < 8 {
		timeStr = timeStr + "00" // 补充毫秒部分
	}
	w.WriteBcd(timeStr[:8], 5)

	// 写入CAN总线数据项
	for _, item := range m.Items {
		w.WriteDWord(uint32(item.CANID))
		w.Write(item.CANData[:])
	}

	return w.Bytes(), nil
}

func (m *T808_0x0705) Decode(data []byte) (int, error) {
	if len(data) < 7 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	// 读取数据项个数
	if m.Count, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read count: %w", err)
	}

	// 读取CAN总线数据接收时间，BCD[5]
	timeStr, err := r.ReadBcd(5)
	if err != nil {
		return 0, fmt.Errorf("read receive time: %w", err)
	}

	// 解析时间格式 hh-mm-ss-msms
	if len(timeStr) >= 8 {
		if m.ReceiveTime, err = time.Parse("15040505", timeStr[:8]); err != nil {
			return 0, fmt.Errorf("parse receive time: %w", err)
		}
	}

	// 读取CAN总线数据项
	m.Items = make([]T808_0x0705_CANItem, 0, int(m.Count))
	for i := 0; i < int(m.Count) && r.Len() >= 12; i++ {
		var item T808_0x0705_CANItem

		// 读取CAN ID
		if canID, err := r.ReadDWord(); err != nil {
			return 0, fmt.Errorf("read CAN ID %d: %w", i, err)
		} else {
			item.CANID = T808_0x0705_CANID(canID)
		}

		// 读取CAN数据，固定8字节
		canData, err := r.Read(8)
		if err != nil {
			return 0, fmt.Errorf("read CAN data %d: %w", i, err)
		}
		copy(item.CANData[:], canData)

		m.Items = append(m.Items, item)
	}

	return len(data) - r.Len(), nil
}
