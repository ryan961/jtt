package jtt

import "fmt"

// T808_0x0802 存储多媒体数据检索应答
type T808_0x0802 struct {
	// 应答流水号，对应的多媒体数据检索消息的流水号
	ReplyMsgSerialNo uint16
	// 检索项
	Items []T808_0x0802_MultimediaItem
}

func (m *T808_0x0802) MsgID() MsgID { return MsgT808_0x0802 }

// T808_0x0802_MultimediaItem 多媒体检索项
type T808_0x0802_MultimediaItem struct {
	// 多媒体ID，值大于0
	MultimediaID uint32
	// 多媒体类型
	// 0: 图像
	// 1: 音频
	// 2: 视频
	MultimediaType byte
	// 通道ID
	ChannelID byte
	// 事件项编码
	// 0: 平台下发指令
	// 1: 定时动作
	// 2: 抢劫报警触发
	// 3: 碰撞侧翻报警触发
	// 其他保留
	EventCode byte
	// 位置信息汇报(0x0200)消息体（28字节），表示拍摄或录制的起始时刻的汇报信息
	Location T808_0x0200
}

func (m *T808_0x0802) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteWord(m.ReplyMsgSerialNo)
	w.WriteWord(uint16(len(m.Items)))

	for _, item := range m.Items {
		w.WriteDWord(item.MultimediaID)
		w.WriteByte(item.MultimediaType)
		w.WriteByte(item.ChannelID)
		w.WriteByte(item.EventCode)

		// 编码位置信息
		locationBytes, err := item.Location.Encode()
		if err != nil {
			return nil, fmt.Errorf("encode location info: %w", err)
		}
		w.Write(locationBytes)
	}

	return w.Bytes(), nil
}

func (m *T808_0x0802) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.ReplyMsgSerialNo, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read reply msg serial no: %w", err)
	}

	var totalCount uint16
	if totalCount, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read total count: %w", err)
	}

	// 读取检索项
	m.Items = make([]T808_0x0802_MultimediaItem, 0, totalCount)
	for i := uint16(0); i < totalCount && r.Len() >= 35; i++ { // 最少35字节：4+1+1+1+28
		var item T808_0x0802_MultimediaItem

		if item.MultimediaID, err = r.ReadDWord(); err != nil {
			return 0, fmt.Errorf("read multimedia id: %w", err)
		}

		if item.MultimediaType, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read multimedia type: %w", err)
		}

		if item.ChannelID, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read channel id: %w", err)
		}

		if item.EventCode, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read event code: %w", err)
		}

		// 读取位置信息，固定28字节
		locationBytes, err := r.Read(28)
		if err != nil {
			return 0, fmt.Errorf("read location info: %w", err)
		}

		if _, err = item.Location.Decode(locationBytes); err != nil {
			return 0, fmt.Errorf("decode location info: %w", err)
		}

		m.Items = append(m.Items, item)
	}

	return len(data) - r.Len(), nil
}
