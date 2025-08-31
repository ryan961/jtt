package jtt

import (
	"fmt"
	"time"
)

// T808_0x8802 存储多媒体数据检索
type T808_0x8802 struct {
	// 多媒体类型
	// 0: 图像
	// 1: 音频
	// 2: 视频
	MultimediaType byte
	// 通道ID，0表示检索该媒体类型的所有通道
	ChannelID byte
	// 事件项编码
	// 0: 平台下发指令
	// 1: 定时动作
	// 2: 抢劫报警触发
	// 3: 碰撞侧翻报警触发
	// 其他保留
	EventCode byte
	// 起始时间，YY-MM-DD-hh-mm-ss
	StartTime time.Time
	// 结束时间，YY-MM-DD-hh-mm-ss
	EndTime time.Time
}

func (m *T808_0x8802) MsgID() MsgID { return MsgT808_0x8802 }

func (m *T808_0x8802) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.MultimediaType)
	w.WriteByte(m.ChannelID)
	w.WriteByte(m.EventCode)
	w.WriteBcdTime(m.StartTime)
	w.WriteBcdTime(m.EndTime)
	return w.Bytes(), nil
}

func (m *T808_0x8802) Decode(data []byte) (int, error) {
	if len(data) < 15 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.MultimediaType, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read multimedia type: %w", err)
	}

	if m.ChannelID, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read channel id: %w", err)
	}

	if m.EventCode, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read event code: %w", err)
	}

	if m.StartTime, err = r.ReadBcdTime(); err != nil {
		return 0, fmt.Errorf("read start time: %w", err)
	}

	if m.EndTime, err = r.ReadBcdTime(); err != nil {
		return 0, fmt.Errorf("read end time: %w", err)
	}

	return len(data) - r.Len(), nil
}
