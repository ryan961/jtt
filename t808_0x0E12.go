package jtt

// T808_0x0E12 驾驶员身份库查询应答
type T808_0x0E12 struct {
	// 人脸库列表个数
	FaceCount byte `json:"faceCount"`
	// 人脸库信息列表
	FaceInfoList []T808_0x0E12_FaceInfo `json:"faceInfoList"`
}

// T808_0x0E12_FaceInfo 人脸信息列表数据格式
type T808_0x0E12_FaceInfo struct {
	// 人脸 ID
	FaceID string `json:"faceID"`
}

func (msg *T808_0x0E12) MsgID() MsgID {
	return MsgT808_0x0E12
}

func (msg *T808_0x0E12) Encode() ([]byte, error) {
	writer := NewWriter()

	// 写入人脸库列表个数
	writer.WriteByte(msg.FaceCount)

	// 写入人脸库信息列表
	for _, faceInfo := range msg.FaceInfoList {
		// 写入人脸 ID 长度
		writer.WriteByte(byte(len(faceInfo.FaceID)))
		// 写入人脸 ID
		err := writer.WriteString(faceInfo.FaceID)
		if err != nil {
			return nil, err
		}
	}

	return writer.Bytes(), nil
}

func (msg *T808_0x0E12) Decode(data []byte) (int, error) {
	reader := NewReader(data)

	// 读取人脸库列表个数
	var err error
	msg.FaceCount, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	// 读取人脸库信息列表
	msg.FaceInfoList = make([]T808_0x0E12_FaceInfo, msg.FaceCount)
	for i := 0; i < int(msg.FaceCount); i++ {
		// 读取人脸 ID 长度
		faceIDLen, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}

		// 读取人脸 ID
		msg.FaceInfoList[i].FaceID, err = reader.ReadString(int(faceIDLen))
		if err != nil {
			return 0, err
		}
	}

	return len(data) - reader.Len(), nil
}
