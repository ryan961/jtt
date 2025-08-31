package jtt

import "fmt"

// T808_0x0302 提问应答
// 0 应答流水号 WORD
// 2 答案ID BYTE
// 对应上行：对 0x8302 的应答

type T808_0x0302 struct {
	RespSerialNo uint16
	AnswerID     byte
}

func (m *T808_0x0302) MsgID() MsgID { return MsgT808_0x0302 }

func (m *T808_0x0302) Encode() ([]byte, error) {
	w := NewWriter()
	w.WriteWord(m.RespSerialNo)
	w.WriteByte(m.AnswerID)
	return w.Bytes(), nil
}

func (m *T808_0x0302) Decode(data []byte) (int, error) {
	if len(data) < 3 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error
	if m.RespSerialNo, err = r.ReadWord(); err != nil {
		return 0, fmt.Errorf("read RespSerialNo: %w", err)
	}
	if m.AnswerID, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read AnswerID: %w", err)
	}
	return len(data) - r.Len(), nil
}
