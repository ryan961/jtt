package jtt

import "fmt"

// T808_0x0801 多媒体数据上传
type T808_0x0801 struct {
	// 多媒体ID，值大于零
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
	MultimediaFormat byte
	// 事件项编码
	// 0: 平台下发指令
	// 1: 定时动作
	// 2: 抢劫报警触发
	// 3: 碰撞侧翻报警触发
	// 4: 打开车门
	// 5: 关闭车门
	EventCode byte
	// 通道ID
	ChannelID byte
	// 位置信息汇报(0x0200)消息体，表示多媒体数据的位置基本信息数据见表23
	Location T808_0x0200
	// 多媒体数据包
	MultimediaData []byte
}

func (m *T808_0x0801) MsgID() MsgID { return MsgT808_0x0801 }

func (m *T808_0x0801) Encode() ([]byte, error) {
	w := NewWriter()

	// 写入多媒体ID
	w.WriteDWord(m.MultimediaID)

	// 写入多媒体类型
	w.WriteByte(m.MultimediaType)

	// 写入多媒体格式编码
	w.WriteByte(m.MultimediaFormat)

	// 写入事件项编码
	w.WriteByte(m.EventCode)

	// 写入通道ID
	w.WriteByte(m.ChannelID)

	// 写入位置信息汇报，固定28字节
	locationBytes, err := m.Location.Encode()
	if err != nil {
		return nil, fmt.Errorf("encode location info: %w", err)
	}
	if len(locationBytes) != 28 {
		return nil, fmt.Errorf("invalid location info length: %d", len(locationBytes))
	}
	w.Write(locationBytes)

	// 写入多媒体数据包
	w.Write(m.MultimediaData)

	return w.Bytes(), nil
}

func (m *T808_0x0801) Decode(data []byte) (int, error) {
	if len(data) < 36 { // 最少36字节：4+1+1+1+1+28
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	// 读取多媒体ID
	if m.MultimediaID, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read multimedia id: %w", err)
	}

	// 读取多媒体类型
	if m.MultimediaType, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read multimedia type: %w", err)
	}

	// 读取多媒体格式编码
	if m.MultimediaFormat, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read multimedia format: %w", err)
	}

	// 读取事件项编码
	if m.EventCode, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read event code: %w", err)
	}

	// 读取通道ID
	if m.ChannelID, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read channel id: %w", err)
	}

	// 读取位置信息汇报，固定28字节
	locationBytes, err := r.Read(28)
	if err != nil {
		return 0, fmt.Errorf("read location info: %w", err)
	}

	if _, err = m.Location.Decode(locationBytes); err != nil {
		return 0, fmt.Errorf("decode location info: %w", err)
	}

	// 读取剩余的多媒体数据包
	remaining := r.Len()
	if remaining > 0 {
		if m.MultimediaData, err = r.Read(remaining); err != nil {
			return 0, fmt.Errorf("read multimedia data: %w", err)
		}
	}

	return len(data) - r.Len(), nil
}
