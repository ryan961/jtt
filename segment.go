package jtt

// Segment 分包消息结构.
type Segment struct {
	PhoneNumber string            `json:"phoneNumber"`
	MsgID       MsgID             `json:"msgId"`
	Total       uint16            `json:"total"`
	Data        map[uint16][]byte `json:"data"`
}

func (s *Segment) IsComplete() bool {
	return s.Total == uint16(len(s.Data))
}

func (s *Segment) Merge(info *SegmentInfo, body []byte) {
	s.Data[info.Index] = body
}

func (s *Segment) GetBody() []byte {
	body := make([]byte, 0, s.Total*1023) // 预分配容量，避免扩容。分包长度固定最大 1023
	for i := uint16(1); i <= s.Total; i++ {
		body = append(body, s.Data[i]...)
	}
	return body
}

func (s *Segment) Reset() {
	s.PhoneNumber = ""
	s.MsgID = MsgID(0)
	s.Total = 0
	// 清空map而不是重新分配，可以重用内存
	for k := range s.Data {
		delete(s.Data, k)
	}
}
