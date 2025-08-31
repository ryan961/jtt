package jtt

import "fmt"

// T808_0x8804 录音开始命令
type T808_0x8804 struct {
	// 录音命令
	// 0: 停止录音
	// 0x01: 开始录音
	Command byte
	// 录音时间，单位为秒(s)，0表示一直录音
	Duration uint16
	// 保存标志
	// 0: 实时上传
	// 1: 保存
	SaveFlag byte
	// 音频采样率
	// 0: 8K
	// 1: 11K
	// 2: 23K
	// 3: 32K
	// 其他保留
	SampleRate byte
}

func (m *T808_0x8804) MsgID() MsgID { return MsgT808_0x8804 }

func (m *T808_0x8804) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteByte(m.Command)
	w.WriteWord(m.Duration)
	w.WriteByte(m.SaveFlag)
	w.WriteByte(m.SampleRate)
	return w.Bytes(), nil
}

func (m *T808_0x8804) Decode(data []byte) (int, error) {
	if len(data) < 5 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	if m.Command, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read command: %w", err)
	}

	if m.Duration, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read duration: %w", err)
	}

	if m.SaveFlag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read save flag: %w", err)
	}

	if m.SampleRate, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read sample rate: %w", err)
	}

	return len(data) - r.Len(), nil
}
