package jtt

import (
	"fmt"
	"math"
)

// T808_0x8800 多媒体数据上传应答
type T808_0x8800 struct {
	// 多媒体ID，>0，如收到全部数据包则没有后续字段
	MultimediaID uint32
	// 重传包ID列表，重传包序号顺序排列，如"包ID1 包ID2...包IDn"
	RetransmitIDs []uint16
}

func (m *T808_0x8800) MsgID() MsgID { return MsgT808_0x8800 }

func (m *T808_0x8800) Encode() ([]byte, error) {
	w := NewWriter()

	// 写入多媒体ID
	w.WriteDWord(m.MultimediaID)

	if len(m.RetransmitIDs) > math.MaxUint8 {
		return nil, fmt.Errorf("retransmit ids count too large: %d", len(m.RetransmitIDs))
	}

	// 写入重传包总数
	w.WriteByte(byte(len(m.RetransmitIDs)))

	// 写入重传包ID列表
	if len(m.RetransmitIDs) > 0 {
		for _, id := range m.RetransmitIDs {
			w.WriteWord(id)
		}
	}

	return w.Bytes(), nil
}

func (m *T808_0x8800) Decode(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error

	// 读取多媒体ID
	if m.MultimediaID, err = r.ReadDWord(); err != nil {
		return 0, fmt.Errorf("read multimedia id: %w", err)
	}

	// 如果还有数据，读取重传包信息
	if r.Len() > 0 {
		// 读取重传包总数
		count, err := r.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("read retransmit count: %w", err)
		}

		m.RetransmitIDs = make([]uint16, 0, int(count))
		for i := 0; i < int(count); i++ {
			id, err := r.ReadWord()
			if err != nil {
				return 0, fmt.Errorf("read retransmit id %d: %w", i, err)
			}
			m.RetransmitIDs = append(m.RetransmitIDs, id)
		}
	}

	return len(data) - r.Len(), nil
}
