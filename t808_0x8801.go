package jtt

import "fmt"

// T808_0x8801 摄像头立即拍摄命令
type T808_0x8801 struct {
	// 通道ID，值大于零
	ChannelID byte
	// 拍摄命令
	// 0: 停止拍摄
	// 0xFFFF: 录像
	// 其他: 拍照张数
	Command uint16
	// 拍摄间隔/录像时间，单位为秒(s)
	// 0: 按最小时间间隔拍照或一直录像
	Interval uint16
	// 保存标志
	// 1: 保存；0: 实时上传
	SaveFlag byte
	// 分辨率
	// 0x00: 最低分辨率
	// 0x01: 320×240
	// 0x02: 640×480
	// 0x03: 800×600
	// 0x04: 1024×768
	// 0x05: 176×144[Qcif]
	// 0x06: 352×288[Cif]
	// 0x07: 704×288[HALF D1]
	// 0x08: 704×576[D1]
	// 0xff: 最高分辨率
	Resolution byte
	// 图像/视频质量
	// 取值范围1~10，1代表质量损失最小，10表示压缩比最大
	Quality byte
	// 亮度，0~255
	Brightness byte
	// 对比度，0~127
	Contrast byte
	// 饱和度，0~127
	Saturation byte
	// 色度，0~255
	Chroma byte
}

func (m *T808_0x8801) MsgID() MsgID { return MsgT808_0x8801 }

func (m *T808_0x8801) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.ChannelID)
	w.WriteWord(m.Command)
	w.WriteWord(m.Interval)
	w.WriteByte(m.SaveFlag)
	w.WriteByte(m.Resolution)
	w.WriteByte(m.Quality)
	w.WriteByte(m.Brightness)
	w.WriteByte(m.Contrast)
	w.WriteByte(m.Saturation)
	w.WriteByte(m.Chroma)
	return w.Bytes(), nil
}

func (m *T808_0x8801) Decode(data []byte) (int, error) {
	if len(data) < 11 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.ChannelID, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read channel id: %w", err)
	}

	if m.Command, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read command: %w", err)
	}

	if m.Interval, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read interval: %w", err)
	}

	if m.SaveFlag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read save flag: %w", err)
	}

	if m.Resolution, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read resolution: %w", err)
	}

	if m.Quality, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read quality: %w", err)
	}

	if m.Brightness, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read brightness: %w", err)
	}

	if m.Contrast, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read contrast: %w", err)
	}

	if m.Saturation, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read saturation: %w", err)
	}

	if m.Chroma, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read chroma: %w", err)
	}

	return len(data) - r.Len(), nil
}
