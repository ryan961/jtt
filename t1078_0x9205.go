package jtt

import (
	"time"
)

// T1078_0x9205 查询资源列表
type T1078_0x9205 = DeviceMediaQuery

func (entity *T1078_0x9205) MsgID() MsgID {
	return MsgT1078_0x9205
}

type DeviceMediaQuery struct {
	LogicChannelID byte      `json:"logicChannelId"` // 逻辑通道号
	StartTime      time.Time `json:"startTime"`      // YY-MM-DD-HH-MM-SS，全 0 表示无起始时间条件
	EndTime        time.Time `json:"endTime"`        // YY-MM-DD-HH-MM-SS，全 0 表示无终止时间条件
	AlarmSign      [2]uint32 `json:"alarmSign"`      // 报警标志位。bit0-bit31为0x0200的报警标志位，bit32-bit63，全0表示无报警类型条件
	MediaType      byte      `json:"mediaType"`      // 音视频类型。0：音视频；1：音频；2：视频；3：视频或音视频
	StreamType     byte      `json:"streamType"`     // 码流类型。0：所有码流；1：主码流；2：子码流
	StorageType    byte      `json:"storageType"`    // 存储器类型。0：所有存储器；1：主存储器；2：灾备存储器
}

func (entity *DeviceMediaQuery) Encode() ([]byte, error) {
	writer := NewWriter()
	writer.WriteByte(entity.LogicChannelID)
	writer.WriteBcdTime(entity.StartTime)
	writer.WriteBcdTime(entity.EndTime)
	writer.WriteUint32(entity.AlarmSign[0])
	writer.WriteUint32(entity.AlarmSign[1])
	writer.WriteByte(entity.MediaType)
	writer.WriteByte(entity.StreamType)
	writer.WriteByte(entity.StorageType)
	return writer.Bytes(), nil
}

func (entity *DeviceMediaQuery) Decode(data []byte) (int, error) {
	if len(data) < 24 {
		return 0, ErrInvalidBody
	}
	reader := NewReader(data)

	var err error
	entity.LogicChannelID, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.StartTime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	entity.EndTime, err = reader.ReadBcdTime()
	if err != nil {
		return 0, err
	}

	var alarmSign [2]uint32
	alarmSign[0], err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	alarmSign[1], err = reader.ReadUint32()
	if err != nil {
		return 0, err
	}
	entity.AlarmSign = alarmSign

	entity.MediaType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.StreamType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	entity.StorageType, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
