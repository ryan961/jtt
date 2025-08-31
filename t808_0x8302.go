package jtt

import (
	"fmt"
)

// T808_0x8302 提问下发
type T808_0x8302 struct {
	// 提问下发标志位 BYTE
	// 0：紧急
	// 1～2：保留
	// 3：终端TTS播读
	// 4：广告屏显示
	// 5～7：保留
	Flag byte
	// 问题文本，经 GBK 编码
	Question string
	// 候选答案列表
	Answers []T808_0x8302_Answer
}

type T808_0x8302_Answer struct {
	// 答案ID BYTE
	// 0～255：答案ID
	ID byte
	// 答案内容，经 GBK 编码
	Content string
}

func (m *T808_0x8302) MsgID() MsgID { return MsgT808_0x8302 }

func (m *T808_0x8302) Encode() ([]byte, error) {
	w := NewWriter()
	// 标志位
	w.WriteByte(m.Flag)
	// 问题长度 + 问题
	ln, err := GB18030Length(m.Question)
	if err != nil {
		return nil, err
	}
	if ln > 0xFF { // 问题字段使用 BYTE 长度
		return nil, fmt.Errorf("question too long: %d", ln)
	}
	w.WriteByte(byte(ln)) // 问题长度
	if ln > 0 {
		if err := w.WriteString(m.Question); err != nil {
			return nil, err
		}
	}
	// 答案列表
	for _, a := range m.Answers {
		w.WriteByte(a.ID)
		aln, err := GB18030Length(a.Content)
		if err != nil {
			return nil, err
		}
		w.WriteWord(uint16(aln)) // 答案内容长度
		if aln > 0 {
			if err := w.WriteString(a.Content); err != nil {
				return nil, err
			}
		}
	}
	return w.Bytes(), nil
}

func (m *T808_0x8302) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, fmt.Errorf("invalid data length: %d", len(data))
	}
	r := NewReader(data)
	var err error
	if m.Flag, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read Flag: %w", err)
	}
	var questionLen byte
	if questionLen, err = r.ReadByte(); err != nil {
		return 0, fmt.Errorf("read question length: %w", err)
	}
	if questionLen > 0 {
		if m.Question, err = r.ReadString(int(questionLen)); err != nil {
			return 0, fmt.Errorf("read Question: %w", err)
		}
	}
	m.Answers = make([]T808_0x8302_Answer, 0)
	for r.Len() > 0 {
		var id byte
		if id, err = r.ReadByte(); err != nil {
			return 0, fmt.Errorf("read id: %w", err)
		}
		var answerLen uint16
		if answerLen, err = r.ReadWord(); err != nil {
			return 0, fmt.Errorf("read answer length: %w", err)
		}
		ans := T808_0x8302_Answer{ID: id}
		if answerLen > 0 {
			if ans.Content, err = r.ReadString(int(answerLen)); err != nil {
				return 0, fmt.Errorf("read answer content: %w", err)
			}
		}
		m.Answers = append(m.Answers, ans)
	}
	return len(data) - r.Len(), nil
}
