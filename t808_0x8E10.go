package jtt

// T808_0x8E10 驾驶员身份识别上报应答
type T808_0x8E10 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16 `json:"replyMsgSerialNo"`
	// 重传包总数
	RetransmitCount uint16 `json:"retransmitCount"`
	// 重传包 ID
	RetransmitIDs []uint16 `json:"retransmitIDs"`
}

func (msg *T808_0x8E10) MsgID() MsgID {
	return MsgT808_0x8E10
}

func (msg *T808_0x8E10) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteWord(msg.ReplyMsgSerialNo)

	// 写入重传包总数
	writer.WriteWord(msg.RetransmitCount)

	// 写入重传包 ID
	for _, id := range msg.RetransmitIDs {
		writer.WriteWord(id)
	}

	return writer.Bytes(), nil
}

func (msg *T808_0x8E10) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取应答流水号
	var err error
	msg.ReplyMsgSerialNo, err = reader.ReadWord()
	if err != nil {
		return 0, err
	}

	// 读取重传包总数
	msg.RetransmitCount, err = reader.ReadWord()
	if err != nil {
		return 0, err
	}

	// 读取重传包 ID
	msg.RetransmitIDs = make([]uint16, msg.RetransmitCount)
	for i := 0; i < int(msg.RetransmitCount); i++ {
		msg.RetransmitIDs[i], err = reader.ReadWord()
		if err != nil {
			return 0, err
		}
	}

	return len(data) - reader.Len(), nil
}
