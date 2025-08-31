package jtt

import "fmt"

// T808_0x0800 多媒体事件信息上传
type T808_0x0800 struct {
	// 多媒体数据ID，值大于0
	MultimediaID uint32
	// 多媒体类型
	// 0: 图像
	// 1: 音频
	// 2: 视频
	MultimediaType byte
	// 多媒体格式编码
	// 0: JPEG
	// 1: TIF
	// 2: MP3
	// 3: WAV
	// 4: WMV
	// 其他保留
	MultimediaFormat byte
	// 事件项编码
	// 0: 平台下发指令
	// 1: 定时动作
	// 2: 抢劫报警触发
	// 3: 碰撞侧翻报警触发
	// 4: 门开拍照
	// 5: 门关拍照
	// 6: 车门由开变关，车速从小于20km到超过20km
	// 7: 定时拍照
	EventCode byte
	// 通道ID
	ChannelID byte
}

func (m *T808_0x0800) MsgID() MsgID { return MsgT808_0x0800 }

func (m *T808_0x0800) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteDWord(m.MultimediaID)
	w.WriteByte(m.MultimediaType)
	w.WriteByte(m.MultimediaFormat)
	w.WriteByte(m.EventCode)
	w.WriteByte(m.ChannelID)
	return w.Bytes(), nil
}

func (m *T808_0x0800) Decode(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.MultimediaID, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read multimedia id: %w", err)
	}

	if m.MultimediaType, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read multimedia type: %w", err)
	}

	if m.MultimediaFormat, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read multimedia format: %w", err)
	}

	if m.EventCode, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read event code: %w", err)
	}

	if m.ChannelID, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read channel id: %w", err)
	}

	return len(data) - r.Len(), nil
}
