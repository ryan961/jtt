package jtt

// T808_0x0E11 驾驶员身份库数据下载应答
type T808_0x0E11 struct {
	// 应答流水号
	ReplyMsgSerialNo uint16 `json:"replyMsgSerialNo"`
	// 应答结果
	// 0：成功，1：失败
	Result byte `json:"result"`
	// 需要下载的总数
	TotalCount byte `json:"totalCount"`
	// 当前下载到第几个文件
	CurrentIndex byte `json:"currentIndex"`
	// 当前下载的人脸 ID
	CurrentFaceID string `json:"currentFaceID"`
}

func (msg *T808_0x0E11) MsgID() MsgID {
	return MsgT808_0x0E11
}

func (msg *T808_0x0E11) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入应答流水号
	writer.WriteWord(msg.ReplyMsgSerialNo)

	// 写入应答结果
	writer.WriteByte(msg.Result)

	// 写入需要下载的总数
	writer.WriteByte(msg.TotalCount)

	// 写入当前下载到第几个文件
	writer.WriteByte(msg.CurrentIndex)

	// 写入当前下载的人脸 ID 长度
	writer.WriteByte(byte(len(msg.CurrentFaceID)))

	// 写入当前下载的人脸 ID
	err := writer.WriteString(msg.CurrentFaceID)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func (msg *T808_0x0E11) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取应答流水号
	var err error
	msg.ReplyMsgSerialNo, err = reader.ReadWord()
	if err != nil {
		return 0, err
	}

	// 读取应答结果
	msg.Result, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取需要下载的总数
	msg.TotalCount, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取当前下载到第几个文件
	msg.CurrentIndex, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取当前下载的人脸 ID 长度
	faceIDLen, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取当前下载的人脸 ID
	msg.CurrentFaceID, err = reader.ReadString(int(faceIDLen))
	if err != nil {
		return 0, err
	}

	return len(data) - reader.Len(), nil
}
